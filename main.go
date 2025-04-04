package main

import (
	"context"
	"crypto/tls"
	"log"
	"os"

	"github.com/IBM/sarama"
	"github.com/aws/aws-lambda-go/lambda"
)

type Event struct {
	Topic string `json:"topic"`
	Key   string `json:"key"`
	Value string `json:"value"`
}

func saramaConfig() *sarama.Config {
	config := sarama.NewConfig()

	// Kafka version
	config.Version = sarama.V3_9_0_0 // Adjust according to your Kafka version

	// SASL/SCRAM Authentication
	config.Net.SASL.Enable = true
	config.Net.SASL.Handshake = true
	config.Net.SASL.User = os.Getenv("KAFKA_USERNAME")
	config.Net.SASL.Password = os.Getenv("KAFKA_PASSWORD")
	config.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA256

	// TLS (optional but recommended)
	config.Net.TLS.Enable = true
	config.Net.TLS.Config = &tls.Config{
		InsecureSkipVerify: true, // Change to false in production
	}

	// Producer settings
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	return config
}

// KafkaProducer function to send messages to Kafka
func KafkaProducer(ctx context.Context, event Event) error {
	kafkaBrokers := []string{os.Getenv("KAFKA_BROKER")}
	config := saramaConfig()
	producer, err := sarama.NewSyncProducer(kafkaBrokers, config)
	if err != nil {
		log.Printf("Failed to start Sarama producer: %v", err)
		return err
	}
	defer producer.Close()

	msg := &sarama.ProducerMessage{
		Topic: event.Topic,
		Key:   sarama.StringEncoder(event.Key),
		Value: sarama.StringEncoder(event.Value),
	}

	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		log.Printf("Failed to send message: %v", err)
		return err
	}

	log.Printf("Message sent to partition %d at offset %d", partition, offset)
	return nil
}

func main() {
	lambda.Start(KafkaProducer)
}
