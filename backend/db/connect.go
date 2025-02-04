package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/CEhresmann/Container-Monitoring/config"
	_ "github.com/lib/pq"
)

var DB *sql.DB

const (
	maxRetries    = 6
	retryInterval = 1 * time.Second
)

func InitDB() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Cfg.Database.Host,
		config.Cfg.Database.User,
		config.Cfg.Database.Password,
		config.Cfg.Database.DBName,
	)

	var err error
	for i := 0; i < maxRetries; i++ {
		DB, err = sql.Open("postgres", dsn)
		if err == nil {
			err = DB.Ping()
			if err == nil {
				log.Println("Подключение к базе данных успешно")
				CreateTable()
				return
			}
		}
		log.Printf("Ошибка подключения к БД (попытка %d/%d): %v", i+1, maxRetries, err)
		time.Sleep(retryInterval)
	}

	log.Fatal("Не удалось подключиться к базе данных после нескольких попыток:", err)
}

func CreateTable() {
	executeSQLFile("db/Create.sql", "Ошибка создания таблицы")
}

func executeSQLFile(filepath string, errorMessage string) {
	sqlScript, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatalf("Ошибка чтения SQL-файла %s: %v", filepath, err)
	}

	_, err = DB.Exec(string(sqlScript))
	if err != nil {
		log.Printf("%s: %v", errorMessage, err)
	}

	log.Printf("SQL-скрипт %s успешно выполнен.", filepath)
}
