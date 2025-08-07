package main

import (
	"log"
	"os"

	"github.com/CicadaHymn/guitar-shop-api/internal/app"
)

func main() {
	if len(os.Args) < 2 {
		printHelp()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "migrate":
		handleMigrate()
	case "server":
		app.StartServer()
	default:
		printHelp()
	}
}

func handleMigrate() {
	// обычный "migrate" без аргументов
	if len(os.Args) == 2 {
		if err := app.RunAllMigrations(); err != nil {
			log.Fatal("error applying migrations:", err)
		}
		log.Println("All migrations applied successfully")
		return
	}

	// подкоманды "up/down"
	switch os.Args[2] {
	case "up":
		log.Println("'migrate up' (на +1 версию)")
		//app.RunSingleMigrationUp()
	case "down":
		if err := app.RunRollback(); err != nil {
			log.Fatal("error rolling back migrations:", err)
		}
		log.Println("last migration rolled back successfully")
	default:
		log.Fatal("invalid command, can be only: migrate, migrate up, migrate down")
	}
}

func printHelp() {
	log.Println(`
  Использование:
  guitar-shop migrate       - Применить ВСЕ миграции
  guitar-shop migrate up    - Применить одну миграцию (+1) 
  guitar-shop migrate down  - Откатить последнюю миграцию
  guitar-shop server        - Запустить веб-сервер`)
}