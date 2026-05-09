package server

import (
	"context"
	"log"
	"os"
)

func Run(server *Server) {
	ctx := context.Background()

	if err := server.Start(ctx); err != nil {
		log.Printf("failed to start server: %v", err)

		os.Exit(1)
	}
}
