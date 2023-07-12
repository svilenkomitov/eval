package main

import (
	"github.com/svilenkomitov/eval/internal/server"
	"log"
)

func main() {
	serverConfig := server.LoadConfig()
	server := server.New(serverConfig)
	if err := server.Start(); err != nil {
		log.Fatalf("server failed to start on port %d: %v", serverConfig.Port, err)
	}
}
