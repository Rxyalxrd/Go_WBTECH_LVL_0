package cache

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"WBTECH/internal/models"

	"github.com/patrickmn/go-cache"
	"gorm.io/gorm"
)

type Cache struct {
	c *cache.Cache
}

func NewCache() *Cache {
	return &Cache{
		c: cache.New(10*time.Minute, 15*time.Minute),
	}
}

func (cache *Cache) AddOrder(db *gorm.DB, orderUID string, jsonData string) error {
	log.Println("Adding to cache with key:", orderUID)
	cache.c.Set(orderUID, jsonData, 10*time.Minute)

	if _, found := cache.c.Get(orderUID); !found {
		log.Println("Failed to cache data for order:", orderUID)
	} else {
		log.Println("Order cached successfully:", orderUID)
	}

	var order models.Order
	err := json.Unmarshal([]byte(jsonData), &order)
	if err != nil {
		log.Println("Error unmarshalling JSON:", err)
		return err
	}

	if order.OrderUID == "" {
		return fmt.Errorf("invalid data: OrderUID is empty")
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ?", order.Delivery.ID).FirstOrCreate(&order.Delivery).Error; err != nil {
			return fmt.Errorf("failed to save Delivery: %w", err)
		}
		order.DeliveryID = order.Delivery.ID

		if err := tx.Where("id = ?", order.Payment.ID).FirstOrCreate(&order.Payment).Error; err != nil {
			return fmt.Errorf("failed to save Payment: %w", err)
		}
		order.PaymentID = order.Payment.ID

		var existingOrder models.Order
		if err := tx.First(&existingOrder, "order_uid = ?", order.OrderUID).Error; err == nil {
			log.Println("Order already exists:", order.OrderUID)
			return nil
		}
		if err := tx.Create(&order).Error; err != nil {
			return fmt.Errorf("failed to save Order: %w", err)
		}

		for i := range order.Items {
			order.Items[i].OrderID = order.ID
			order.Items[i].ID = 0
		}

		if len(order.Items) > 0 {
			if err := tx.Create(&order.Items).Error; err != nil {
				return fmt.Errorf("failed to save Items: %w", err)
			}
		}

		return nil
	})

	if err != nil {
		log.Println("Error saving order to DB:", err)
		return err
	}

	log.Println("Order saved successfully:", orderUID)
	return nil
}

func (cache *Cache) RestoreCache(db *gorm.DB) {
	var orders []models.Order

	if err := db.Preload("Delivery").Preload("Payment").Preload("Items").Find(&orders).Error; err != nil {
		log.Println("Error restoring cache:", err)
		return
	}

	for _, order := range orders {
		jsonData, err := json.Marshal(order)
		if err != nil {
			log.Println("Error marshalling order to JSON:", err)
			continue
		}
		cache.c.Set(order.OrderUID, string(jsonData), 10*time.Minute)
	}

	log.Println("Cache restored from DB")
}

func (c *Cache) Get(orderUID string) (string, bool) {
	log.Println("Fetching from cache with key:", orderUID)

	val, found := c.c.Get(orderUID)
	if !found {
		log.Println("Order not found in cache:", orderUID)
		return "", false
	}

	if strVal, ok := val.(string); ok {
		log.Println("Order found in cache:", orderUID)
		return strVal, true
	}

	log.Println("Invalid data type in cache for key:", orderUID)
	return "", false
}
