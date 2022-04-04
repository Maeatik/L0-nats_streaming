package db

import (
	"L0/back/cache"
	models2 "L0/back/models"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"sync"
)

//константы для подключения к БД
const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "pass"
	dbname   = "testModel"
)

// PostgresRepository - структура для работы с БД
type PostgresRepository struct {
	Db *sql.DB
}

// OpenConnection - Настройка подключения к БД
func OpenConnection() (*PostgresRepository, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	//psqlInfo := "host=localhost port=5432 user=postgres password=pass dbname=testModel sslmode=disable"
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Println(err)
	}
	//Проверка подключения
	err = db.Ping()
	if err != nil {
		log.Println(err)
	}
	//Возвращение указателя на новую БД
	return &PostgresRepository{Db: db}, err
}

// Close - Закрытие подключения к БД
func (r *PostgresRepository) Close() {
	r.Db.Close()
}

// InsertModel - Добавление записей в таблицы БД
func (r *PostgresRepository) InsertModel(wg *sync.WaitGroup, model models2.Model) error {

	//Запросы в БД для добавления данных
	modelStatement := `INSERT INTO model VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	deliveryStatement := `INSERT INTO delivery VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	paymentStatement := `INSERT INTO payment VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	itemStatement := `INSERT INTO items VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	//Начало работы с БД
	toDB, err := r.Db.Begin()
	//Если соединение с БД не было установлено, состояение БД откатывается
	if err != nil {
		log.Println(err)
		toDB.Rollback()
		return err
	}
	//Добавление данных в таблицу model
	_, err = toDB.Exec(modelStatement,
		model.Order_uid,
		model.Track_number,
		model.Entry,
		model.Locale,
		model.Internal_signature,
		model.Customer_id,
		model.Delivery_service,
		model.Shardkey,
		model.Sm_id,
		model.Date_created,
		model.Oof_shard)
	//Если возникла ошибка, состояние БД откатывается
	if err != nil {
		log.Println(modelStatement, err)
		toDB.Rollback()
		return err
	}

	//Добавления данных в таблицу delivery
	_, err = toDB.Exec(deliveryStatement,
		model.Order_uid,
		model.Delivery.Name,
		model.Delivery.Phone,
		model.Delivery.Zip,
		model.Delivery.City,
		model.Delivery.Address,
		model.Delivery.Region,
		model.Delivery.Email)
	//Если возникла ошибка, состояние БД откатывается
	if err != nil {
		log.Println(deliveryStatement, err)
		toDB.Rollback()
		return err
	}
	//Добавления данных в таблицу payment
	_, err = toDB.Exec(paymentStatement,
		model.Payment.Transaction,
		model.Payment.Request_id,
		model.Payment.Currency,
		model.Payment.Provider,
		model.Payment.Amount,
		model.Payment.Payment_dt,
		model.Payment.Bank,
		model.Payment.Delivery_cost,
		model.Payment.Goods_total,
		model.Payment.Custom_fee)
	//Если возникла ошибка, состояние БД откатывается
	if err != nil {
		log.Println(paymentStatement, err)
		toDB.Rollback()
		return err
	}
	//Подготовка к добавление записей в таблицу items, так как одной модели могут принадлежать несколько item'ов
	insertValues, err := toDB.Prepare(itemStatement)
	//Если возникла ошибка, состояние БД откатывается
	if err != nil {
		log.Println(err)
		toDB.Rollback()
		return err
	}
	//закрытие подсоединения к таблице items
	defer insertValues.Close()
	//Добавление item'ов через цикл в таблицу items
	for _, item := range model.Items {
		_, err = insertValues.Exec(
			item.Chrt_id,
			item.Track_number,
			item.Price,
			item.Rid,
			item.Name,
			item.Sale,
			item.Size,
			item.Total_price,
			&item.Nm_id,
			item.Brand,
			item.Status,
		)
		//Если возникла ошибка, состояние БД откатывается
		if err != nil {
			log.Println(itemStatement, err)
			toDB.Rollback()
			return err
		}
	}
	//Если состояние БД не откатилось, данные заливаются в БД
	err = toDB.Commit()
	//Уменьшения счетчика ожидания
	wg.Done()
	return err
}

// Recovery - Воставновление\Получение данных из кеша
func (r *PostgresRepository) Recovery(modelMap *cache.Cache) error {
	//Запросы в БД для получения данных из таблиц
	modelStatement := "SELECT * FROM model"
	deliveryStatement := "SELECT name, phone, zip, city, address, region, email FROM delivery WHERE order_uid = $1"
	paymentStatement := "SELECT * FROM payment WHERE transaction = $1"
	itemStatement := "SELECT * FROM items WHERE track_number = $1"

	//Подготовка получения данных из главной таблицы model
	getModels, err := r.Db.Prepare(modelStatement)
	if err != nil {
		return err
	}
	defer getModels.Close()

	//Подготовка получения данных из таблицы delivery
	getDeliveries, err := r.Db.Prepare(deliveryStatement)
	if err != nil {
		return err
	}
	defer getDeliveries.Close()

	//Подготовка получения данных из таблицы payment
	getPayments, err := r.Db.Prepare(paymentStatement)
	if err != nil {
		return err
	}
	defer getPayments.Close()

	//Подготовка получения данных из таблицы items
	getItems, err := r.Db.Prepare(itemStatement)
	if err != nil {
		return err
	}
	defer getItems.Close()

	//Получение строк таблицы model
	rows, err := r.Db.Query(modelStatement)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer rows.Close()

	//Цикл чтение полученных из таблицы строк
	for rows.Next() {
		//Чтение полученных из таблицы model строк
		model := new(models2.Model)
		if err = rows.Scan(
			&model.Order_uid,
			&model.Track_number,
			&model.Entry,
			&model.Locale,
			&model.Internal_signature,
			&model.Customer_id,
			&model.Delivery_service,
			&model.Shardkey,
			&model.Sm_id,
			&model.Date_created,
			&model.Oof_shard); err != nil {
			log.Println(err)
			return err
		}
		//Чтение полученных из таблицы delivery строк
		//Данные подбираются по значению order_uid из таблицы model
		if err = getDeliveries.QueryRow(model.Order_uid).Scan(
			&model.Delivery.Name,
			&model.Delivery.Phone,
			&model.Delivery.Zip,
			&model.Delivery.City,
			&model.Delivery.Address,
			&model.Delivery.Region,
			&model.Delivery.Email); err != nil {
			log.Println(err)
			return err
		}
		//Чтение полученных из таблицы payment строк
		//Данные подбираются по значению order_uid из таблицы model
		if err = getPayments.QueryRow(model.Order_uid).Scan(
			&model.Payment.Transaction,
			&model.Payment.Request_id,
			&model.Payment.Currency,
			&model.Payment.Provider,
			&model.Payment.Amount,
			&model.Payment.Payment_dt,
			&model.Payment.Bank,
			&model.Payment.Delivery_cost,
			&model.Payment.Goods_total,
			&model.Payment.Custom_fee,
		); err != nil {
			log.Println(err)
			return err
		}
		//Так как строк в таблице items не обязательно соответствует кол-ву строк в таблице model, создается отдельный
		//читатель строк этой таблицы
		itemRows, err := getItems.Query(model.Track_number)
		if err != nil {
			log.Println(err)
			return err
		}
		defer itemRows.Close()

		//Чтение полученных из таблицы items строк
		for itemRows.Next() {
			item := new(models2.Items)
			//Чтение записей из таблицы items связанных с model
			if err = itemRows.Scan(
				&item.Chrt_id,
				&item.Track_number,
				&item.Price,
				&item.Rid,
				&item.Name,
				&item.Sale,
				&item.Size,
				&item.Total_price,
				&item.Nm_id,
				&item.Brand,
				&item.Status); err != nil {
				log.Println(err)
				return err
			}
			//Полученная запись добавляется в слайс item'ов
			model.Items = append(model.Items, *item)
		}
		//Все полученные данные добавляются в Кеш
		modelMap.AddModelCache(*model)
	}
	return err
}
