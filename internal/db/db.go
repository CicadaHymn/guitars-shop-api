package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

func InitDB() error {
	connStr := "postgres://admin:000@localhost:5432/shop?sslmode=disable"
	
	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return fmt.Errorf("ошибка парсинга DSN: %v", err)
	}

	Pool, err = pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return fmt.Errorf("ошибка подключения: %v", err)
	}

	return Pool.Ping(context.Background())
}