package db

import (
	"database/sql"
	"fmt"
	"github.com/CEhresmann/Container-Monitoring/config"
	_ "github.com/lib/pq"
	"log"
	"os"
)

var DB *sql.DB

func InitDB() {
	dsn := "host=" + config.Cfg.Database.Host + "user=" + config.Cfg.Database.User +
		"password=" + config.Cfg.Database.Password + "dbname=" + config.Cfg.Database.DBName +
		"sslmode=disable"
	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer DB.Close()
}

func CreateUser() {
	sqlScript, err := os.ReadFile("db/User.sql")
	if err != nil {
		log.Println("failed to read SQL script: ", err)
	}
	createTableSQL := sqlScript
	_, err = DB.Exec(string(createTableSQL))
	if err != nil {
		log.Fatalf("Failed to create user: %v", err)
	}
	fmt.Println("User created successfully.")
}

func CreateTable() {
	sqlScript, err := os.ReadFile("db/Create.sql")
	if err != nil {
		log.Println("failed to read SQL script: ", err)
	}
	createTableSQL := sqlScript
	_, err = DB.Exec(string(createTableSQL))
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}
	fmt.Println("Table created successfully.")
}
