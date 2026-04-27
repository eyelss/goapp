package main

import (
	"context"
	"goapp/framework"
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
	listener, server := framework.Load()

	userpb.RegisterUserServer(server, &userServer{})

	log.Printf("server listening at %v", listener.Addr())
	if err := server.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
