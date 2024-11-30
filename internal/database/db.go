package database

import (
	config "WBTECH/configs"
	"flag"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ConnectDB подключается к базе данных
func ConnectDB() (*gorm.DB, error) {

	config.LoadEnv()

	dbURL, err := config.GetDBURL()
	if err != nil {
		return nil, fmt.Errorf("failed to get DB URL: %w", err)
	}

	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("Database connection established successfully!")
	return db, nil
}

// AutoMigrateTables выполняет миграции для указанных моделей
func AutoMigrateTables(db *gorm.DB, models ...interface{}) error {
	err := db.AutoMigrate(models...)
	if err != nil {
		return fmt.Errorf("error during migration: %v", err)
	}

	log.Println("Migration completed successfully!")
	return nil
}

// InitDBAndMigrate инициализирует базу данных и выполняет миграции при запуске с флагом -migrate
func InitDBAndMigrate(models ...interface{}) *gorm.DB {
	config.LoadEnv()

	dbURL, err := config.GetDBURL()
	if err != nil {
		log.Fatalf("Error getting DB URL: %v", err)
	}

	migrate := flag.Bool("migrate", false, "Run database migrations")
	flag.Parse()

	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if *migrate {
		err := AutoMigrateTables(db, models...)
		if err != nil {
			log.Fatalf("Error during migration: %v", err)
		}
		log.Println("Migration completed successfully. Exiting...")
		os.Exit(0)
	}

	return db
}
