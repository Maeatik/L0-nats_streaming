package cache

import (
	"L0/back/models"
	"sync"
)

// Cache - структура Кэша
type Cache struct {
	sync.RWMutex
	modelMap map[string]models.Model
}

// CreateCache - Создание нового кеша, путем создания мапы
func CreateCache() *Cache {
	return &Cache{
		modelMap: make(map[string]models.Model),
	}
}

// AddModelCache - Добавление модели в кэш. Во время добавления модель недоступна для изменения.
// После добавления вновь становится доступо
func (c *Cache) AddModelCache(model models.Model) {
	c.Lock()
	c.modelMap[*model.Order_uid] = model
	c.Unlock()
}

// GetModelCache - Получение одной модели из кеша.
func (c *Cache) GetModelCache(order_uid string) (models.Model, bool) {
	//Блокировка для новых читателей
	c.RLock()
	//Разблокировка читателей, после return
	defer c.RUnlock()
	//Уже заданный читатель берет данные из Кеша по order_uid
	model, flag := c.modelMap[order_uid]
	return model, flag
}

// GetModelsCache - Получение всех моделей
func (c *Cache) GetModelsCache() []models.Model {
	//создание слайлса из моделей по их количеству
	res := make([]models.Model, 0, len(c.modelMap))
	c.RLock()
	defer c.RUnlock()
	//добавление всех доступных моделей для чтения в слайс
	for _, res_model := range c.modelMap {
		res = append(res, res_model)
	}
	return res
}
