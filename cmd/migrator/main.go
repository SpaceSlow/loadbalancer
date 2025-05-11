package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"

	"github.com/SpaceSlow/loadbalancer/config"
)

func main() {
	cfg, err := config.LoadConfig(os.Getenv("CONFIG_PATH"))
	if err != nil {
		log.Fatal(fmt.Sprintf("Load config error: %s", err.Error()))
	}
	db, err := sql.Open(cfg.DB.DBMS, cfg.DB.DSN)
	if err != nil {
		log.Fatal(fmt.Sprintf("Open db error: %s", err.Error()))
	}

	var driver database.Driver
	switch cfg.DB.DBMS {
	case "postgres":
		driver, err = postgres.WithInstance(db, &postgres.Config{})
	default:
		log.Fatal("Unknown database management system")
	}
	if err != nil {
		log.Fatal(fmt.Sprintf("Initialize driver error: %s", err.Error()))
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file:///app/migrations",
		cfg.DB.DBMS,
		driver,
	)
	if err != nil {
		log.Fatal(fmt.Sprintf("Initialize database migrate instance error: %s", err.Error()))
	}

	err = m.Up()
	if err != nil {
		log.Fatal(fmt.Sprintf("Migrate up error: %s", err.Error()))
	}
	log.Print("Success migrate up")
}
