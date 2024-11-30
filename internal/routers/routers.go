package routers

import (
	config "WBTECH/configs"
	"WBTECH/internal/cache"
	"WBTECH/internal/kafka"
	"WBTECH/internal/models"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	kfk "github.com/segmentio/kafka-go"
	"gorm.io/gorm"
)

func SetupRouter(c *cache.Cache, db *gorm.DB) *gin.Engine {
	r := gin.Default()

	r.GET("/orders/:order_uid", func(ctx *gin.Context) {
		orderUID := ctx.Param("order_uid")
		log.Println("Fetching order with orderUID:", orderUID)

		val, found := c.Get(orderUID)
		if !found {
			log.Println("Order not found in cache for orderUID:", orderUID)
			ctx.JSON(http.StatusNotFound, gin.H{"message": "Order not found"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"order": val})
	})

	r.POST("/orders", func(ctx *gin.Context) {
		var order models.Order
		if err := ctx.ShouldBindJSON(&order); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}

		jsonData, err := json.Marshal(order)
		if err != nil {
			log.Println("Error marshalling JSON:", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process data"})
			return
		}

		kafkaHost, err := config.GetKafkaURL()

		if err != nil {
			log.Println("Error kafka url")
		}

		writer := kfk.NewWriter(kfk.WriterConfig{
			Brokers: []string{kafkaHost},
			Topic:   "orders",
		})
		defer writer.Close()

		err = kafka.ProduceKafkaMessage(writer, order.OrderUID, jsonData)
		if err != nil {
			log.Println("Error writing to Kafka:", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send data to Kafka"})
			return
		}

		log.Println("Message sent to Kafka successfully")

		if err := c.AddOrder(db, order.OrderUID, string(jsonData)); err != nil {
			log.Println("Error saving data to DB or Cache:", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save data to DB or Cache"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Order processed and saved successfully"})
	})

	return r

}
