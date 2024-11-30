// internal/kafka/kafka.go
package kafka

import (
	"context"
	"encoding/json"
	"log"

	"WBTECH/internal/cache"
	"WBTECH/internal/models"

	"github.com/segmentio/kafka-go"
	"gorm.io/gorm"
)

func ProduceKafkaMessage(writer *kafka.Writer, key string, value []byte) error {
	err := writer.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte(key),
		Value: value,
	})
	if err != nil {
		log.Printf("Failed to send Kafka message: %v", err)
		return err
	}

	log.Printf("Message sent to Kafka: key=%s, value=%s", key, string(value))
	return nil
}

func ConsumeKafkaMessages(broker, topic string, db *gorm.DB, cache *cache.Cache) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{broker},
		Topic:   topic,
		GroupID: "order_service_group",
	})
	defer r.Close()

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Println("Error reading message:", err)
			continue
		}

		log.Printf("Message received: %s", string(m.Value))

		// Сохраняем данные в кэш
		cache.AddOrder(db, string(m.Key), string(m.Value))

		// Парсим JSON в структуру Order и сохраняем в БД
		var order models.Order
		if err := json.Unmarshal(m.Value, &order); err != nil {
			log.Println("Error unmarshalling order:", err)
			continue
		}

		// Сохраняем заказ в базу данных через транзакцию
		err = saveOrderToDB(db, &order)
		if err != nil {
			log.Println("Error saving order to DB:", err)
		} else {
			log.Printf("Order saved to DB: %s", string(rune(order.ID)))
		}
	}
}

// Функция для сохранения заказа в базу данных с транзакцией
func saveOrderToDB(db *gorm.DB, order *models.Order) error {
	return db.Transaction(func(tx *gorm.DB) error {
		// Сохраняем основной заказ
		if err := tx.Create(order).Error; err != nil {
			return err
		}

		// Сохраняем связанные данные
		if err := tx.Create(&order.Delivery).Error; err != nil {
			return err
		}
		if err := tx.Create(&order.Payment).Error; err != nil {
			return err
		}
		if len(order.Items) > 0 {
			if err := tx.Create(&order.Items).Error; err != nil {
				return err
			}
		}

		return nil
	})
}
