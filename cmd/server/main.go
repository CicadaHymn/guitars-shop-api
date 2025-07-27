		package main

		import (
			"context"
			"log"
			"net/http"
			"os"

			"github.com/gin-gonic/gin"
			"github.com/CicadaHymn/guitar-shop-api/internal/db"
			"github.com/CicadaHymn/guitar-shop-api/internal/api"

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
			// экземпляр роутера
			r := gin.Default()

			api.SetupRouters(r)
			
			r.GET("/", func(c *gin.Context) {
				c.String(http.StatusOK, "Hello, World!")	
			})

			r.Run(":8080")

		}

		func runMigrations() error {
			log.Println("Applying migrations...")
			return db.ApplyMigrations(context.Background())
		}
