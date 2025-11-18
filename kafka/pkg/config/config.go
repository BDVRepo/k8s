package config

import "os"

type KafkaConfig struct {
	Brokers []string
}

func LoadKafka() KafkaConfig {
	brokers := os.Getenv("KAFKA_BROKERS")
	if brokers == "" {
		brokers = "localhost:9092"
	}

	return KafkaConfig{
		Brokers: []string{brokers},
	}
}

