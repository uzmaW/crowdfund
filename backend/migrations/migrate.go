package main

import (
    "database/sql"
    "fmt"
    "log"
    "os"

    _ "github.com/lib/pq"
    "github.com/golang-migrate/migrate/v4"
    "github.com/golang-migrate/migrate/v4/database/postgres"
    _ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
    dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
        os.Getenv("POSTGRES_USER"),
        os.Getenv("POSTGRES_PASSWORD"),
        os.Getenv("POSTGRES_HOST"),
        os.Getenv("POSTGRES_PORT"),
        os.Getenv("POSTGRES_DB"),
    )

    db, err := sql.Open("postgres", dbURL)
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
    defer db.Close()

    driver, err := postgres.WithInstance(db, &postgres.Config{})
    if err != nil {
        log.Fatalf("Failed to create migration driver: %v", err)
    }

    m, err := migrate.NewWithDatabaseInstance(
        "file://migrations", // Directory containing migration files
        "postgres", driver)
    if err != nil {
        log.Fatalf("Failed to create migration instance: %v", err)
    }

    if len(os.Args) > 1 {
        switch os.Args[1] {
        case "up":
            if err := m.Up(); err != nil {
                log.Fatalf("Migration failed: %v", err)
            }
            log.Println("Migrations applied successfully")
        case "down":
            if err := m.Down(); err != nil {
                log.Fatalf("Rollback failed: %v", err)
            }
            log.Println("Rollback successful")
        }
    } else {
        log.Println("Usage: go run migrate.go [up|down]")
    }
}