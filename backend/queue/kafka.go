package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/CEhresmann/Container-Monitoring/config"
	"log"
	"time"

	"github.com/CEhresmann/Container-Monitoring/db"
	"github.com/segmentio/kafka-go"
)

var (
	brokerAddress = config.Cfg.Server.Port
	topic         = "ip_status_topic"
	partition     = 0
)

type IPStatus struct {
	IP       string `json:"ip"`
	PingTime int    `json:"ping_time"`
	LastOK   string `json:"last_ok"`
}

func ProduceMessage() {
	conn, err := kafka.DialLeader(context.Background(), "tcp", brokerAddress, topic, partition)
	if err != nil {
		log.Fatal("Ошибка подключения к Kafka:", err)
	}
	defer conn.Close()

	rows, err := db.DB.Query("SELECT ip, ping_time, last_ok FROM ips")
	if err != nil {
		log.Fatal("Ошибка запроса данных из БД:", err)
	}
	defer rows.Close()

	for rows.Next() {
		var ipStatus IPStatus
		err = rows.Scan(&ipStatus.IP, &ipStatus.PingTime, &ipStatus.LastOK)
		if err != nil {
			log.Println("Ошибка сканирования строки:", err)
			continue
		}

		message, err := json.Marshal(ipStatus)
		if err != nil {
			log.Println("Ошибка сериализации JSON:", err)
			continue
		}

		err = conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
		if err != nil {
			log.Println("Ошибка записи с дедлайном", err)
			continue
		}
		_, err = conn.WriteMessages(kafka.Message{Value: message})
		if err != nil {
			log.Println("Ошибка записи в Kafka:", err)
		}
	}

	fmt.Println("Сообщения отправлены в Kafka")
}

func ConsumeMessages() {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{brokerAddress},
		Topic:    topic,
		GroupID:  "ip_status_group",
		MaxBytes: 10e6, // 10MB
	})

	for {
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
