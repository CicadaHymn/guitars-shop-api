package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

func InitDB() error {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable&client_encoding=UTF8",
		dbUser,
		dbPassword,
		dbHost,
		dbPort,
		dbName,
	)
	
	log.Printf("Connecting to DB: postgres://%s:***@%s:%s/%s", dbUser, dbHost, dbPort, dbName)

	// задаем соединение
	var err error
	Pool, err = pgxpool.New(context.Background(), connStr)
	if err != nil {
		return fmt.Errorf("missing database connection: %v", err)
	}

	// проверка соединения
	if err := Pool.Ping(context.Background()); err != nil {
		return fmt.Errorf("the connection to the database could not be verified: %v", err)
	}

	return nil
}