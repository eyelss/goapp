package framework

import (
	"fmt"
	"log"
	"net"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

func Load() (listener net.Listener, server *grpc.Server) {
	loadConfig()

	port := viper.GetInt("port")

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server = grpc.NewServer()

	if viper.GetBool("debug") {
		reflection.Register(server)
	}

	return
}

func Connect(serviceName string) (conn *grpc.ClientConn, err error) {
	conn, err = grpc.NewClient(
		fmt.Sprintf("%s:%d", serviceName, viper.GetInt("port")),
		grpc.WithTransportCredentials(insecure.NewCredentials()))

	return
}
