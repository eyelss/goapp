package server

import (
	"log"
	"os"
)

func Run(server IServer) {
	if err := server.Start(); err != nil {
		log.Printf("failed to start server: %v", err)

		os.Exit(1)
	}
}
