package main

import (
	"log"

	"codecrafters-http-server-go/app/internal/config"
	"codecrafters-http-server-go/app/internal/server"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Create and start the server
	srv := server.New(cfg)

	log.Printf("Starting HTTP server on %s", cfg.Address)
	if err := srv.Start(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
