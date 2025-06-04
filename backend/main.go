package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"

	"backend/db"
	"backend/handlers"
)

func main() {
	// Инициализируем БД
	if err := db.InitDB(); err != nil {
		log.Fatalf("❌ Ошибка подключения к базе данных: %v", err)
	}

	//Подгружаю из env
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️ .env не найден, продолжаем без него")
	}

	app := fiber.New()
	// Разрешаем доступ с фронта (CORS)
	app.Use(cors.New())

	app.Get("/users", handlers.GetUsers)

	// Простой роут 2
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Backend работает и подключён к БД 🎯")
	})

	// Запуск сервера
	if err := app.Listen(":3000"); err != nil {
		panic(err)
	}
}
