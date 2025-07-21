package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/CicadaHymn/guitar-shop-api/internal/db"
	"github.com/CicadaHymn/guitar-shop-api/internal/handlers"
)

func main() {
	if err := db.InitDB(); err != nil {
		log.Fatal("DB init error:", err)
	}
	defer db.Pool.Close()

	// проверка аргумента migrate 
	if len(os.Args) > 1 && os.Args[1] == "migrate" {
		if err := runMigrations(); err != nil {
			log.Fatal("Migration failed:", err)
		}
		log.Println("Migrations applied successfully")
		return
	}

	startServer()
}

func runMigrations() error {
	log.Println("Applying migrations...")
	return db.ApplyMigrations(context.Background())
}

func startServer() {
	http.HandleFunc("/guitars", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			handlers.GetGuitars(w, r)
		case "POST":
			handlers.AddGuitar(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}