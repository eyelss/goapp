package main

import (
	"context"
	userpb "goapp/gen/goapp/user"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type userServer struct {
	userpb.UnimplementedUserServer
}

func (s *userServer) Check(ctx context.Context, in *userpb.UserRequest) (*userpb.UserResponse, error) {
	log.Printf("Received: %v", in.GetReq())

	return &userpb.UserResponse{Res: "OK: " + in.GetReq()}, nil
}

func main() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()

	userpb.RegisterUserServer(server, &userServer{})

	reflection.Register(server)

	log.Printf("server listening at %v", listener.Addr())
	if err := server.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
