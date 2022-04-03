package cache

import (
	"L0/back/models"
	"sync"
)

type Cache struct {
	sync.RWMutex
	modelMap map[string]models.Model
}

func CreateCache() *Cache {
	return &Cache{
		modelMap: make(map[string]models.Model),
	}
}

func (c *Cache) AddModelCache(model models.Model)   {
	c.Lock()
	c.modelMap[*model.Order_uid] = model
	c.Unlock()
}

func (c *Cache) GetModelCache(order_uid string) (models.Model, bool)  {
	c.RLock()
	defer c.RUnlock()
	model, flag := c.modelMap[order_uid]
	return model, flag
}

func (c *Cache) GetModelsCache() []models.Model {
	res := make([]models.Model, 0, len(c.modelMap))
	c.RLock()
	defer c.RUnlock()

	for _, res_model := range c.modelMap{
		res = append(res, res_model)
	}

	return res
}

