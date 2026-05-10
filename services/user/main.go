package main

import (
	"context"
	"goapp/framework/lib/kafka"
	"goapp/framework/server"
	userpb "goapp/gen/goapp/user"
	"log"
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

	log.Println("Before")
	kafka.Write("some-topic")
	log.Println("After")

	select {}
}
