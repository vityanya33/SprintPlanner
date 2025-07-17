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
		log.Fatalf("❌ Failed to connect to DB: %v", err)
	}

	//Подгружаю из env
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️ .env doesnt found")
	}

	app := fiber.New()
	// Разрешаем доступ с фронта (CORS)
	app.Use(cors.New())

    //Запуск UI для openapi
    app.Static("/", "./docs")

	//Запросы по пользователям
	app.Get("/users", handlers.GetUsers)
	app.Get("/users/:id", handlers.GetUserByID)
	app.Post("/users", handlers.PostUsers)
	app.Delete("/users/:id", handlers.DeleteUsers)
	app.Patch("/users/:id", handlers.PatchUsers)

	//Запросы по задачам
	app.Get("/tasks", handlers.GetTasks)
	app.Get("/tasks/:id", handlers.GetTaskByID)
	app.Post("/tasks", handlers.PostTasks)
	app.Delete("/tasks/:id", handlers.DeleteTasks)
	app.Patch("/tasks/:id", handlers.PatchTasks)

	// Запуск сервера
	if err := app.Listen("0.0.0.0:3000"); err != nil {
		panic(err)
	}

}
