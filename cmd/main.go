package main

import (
	"github.com/golang-migrate/migrate/v4"
	"github.com/svilenkomitov/eval/internal/server"
	"github.com/svilenkomitov/eval/internal/storage"
	"github.com/svilenkomitov/eval/internal/storage/migration"
	"log"
)

func main() {
	dbConfig := storage.LoadConfig()
	db, err := storage.Connect(dbConfig)
	if err != nil {
		log.Fatalf("connecting to database failed: %v", err)
	}

	migrationService := migration.New(dbConfig)
	if err := migrationService.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("db migration failed: %v", err)
	}

	serverConfig := server.LoadConfig()
	server := server.New(serverConfig, db)
	if err := server.Start(); err != nil {
		log.Fatalf("server failed to start on port %d: %v", serverConfig.Port, err)
	}
}
