package main

import (
	"fmt"
	"log"

	"github.com/mxV03/wms/internal/features/interfaces/cli"
	"github.com/mxV03/wms/internal/products"
	_ "github.com/mxV03/wms/internal/products/constraints"

	"github.com/mxV03/wms/internal/storage"
	"github.com/mxV03/wms/internal/storage/sqlite"
)

func main() {
	// init db
	client, err := sqlite.InitDB()
	if err != nil {
		log.Fatalf("failed initializing database: %v", err)
	}
	defer client.Close()

	storage.SetClient(client)

	fmt.Printf("Starting product: %s (tags: %v)\n", products.Name, products.EnabledTags)

	if err := cli.Run(); err != nil {
		log.Fatalf("CLI error: %v", err)
	}
}
