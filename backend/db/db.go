package db

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

func InitDB() error {
	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		dbUrl = "postgres://postgres:admin123@localhost:5432/sprintplanner?sslmode=disable"
	}

	var err error
	const maxAttempts = 10
	const delay = 2 * time.Second

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
