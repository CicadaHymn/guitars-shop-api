package app

import (
	"context"
	"fmt"
	"log"


	"github.com/CicadaHymn/guitar-shop-api/internal/api"
	"github.com/CicadaHymn/guitar-shop-api/internal/db"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

//  запускает веб-сервер
func StartServer() {
	godotenv.Load()

	if err := db.InitDB(); err != nil {
		log.Fatal("DB init error:", err)
	}
	defer db.Pool.Close()

	r := gin.Default()
	setupRoutes(r)

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func setupRoutes(r *gin.Engine) {
	api.SetupRouters(r)
}

// выполнение всех невыполненных миграций
func RunAllMigrations() error {
	if err := godotenv.Load(); err != nil {
		return fmt.Errorf("failed to load .env: %w", err)
	}
	
	if err := db.InitDB(); err != nil {
		return fmt.Errorf("failed to init DB: %w", err)
	}
	defer db.Pool.Close()
	log.Println("Applying migrations...")
	return db.ApplyMigrations(context.Background())
}

// откат миграции на -1 версию
func RunRollback() error {
    if err := godotenv.Load(); err != nil {
        return fmt.Errorf("failed to load .env: %w", err)
    }
    
    if err := db.InitDB(); err != nil {
        return fmt.Errorf("failed to init DB: %w", err)
    }
    defer db.Pool.Close()

    log.Println("Откат последней миграции...")
    return db.RollBackLastMigration(context.Background())
}