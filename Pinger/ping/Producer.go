package ping

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"log"
)

func ProduceToKafka(ipStatus IPStatus, topic string, broker string) {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{broker},
		Topic:   topic,
		Async:   false,
	})
	defer writer.Close()

	message, err := json.Marshal(ipStatus)
	if err != nil {
		log.Println("Ошибка кодирования JSON:", err)
		return
	}

	err = writer.WriteMessages(context.Background(), kafka.Message{
		Value: message,
	})
	if err != nil {
		log.Println("Ошибка отправки в Kafka:", err)
	} else {
		log.Println("Данные отправлены в Kafka:", ipStatus)
	}
}
