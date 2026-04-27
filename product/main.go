package main

import (
	"context"
	productpb "goapp/gen/goapp/product"
	userpb "goapp/gen/goapp/user"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

type productServer struct {
	productpb.UnimplementedProductServer
}

func (s *productServer) Check(ctx context.Context, in *productpb.ProductRequest) (*productpb.ProductResponse, error) {
	log.Printf("Received: %v (Product)", in.GetReq())

	conn, err := grpc.NewClient(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("could not connect users service: %v", err)
	}
	defer conn.Close()
	c := userpb.NewUserClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.Check(ctx, &userpb.UserRequest{Req: in.GetReq()})
	if err != nil {
		log.Fatalf("could not check user: %v", err)
	}

	log.Printf("Check result: %v", r.Res)
	return &productpb.ProductResponse{Res: "OK: " + in.GetReq()}, nil
}

func main() {
	listener, err := net.Listen("tcp", ":50050")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()

	productpb.RegisterProductServer(server, &productServer{})

	reflection.Register(server)

	log.Printf("server listening at %v", listener.Addr())
	if err := server.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
