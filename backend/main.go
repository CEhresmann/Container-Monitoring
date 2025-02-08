package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/CEhresmann/Container-Monitoring/config"
	"github.com/CEhresmann/Container-Monitoring/db"
	"github.com/CEhresmann/Container-Monitoring/handlers"
	"github.com/CEhresmann/Container-Monitoring/queue"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	requests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "status"},
	)
)

func init() {
	prometheus.MustRegister(requests)
}

func main() {
	config.LoadConfig()

	db.InitDB()

	var wg sync.WaitGroup
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	wg.Add(1)
	go func() {
		defer wg.Done()
		queue.ConsumeMessages()
	}()

	router := mux.NewRouter()
	router.HandleFunc("/api/ip", handlers.GetIPStatuses).Methods("GET")
	router.HandleFunc("/api/ip", handlers.AddIPStatus).Methods("POST")

	router.Handle("/metrics", promhttp.Handler())

	server := &http.Server{
		Addr:    ":" + config.Cfg.Server.Port,
		Handler: router,
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Println("Сервер запущен на:", config.Cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Ошибка запуска сервера: %v", err)
		}
	}()

	<-stopChan
	log.Println("Получен сигнал остановки, завершаем работу...")

	if err := server.Close(); err != nil {
		log.Printf("Ошибка при закрытии сервера: %v", err)
	}

	wg.Wait()
	log.Println("Все горутины завершены, выход.")
}
