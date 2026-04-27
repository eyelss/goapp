package main

import (
	"context"
	"goapp/framework"
	productpb "goapp/gen/goapp/product"
	userpb "goapp/gen/goapp/user"
	"log"
	"time"
)

type productServer struct {
	productpb.UnimplementedProductServer
}

func (s *productServer) Check(ctx context.Context, in *productpb.ProductRequest) (*productpb.ProductResponse, error) {
	log.Printf("Received: %v (Product)", in.GetReq())

	conn, err := framework.Connect("user")

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
	listener, server := framework.Load()

	productpb.RegisterProductServer(server, &productServer{})

	log.Printf("server listening at %v", listener.Addr())
	if err := server.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
