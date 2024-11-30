package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

func GetKafkaURL() (string, error) {
	kafkaHost := os.Getenv("KAFKA_HOST")
	kafkaPort := os.Getenv("KAFKA_PORT")

	if kafkaHost == "" || kafkaPort == "" {
		return "", fmt.Errorf("one or more required environment variables (KAFKA_HOST, KAFKA_PORT) are missing")
	}

	kafkaURL := fmt.Sprintf("%s:%s", kafkaHost, kafkaPort)
	return kafkaURL, nil
}

func GetDBURL() (string, error) {
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	name := os.Getenv("DB_NAME")
	sslMode := os.Getenv("DB_SSLMODE")

	if user == "" || password == "" || host == "" || port == "" || name == "" {
		return "", fmt.Errorf("one or more required environment variables are missing")
	}

	if sslMode == "" {
		sslMode = "disable"
	}

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", user, password, host, port, name, sslMode)
	return dbURL, nil
}
