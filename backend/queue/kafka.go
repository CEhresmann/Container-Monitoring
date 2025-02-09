package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/CEhresmann/Container-Monitoring/config"
	"log"
	"os"

	"github.com/CEhresmann/Container-Monitoring/db"
	"github.com/segmentio/kafka-go"
)

/*
var (

	brokerAddress = config.Cfg.Queue.Broker

)
*/
type IPStatus struct {
	IP       string `json:"ip"`
	PingTime int    `json:"ping_time"`
	LastOK   string `json:"last_ok"`
}

func ConsumeMessages() {
	broker := os.Getenv("QUEUE_BROKER")
	if broker == "" {
		broker = config.Cfg.Queue.Broker 
	}

	topic := os.Getenv("PING_TOPIC")
	if topic == "" {
		topic = config.Cfg.Queue.Topic
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{broker},
		Topic:    topic,
		GroupID:  "AwyYJhkoTAmpeLAJ_BMgHg",
		MaxBytes: 10e6, // 10MB
	})

	for {
		log.Printf("Broker: [%s], Topic: [%s]", broker, topic)
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Println("Ошибка чтения сообщения из Kafka:", err)
			continue
		}

		var ipStatus IPStatus
		err = json.Unmarshal(m.Value, &ipStatus)
		if err != nil {
			log.Println("Ошибка десериализации JSON:", err)
			continue
		}

		_, err = db.DB.Exec("INSERT INTO ips (ip, ping_time, last_ok) VALUES ($1, $2, $3)",
			ipStatus.IP, ipStatus.PingTime, ipStatus.LastOK)
		if err != nil {
			log.Println("Ошибка вставки в БД:", err)
			continue
		}

		fmt.Println("Данные добавлены в БД:", ipStatus)
	}
}
