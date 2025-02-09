package queue

import (
	"context"
	"encoding/json"
	"os"

	"go.uber.org/zap"

	"github.com/CEhresmann/Container-Monitoring/config"
	"github.com/CEhresmann/Container-Monitoring/db"
	"github.com/segmentio/kafka-go"
)

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

	zap.L().Info("Запуск Kafka-потребителя", zap.String("broker", broker), zap.String("topic", topic))

	for {
		zap.L().Info("Чтение сообщения из Kafka", zap.String("broker", broker), zap.String("topic", topic))
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			zap.L().Error("Ошибка чтения сообщения из Kafka", zap.Error(err))
			continue
		}

		var ipStatus IPStatus
		if err := json.Unmarshal(m.Value, &ipStatus); err != nil {
			zap.L().Error("Ошибка десериализации JSON", zap.Error(err))
			continue
		}

		_, err = db.DB.Exec("INSERT INTO ips (ip, ping_time, last_ok) VALUES ($1, $2, $3)",
			ipStatus.IP, ipStatus.PingTime, ipStatus.LastOK)
		if err != nil {
			zap.L().Error("Ошибка вставки в БД", zap.Error(err))
			continue
		}

		zap.L().Info("Данные успешно добавлены в БД", zap.Any("ipStatus", ipStatus))
	}
}
