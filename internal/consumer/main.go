package main

import (
	"database/sql"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	_ "github.com/lib/pq"
	"log"
	"os"
	"os/signal"
)

const (
	kafkaBrokers  = "localhost:9092"
	kafkaTopic    = "my_topic"
	psqlHost      = "localhost"
	psqlPort      = 5432
	psqlUser      = "your_username"
	psqlPassword  = "your_password"
	psqlDBName    = "your_db_name"
	psqlTableName = "messages"
)

func main() {
	// Инициализация соединения с Kafka
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": kafkaBrokers,
		"group.id":          "my_consumer_group",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		log.Fatalf("Failed to create Kafka consumer: %v", err)
	}
	defer consumer.Close()

	// Подписка на топик Kafka
	err = consumer.SubscribeTopics([]string{kafkaTopic}, nil)
	if err != nil {
		log.Fatalf("Failed to subscribe to Kafka topic: %v", err)
	}

	// Инициализация соединения с PostgreSQL
	psqlConnStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		psqlHost, psqlPort, psqlUser, psqlPassword, psqlDBName)
	db, err := sql.Open("postgres", psqlConnStr)
	if err != nil {
		log.Fatalf("Failed to open PostgreSQL connection: %v", err)
	}
	defer db.Close()

	// Создание таблицы в PostgreSQL, если она не существует
	if err := createTableIfNotExists(db); err != nil {
		log.Fatalf("Failed to create table in PostgreSQL: %v", err)
	}

	// Обработка сообщений из Kafka и запись в PostgreSQL
	run := true
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	for run {
		select {
		case sig := <-signals:
			log.Printf("Received signal: %v", sig)
			run = false
		default:
			ev := consumer.Poll(100)
			if ev == nil {
				continue
			}

			switch e := ev.(type) {
			case *kafka.Message:
				message := string(e.Value)
				log.Printf("Received message: %s", message)

				// Запись сообщения в PostgreSQL
				if err := insertMessageToPostgreSQL(db, message); err != nil {
					log.Printf("Failed to insert message to PostgreSQL: %v", err)
				}

			case kafka.Error:
				log.Printf("Received Kafka error: %v", e)

			default:
				log.Printf("Ignoring event: %v", e)
			}
		}
	}

	// Отписка от топика и остановка потребителя Kafka
	consumer.Unsubscribe()
	consumer.Close()
}

func createTableIfNotExists(db *sql.DB) error {
	query := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			id SERIAL PRIMARY KEY,
			message TEXT
		)`, psqlTableName)

	_, err := db.Exec(query)
	return err
}

func insertMessageToPostgreSQL(db *sql.DB, message string) error {
	query := fmt.Sprintf("INSERT INTO %s (message) VALUES ($1)", psqlTableName)
	_, err := db.Exec(query, message)
	return err
}
