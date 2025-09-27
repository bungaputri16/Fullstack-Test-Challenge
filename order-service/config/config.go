package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type Config struct {
	DBHost             string
	DBPort             string
	DBUser             string
	DBPassword         string
	DBName             string
	RedisHost          string
	RedisPort          string
	RabbitMQHost       string
	RabbitMQUser       string
	RabbitMQPass       string
	ServicePort        string
	ProductServiceHost string
	ProductServicePort string
}

func Load() Config {
	return Config{
		DBHost:             getEnv("POSTGRES_HOST", "localhost"),
		DBPort:             getEnv("POSTGRES_PORT", "5432"),
		DBUser:             getEnv("POSTGRES_USER", "user"),
		DBPassword:         getEnv("POSTGRES_PASSWORD", "password"),
		DBName:             getEnv("POSTGRES_DB", "appdb"),
		RedisHost:          getEnv("REDIS_HOST", "localhost"),
		RedisPort:          getEnv("REDIS_PORT", "6379"),
		RabbitMQHost:       getEnv("RABBITMQ_HOST", "localhost"),
		RabbitMQUser:       getEnv("RABBITMQ_USER", "admin"),
		RabbitMQPass:       getEnv("RABBITMQ_PASS", "secret"),
		ServicePort:        getEnv("ORDER_SERVICE_PORT", "4000"),
		ProductServiceHost: getEnv("PRODUCT_SERVICE_HOST", "localhost"),
		ProductServicePort: getEnv("PRODUCT_SERVICE_PORT", "3000"),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func InitDB(cfg Config) *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("DB not reachable:", err)
	}

	log.Println("Connected to PostgreSQL")
	return db
}
