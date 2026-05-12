package main

import (
	"context"
	"goapp/framework/lib/kafka"
	"goapp/framework/server"
	userpb "goapp/gen/goapp/user"
	"log"
	"time"
)

type userServer struct {
	userpb.UnimplementedUserServer
}

func (s *userServer) Check(ctx context.Context, in *userpb.UserRequest) (*userpb.UserResponse, error) {
	log.Printf("Received: %v (User)", in.GetReq())

	return &userpb.UserResponse{Res: "OK: " + in.GetReq()}, nil
}

func main() {
	srv, err := server.New()

	if err != nil {
		log.Fatal(err)
	}

	userpb.RegisterUserServer(srv, &userServer{})

	server.Run(srv)

	kafka.Listen(kafka.Basic, func(message kafka.Message) {
		log.Printf("GOT MESSAGE: %v\n", message.Topic)
		log.Printf("%s => %s", string(message.Key), string(message.Value))
	})

	time.Sleep(time.Second * 15)

	log.Println("Publishing message!")
	errp := kafka.Publish(kafka.Basic, kafka.Message{
		Key:   []byte("some-key"),
		Value: []byte("some-value"),
	})

	if errp != nil {
		log.Fatal(errp)
	}

	select {}
}
