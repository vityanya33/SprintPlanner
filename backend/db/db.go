package db

import (
	"context" // Для управления временем выполнения (таймауты, отмена)
	"fmt"     // Для форматирования ошибок и вывода
	"os"      // Для чтения переменных окружения
	"time"    // Для задания таймаутов

	"github.com/jackc/pgx/v5/pgxpool"
	// библиотека для подключения к PostgreSQL через пул соединений
)

// Объявляется глобальная переменная Pool
// Хранит пул соединений с БД PostgreSQL

var Pool *pgxpool.Pool

/*
Инициализирует подключение к базе данных.
Возвращает error, если что-то пошло не так.

Сначала пытается взять строку подключения к БД из переменной окружения DATABASE_URL.
Если она не задана — использует дефолтную строку для подключения к локальной БД PostgreSQL
*/
func InitDB() error {
	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		dbUrl = "postgres://postgres:admin123@localhost:5432/sprintplanner?sslmode=disable"
	}

	//Создаёт контекст с таймаутом в 5 секунд.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	//defer cancel() гарантирует, что ресурсы освобождаются после завершения.
	defer cancel()

	//Инициализирует пул соединений к БД (объект pgxpool.pool)
	pool, err := pgxpool.New(ctx, dbUrl)
	if err != nil {
		return fmt.Errorf("Error to connect to DB: %w", err)
	}
	//Пингует, чтобы понять можно ли подключиться
	if err := pool.Ping(ctx); err != nil {
		return fmt.Errorf("Failed connection to DB: %w", err)
	}

	fmt.Println("✅Connected to PostgreSQL")

	Pool = pool
	return nil
}
