package config

import (
	"flag"
)

// Config holds the application configuration
type Config struct {
	// FilesDirectory is the directory to serve files from
	FilesDirectory string

	// Address is the server's listening address
	Address string
}

// Load reads configuration from various sources and returns a config struct
func Load() *Config {
	cfg := &Config{
		Address: "0.0.0.0:4221",
	}

	// Parse command line flags
	flag.StringVar(&cfg.FilesDirectory, "directory", "", "Directory to serve files from")
	flag.Parse()

	return cfg
}
