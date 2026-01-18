package main

import (
	"log"

	server "github.com/yogusita/to-adhdo/server"
)

func main() {
	if err := server.Start(); err != nil {
		log.Printf("Failed to start server: %v", err)
	}
}
