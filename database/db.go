package database

import (
	"os"
	"log"
	"database/sql"
	"fmt"
	"kafka-go-service/models"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	var err error

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
	dbHost, dbPort, dbUser, dbPassword, dbName)

	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	
	_, err = DB.Exec(`
	CREATE TABLE IF NOT EXISTS messages (
		id SERIAL PRIMARY KEY,
		content TEXT,
		processed BOOLEAN DEFAULT FALSE
		)
	`)

	if err != nil {
		log.Fatal(err)
	}
}

func StoreMessage(msg *models.Message) error {
	err := DB.QueryRow("INSERT INTO messages (content) VALUES ($1) RETURNING id", msg.Content).Scan(&msg.ID)
	return err
}

func UpdateMessageStatus(id int) error {
	_, err := DB.Exec("UPDATE messages SET processed = TRUE WHERE id = $1", id)
	return err
}

func GetStats() (models.Stats, error) {
	var stats models.Stats
	err := DB.QueryRow("SELECT COUNT(*), COUNT(CASE WHEN processed THEN 1 END) FROM messages").Scan(&stats.TotalMessages, &stats.ProcessedMessages)
	return stats, err
}
