package main

import (
	"flag"
	"log"

	"github.com/codecrafters-io/http-server-starter-go/app/internal/config"
	"github.com/codecrafters-io/http-server-starter-go/app/internal/server"
)

func main() {
	// Parse command line flags
	directory := flag.String("directory", "", "Directory to serve files from")
	flag.Parse()

	// Create config
	cfg := &config.Config{
		FilesDirectory: *directory,
		Address:        "0.0.0.0:4221",
	}

	// Create and start the server
	srv := server.New(cfg)

	log.Printf("Starting HTTP server on %s", cfg.Address)
	if err := srv.Start(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
