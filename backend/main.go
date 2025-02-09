package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/CEhresmann/Container-Monitoring/config"
	"github.com/CEhresmann/Container-Monitoring/db"
	"github.com/CEhresmann/Container-Monitoring/my_handler"
	"github.com/CEhresmann/Container-Monitoring/queue"
	"github.com/gorilla/handlers"
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
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	zap.ReplaceGlobals(logger)
	logger.Info("Запуск приложения")

	config.LoadConfig()
	db.InitDB()

	prometheus.MustRegister(requests)

	var wg sync.WaitGroup
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.Info("Запуск Kafka-потребителя")
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

	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.Info("Сервер запущен", zap.String("port", config.Cfg.Server.Port))
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("Ошибка запуска сервера", zap.Error(err))
		}
	}()

	<-stopChan
	logger.Info("Получен сигнал остановки, завершаем работу...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Ошибка при остановке сервера", zap.Error(err))
	}

	wg.Wait()
	logger.Info("Все горутины завершены, выход.")
}
