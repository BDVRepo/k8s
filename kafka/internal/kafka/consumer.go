package kafka

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader *kafka.Reader
}

func NewConsumer(brokers []string, topic string, groupID string) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:  brokers,
			Topic:    topic,
			GroupID:  groupID,
			MinBytes: 10e3, // 10KB
			MaxBytes: 10e6, // 10MB
		}),
	}
}

func (c *Consumer) ReadMessage(ctx context.Context) (kafka.Message, error) {
	return c.reader.ReadMessage(ctx)
}

func (c *Consumer) UnmarshalMessage(msg kafka.Message, v interface{}) error {
	return json.Unmarshal(msg.Value, v)
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}

func (c *Consumer) LogMessage(msg kafka.Message) {
	log.Printf("[consumer] received message: topic=%s partition=%d offset=%d key=%s",
		msg.Topic, msg.Partition, msg.Offset, string(msg.Key))
}

