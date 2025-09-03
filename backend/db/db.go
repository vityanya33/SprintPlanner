package db

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

var Pool *pgxpool.Pool

//go:embed migrations/*.sql
var migrationFiles embed.FS

func InitDB() error {
	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		// Для Docker-среды используем host=db вместо localhost
		dbUrl = "postgres://postgres:admin123@localhost:5432/sprintplanner?sslmode=disable"
	}

	var err error
	const maxAttempts = 10
	const delay = 2 * time.Second

	fmt.Printf("⏱️ Initializing DB connection (maxAttempts=%d, delay=%s)\n", maxAttempts, delay)

	for i := 1; i <= maxAttempts; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		// Попытка создать пул соединений
		pool, connErr := pgxpool.New(ctx, dbUrl)
		if connErr == nil {
			pingErr := pool.Ping(ctx)
			if pingErr == nil {
				fmt.Println("✅ Connected to PostgreSQL")
				Pool = pool
				cancel()

				// Запускаем миграции после успешного подключения
				fmt.Println("🚀 Starting DB migrations with goose...")
				if migrationErr := runMigrations(); migrationErr != nil {
					pool.Close()
					return fmt.Errorf("❌ Migration failed: %w", migrationErr)
				}
				fmt.Println("✅ Migrations completed successfully")

				return nil
			} else {
				err = pingErr
			}
			pool.Close()
		} else {
			err = connErr
		}

		cancel()
		fmt.Printf("⏳ Attempt %d/%d: can't connect to DB: %v\n", i, maxAttempts, err)
		time.Sleep(delay)
	}

	return fmt.Errorf("❌ Failed to connect to DB after %d attempts: %w", maxAttempts, err)
}

// runMigrations выполняет все миграции из папки migrations c помощью goose
func runMigrations() error {
	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		dbUrl = "postgres://postgres:admin123@localhost:5432/sprintplanner?sslmode=disable"
	}

	// Открываем подключение через database/sql с драйвером pgx для goose
	sqlDB, err := sql.Open("pgx", dbUrl)
	if err != nil {
		return err
	}
	defer sqlDB.Close()

	if err := sqlDB.Ping(); err != nil {
		return err
	}

	// Используем встроенные файлы миграций
	goose.SetBaseFS(migrationFiles)
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	if err := goose.Up(sqlDB, "migrations"); err != nil {
		return err
	}

	return nil
}
