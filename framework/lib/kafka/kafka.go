package kafka

import (
	"context"
	config "goapp/framework/lib"
	"log"

	"github.com/segmentio/kafka-go"
)

func getConnection(topic string, partition int) (*kafka.Conn, error) {
	kafkaAddr := config.Get[string]("kafka.address")

	return kafka.DialLeader(context.Background(), "tcp", kafkaAddr, topic, partition)
}

func Write(topic string) {
	ctx := context.Background()

	kafkaAddress := config.Get[string]("kafka-addr")
	writer := kafka.Writer{
		Addr:                   kafka.TCP(kafkaAddress),
		Topic:                  topic,
		Balancer:               &kafka.LeastBytes{},
		AllowAutoTopicCreation: true,
	}

	err := writer.WriteMessages(ctx,
		kafka.Message{
			Key:   []byte("Key-A"),
			Value: []byte("Value-A"),
		},
	)

	if err != nil {
		log.Fatal(err)
	}

	log.Print("Kafka wrote")
}
