package main

import (
	"github.com/CEhresmann/Container-Monitoring/Pinger/ping"
	"log"
	"os"
)

func main() {
	broker := os.Getenv("QUEUE_BROKER")
	topic := os.Getenv("PING_TOPIC")

	log.Println("Запуск Pinger...")

	if broker == "" || topic == "" {
		log.Fatal("Ошибка: QUEUE_BROKER и PING_TOPIC должны быть заданы в переменных окружения")
	}

	ping.RunPinger(broker, topic)
}
