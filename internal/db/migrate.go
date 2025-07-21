package db

import (
	"context"
	"embed"
	"fmt"
    "errors"
	"io/fs"
	"log"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

func ApplyMigrations(ctx context.Context) error {
    // 1. Получаем все файлы миграций
    files, err := fs.Glob(migrationsFS, "migrations/*.sql")
    if err != nil {
        return fmt.Errorf("ошибка чтения миграций: %w", err)
    }
    sort.Strings(files)

    conn, err := Pool.Acquire(ctx)
    if err != nil {
        return err
    }
    defer conn.Release()

    for _, file := range files {
        // 2. Извлекаем номер миграции
        versionStr := strings.Split(filepath.Base(file), "_")[0]
        version, err := strconv.Atoi(versionStr)
        if err != nil {
            return fmt.Errorf("неверный формат номера в %s: %w", file, err)
        }

        // 3. Проверяем, была ли уже применена миграция
        var dummy int
        err = conn.QueryRow(ctx,
            "SELECT 1 FROM schema_versions WHERE version = $1",
            version,
        ).Scan(&dummy)

        if err == nil {
            continue
        } else if !errors.Is(err, pgx.ErrNoRows) {
            return fmt.Errorf("ошибка проверки миграции %03d: %w", version, err)
        }

        // 4. Читаем SQL файл
        sql, err := migrationsFS.ReadFile(file)
        if err != nil {
            return fmt.Errorf("ошибка чтения файла %s: %w", file, err)
        }

        // 5. Выполняем SQL
        if _, err := conn.Exec(ctx, string(sql)); err != nil {
            return fmt.Errorf("ошибка выполнения миграции %03d: %w", version, err)
        }

        // 6. Записываем версию в schema_versions
        if _, err := conn.Exec(ctx,
            "INSERT INTO schema_versions (version) VALUES ($1)",
            version,
        ); err != nil {
            return fmt.Errorf("ошибка записи миграции %03d: %w", version, err)
        }

        log.Printf("Применена миграция %03d: %s", version, file)
    }

    return nil
}
