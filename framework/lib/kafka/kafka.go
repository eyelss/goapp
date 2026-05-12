package kafka

import (
	"context"
	config "goapp/framework/lib"
	"log"

	"github.com/segmentio/kafka-go"
)

type Message = kafka.Message

type EnqMessage[T any] struct {
	Key   []byte
	Value T
}

func Publish(topic Topic, messages ...Message) error {
	ctx := context.Background()

	kafkaAddress := config.Get[string]("kafka-addr")
	writer := kafka.Writer{
		Addr:     kafka.TCP(kafkaAddress),
		Topic:    string(topic),
		Balancer: &kafka.LeastBytes{},
	}

	ProtoSerializer{}.Encode(messages)

	err := writer.WriteMessages(ctx, messages...)

	return err
}

func Listen(topic Topic, handler func(Message)) {
	kafkaAddress := config.Get[string]("kafka-addr")
	appName := config.Get[string]("app-name")

	reader := kafka.NewReader(kafka.ReaderConfig{
		GroupID:  appName,
		Brokers:  []string{kafkaAddress},
		Topic:    string(topic),
		MaxBytes: 10e6,
	})

	go func() {
		for {
			m, err := reader.ReadMessage(context.Background())

			if err != nil {
				log.Printf("error reading kafka message: %s\n", err)
			} else {
				handler(m)
			}
		}
	}()
}
