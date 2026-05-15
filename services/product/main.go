package main

import (
	"context"
	"goapp/framework/server"
	"goapp/framework/server/grpc"
	productpb "goapp/gen/goapp/product"
	"log"
)

type productServer struct {
	productpb.UnimplementedProductServer
}

func (s *productServer) Check(ctx context.Context, in *productpb.ProductRequest) (*productpb.ProductResponse, error) {
	log.Printf("Received: %v (Product)", in.GetReq())

	return &productpb.ProductResponse{Res: "OK: " + in.GetReq()}, nil
}

func main() {
	srv, err := grpc.New()

	if err != nil {
		log.Fatal(err)
	}

	productpb.RegisterProductServer(srv, &productServer{})

	server.Run(srv)

	select {}
}
