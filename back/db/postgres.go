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

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "pass"
	dbname   = "testModel"
)

type PostgresRepository struct {
	Db *sql.DB
}

func OpenConnection() (*PostgresRepository, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	//psqlInfo := "host=localhost port=5432 user=postgres password=pass dbname=MortyGRAB sslmode=disable"
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Println(err)
	}
	err = db.Ping()
	if err != nil {
		log.Println(err)
	}
	return &PostgresRepository{Db: db}, err
}

func (r *PostgresRepository) Close() {
	r.Db.Close()
}

func (r *PostgresRepository) InsertModel(wg *sync.WaitGroup, model models2.Model) error {

	modelStatement := `INSERT INTO model VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	deliveryStatement := `INSERT INTO delivery VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	paymentStatement := `INSERT INTO payment VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	itemStatement := `INSERT INTO items VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	toDB, err := r.Db.Begin()

	if err != nil{
		log.Println(err)
		toDB.Rollback()
		return err
	}

	insertValues, err := toDB.Prepare(itemStatement)
	if err != nil {
		log.Println(err)
		toDB.Rollback()
		return err
	}
	defer insertValues.Close()

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
	if err != nil{
		log.Println(modelStatement, err)
		toDB.Rollback()
		return err
	}

	_, err = toDB.Exec(deliveryStatement,
		model.Order_uid,
		model.Delivery.Name,
		model.Delivery.Phone,
		model.Delivery.Zip,
		model.Delivery.City,
		model.Delivery.Address,
		model.Delivery.Region,
		model.Delivery.Email)
	if err != nil{
		log.Println(deliveryStatement, err)
		toDB.Rollback()
		return err
	}

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

	if err != nil{
		log.Println(paymentStatement, err)
		toDB.Rollback()
		return err
	}

	for _, item := range model.Items{
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
		if err != nil{
			log.Println(itemStatement, err)
			toDB.Rollback()
			return err
		}
	}

	err = toDB.Commit()
	wg.Done()
	return err
}

func (r *PostgresRepository) Recovery(modelMap *cache.Cache)  error {

	modelStatement := "SELECT * FROM model"
	deliveryStatement := "SELECT name, phone, zip, city, address, region, email FROM delivery WHERE order_uid = $1"
	paymentStatement := "SELECT * FROM payment WHERE transaction = $1"
	itemStatement := "SELECT * FROM items WHERE track_number = $1"
	getModels, err := r.Db.Prepare(modelStatement)
	if err != nil{
		return err
	}
	defer getModels.Close()

	getDeliveries, err := r.Db.Prepare(deliveryStatement)
	if err != nil{
		return err
	}
	defer getDeliveries.Close()

	getPayments, err := r.Db.Prepare(paymentStatement)
	if err != nil{
		return err
	}
	defer getPayments.Close()

	getItems, err := r.Db.Prepare(itemStatement)
	if err != nil{
		return err
	}
	defer getItems.Close()

	rows, err := r.Db.Query(modelStatement)
	if err != nil{
		log.Fatal(err)
		return err
	}
	defer rows.Close()


	for rows.Next() {
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

		itemRows, err := getItems.Query(model.Track_number)
		if err != nil{
			log.Println(err)
			return err
		}
		defer itemRows.Close()

		for itemRows.Next() {
			item := new(models2.Items)

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
				&item.Status); err != nil{
				log.Println(err)
				return err
			}
			model.Items = append(model.Items, *item)
		}

		modelMap.AddModelCache(*model)
	}
	return err
}