package main

import (
	"fmt"
	"os"

	"github.com/fokal/fokal-core/pkg/conn"
	migrate "github.com/rubenv/sql-migrate"
)

func main() {
	postgresURL := os.Getenv("DATABASE_URL")
	if postgresURL == "" {
		fmt.Fprintf(os.Stderr, "Postgres URL not set at POSTGRES_URL")
		os.Exit(1)
	}
	db := conn.DialPostgres(postgresURL)

	migrations := &migrate.FileMigrationSource{
		Dir: "migrations",
	}

	n, err := migrate.Exec(db.DB, "postgres", migrations, migrate.Up)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to apply migrations: %+v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Applied %d migrations.\n", n)
}