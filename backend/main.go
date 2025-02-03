package main

import (
	"github.com/CEhresmann/Container-Monitoring/config"
	"github.com/CEhresmann/Container-Monitoring/db"
	"github.com/CEhresmann/Container-Monitoring/handlers"
	"github.com/CEhresmann/Container-Monitoring/queue"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

func main() {
	config.LoadConfig()

	db.InitDB()
	db.CreateUser()
	db.CreateTable()

	go queue.ConsumeMessages()

	go func() {
		for {
			queue.ProduceMessage()
			time.Sleep(15 * time.Second)
		}
	}()

	router := mux.NewRouter()
	router.HandleFunc("/api/ip", handlers.GetIPStatuses).Methods("GET")
	router.HandleFunc("/api/ip", handlers.AddIPStatus).Methods("POST")

	log.Println("Server running on port", config.Cfg.Server.Port)
	log.Fatal(http.ListenAndServe(":"+config.Cfg.Server.Port, router))
}
