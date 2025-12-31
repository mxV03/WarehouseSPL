package sqlite

import (
	"context"
	"database/sql"
	"log"
	"os"
	"path/filepath"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/mxV03/wms/ent"
	_ "modernc.org/sqlite"
)

func InitDB() (*ent.Client, error) {
	dir := "data"
	if err := os.MkdirAll(dir, 0o755); err != nil {
		log.Fatalf("failed creating db dir: %v", err)
	}

	dbPath := filepath.Join(dir, "warehouse.db")
	dsn := "file:" + dbPath + "?_fk=1"

	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		log.Fatalf("failed opening connection to sqlite: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("failed pinging sqlite: %v", err)
	}

	if _, err := db.Exec(`PRAGMA foreign_keys = ON`); err != nil {
		log.Fatalf("failed enabling foreign keys: %v", err)
	}

	drv := entsql.OpenDB(dialect.SQLite, db)

	client := ent.NewClient(ent.Driver(drv))

	// run auto migration
	if err := client.Schema.Create(context.Background()); err != nil {
		client.Close()
		log.Fatalf("failed creating schema resources: %v", err)
	}

	return client, nil
}
