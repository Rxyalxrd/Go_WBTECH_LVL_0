package main

import (
	"log"

	"WBTECH/internal/cache"
	"WBTECH/internal/database"
	"WBTECH/internal/models"
	"WBTECH/internal/routers"
)

func main() {
	db, err := database.ConnectDB()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	database.InitDBAndMigrate(&models.Order{}, &models.Delivery{}, &models.Payment{}, &models.Item{})

	err = database.AutoMigrateTables(db, &models.Order{}, &models.Delivery{}, &models.Payment{}, &models.Item{})
	if err != nil {
		log.Fatalf("Error during migration: %v", err)
	}

	log.Println("DB instance ready and migrations applied.")

	myCache := cache.NewCache()

	myCache.RestoreCache(db)

	r := routers.SetupRouter(myCache, db)

	log.Println("Starting server on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
