package main

import (
	"context"
	"goapp/framework/registry"
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
	consulRegistry, err := registry.NewConsulRegistry("localhost:8500")

	if err != nil {
		log.Fatal(err)
	}

	srv, err := server.New(
		server.WithName("user"),
		server.WithAddr("localhost:50051"),
		server.WithRegistry(consulRegistry),
	)

	if err != nil {
		log.Fatal(err)
	}

	userpb.RegisterUserServer(srv, &userServer{})

	server.Run(srv)
}
