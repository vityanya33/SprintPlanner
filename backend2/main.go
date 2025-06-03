package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New()

	// Разрешаем доступ с фронта (CORS)
	app.Use(cors.New())

	// Простой роут 2
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("SprintPlanner backend работает  🎯🎯")
	})

	// Запуск сервера
	err := app.Listen(":3000")
	if err != nil {
		panic(err)
	}
}
