package db

import (
	"context"
	"embed"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

func ApplyMigrations(ctx context.Context) error {
    connStr := fmt.Sprintf(
        "postgres://%s:%s@%s:%s/%s?sslmode=disable&client_encoding=UTF8",
        os.Getenv("DB_USER"),
        os.Getenv("DB_PASSWORD"),
        os.Getenv("DB_HOST"),
        os.Getenv("DB_PORT"),
        os.Getenv("DB_NAME"))


    fmt.Println("DB Connection String:", connStr)
    // инициализация миграций
    source, err := iofs.New(migrationsFS, "migrations")
    if err != nil {
        return fmt.Errorf("failed to init migrations source: %w", err)
    }

    // создание мигратора
    m, err := migrate.NewWithSourceInstance("iofs", source, connStr)
    if err != nil {
        return fmt.Errorf("failed to create migrator: %w", err)
    }
    defer m.Close()

    // применяем миграции
    if err := m.Up(); err != nil && err != migrate.ErrNoChange {
        return fmt.Errorf("failed to apply migrations: %w", err)
    }

    // лог результата
    version, dirty, _ := m.Version()
    log.Printf("Migrations applied. Version: %d (dirty: %v)", version, dirty)
    return nil
}

func RollBackLastMigration(ctx context.Context) error {
	config := Pool.Config()
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable&client_encoding=UTF8",
		config.ConnConfig.User,
		config.ConnConfig.Password,
		config.ConnConfig.Host,
		config.ConnConfig.Port,
		config.ConnConfig.Database,
	)

	source, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		return fmt.Errorf("failed to init migrations source: %w", err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", source, connStr)
	if err != nil {
		return fmt.Errorf("failed to create migrator: %w", err)
	}
	defer m.Close()

	// проверка текущей версии миграций
  currentVersion, dirty, err := m.Version()
    if err != nil {
        if err == migrate.ErrNilVersion {
            log.Println("no migrations applied yet - nothing to rollback")
            return nil
    }
        return fmt.Errorf("failed to get current migration version: %w", err)
    }
    // проверка ДБ на статус dirty
    if dirty {
        return fmt.Errorf("cannot rollback, database is dirty")
    }

	if err := m.Steps(-1); err != nil {
        if err == migrate.ErrNoChange {
            log.Println("no migrations applied yet - nothing to rollback")
            return nil
        }
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	newVersion, _, _ := m.Version()
    log.Printf("Successfully rolled back from version %d to %d", currentVersion, newVersion)

	return nil
}
