package db

import (
	"database/sql"
	"embed"
	"fmt"
	"log"
	"net/http"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/httpfs"
	_ "github.com/mattn/go-sqlite3"
)

func Migrate(db *sql.DB, fs embed.FS) error {
	src, err := httpfs.New(http.FS(fs), "db/migration")
	if err != nil {
		return fmt.Errorf("invalid source instance, %v", err)
	}

	trg, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		return fmt.Errorf("creating sqlite3 db driver failed %v", err)
	}
	m, err := migrate.NewWithInstance("httpfs", src, "sqlite3", trg)
	if err != nil {
		return fmt.Errorf("failed to initialize migrate instance, %v", err)
	}

	err = m.Up()
	if err != nil {
		if err == migrate.ErrNoChange {
			log.Printf("[INFO] database schema is up to date")
			return nil
		}
		return fmt.Errorf("migrating database failed %v", err)
	}

	log.Printf("[INFO] database migrated to latest schema")
	return nil
}
