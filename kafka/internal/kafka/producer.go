package kafka

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(brokers []string) *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:                   kafka.TCP(brokers...),
			Balancer:               &kafka.LeastBytes{},
			AllowAutoTopicCreation: true,
		},
	}
}

func (p *Producer) Publish(ctx context.Context, topic string, key string, value interface{}) error {
	valueJSON, err := json.Marshal(value)
	if err != nil {
		return err
	}

	msg := kafka.Message{
		Topic: topic,
		Key:   []byte(key),
		Value: valueJSON,
	}

	err = p.writer.WriteMessages(ctx, msg)
	if err != nil {
		log.Printf("[producer] failed to write message to topic %s: %v", topic, err)
		return err
	}

	log.Printf("[producer] published message to topic %s, key=%s", topic, key)
	return nil
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
