package main

import (
	"context"
	"errors"
	"github.com/gorilla/handlers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/CEhresmann/Container-Monitoring/config"
	"github.com/CEhresmann/Container-Monitoring/db"
	"github.com/CEhresmann/Container-Monitoring/my_handler"
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
	corsHandler := handlers.CORS()(router)

	router.HandleFunc("/api/ip", my_handler.GetIPStatuses).Methods("GET")
	router.HandleFunc("/api/ip", my_handler.AddIPStatus).Methods("POST")
	router.Handle("/metrics", promhttp.Handler())

	server := &http.Server{
		Addr:    ":" + config.Cfg.Server.Port,
		Handler: corsHandler,
	}
	go func() {
		prometheus.MustRegister(requests)
	}()
	go func() {
		defer wg.Done()
		log.Println("Сервер запущен на:", config.Cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Ошибка запуска сервера: %v", err)
		}
	}()

	<-stopChan
	log.Println("Получен сигнал остановки, завершаем работу...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Ошибка при остановке сервера: %v", err)
	}

	wg.Wait()
	log.Println("Все горутины завершены, выход.")
}
