package main

import (
	"log"

	"github.com/mxV03/warhousemanagementsystem/internal/storage"
	"github.com/mxV03/warhousemanagementsystem/internal/storage/sqlite"
)

func main() {
	// init db
	client, err := sqlite.InitDB()
	if err != nil {
		log.Fatalf("failed initializing database: %v", err)
	}
	defer client.Close()

	storage.SetClient(client)

	log.Println("Database initialized successfully")
}
