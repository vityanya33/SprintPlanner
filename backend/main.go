package main

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"

	"backend/db"
	"backend/handlers"
	"backend/repositories"
	"backend/services"
)

func main() {
	// Инициализируем БД
	if err := db.InitDB(); err != nil {
		log.Fatalf("❌ Failed to connect to DB: %v", err)
	}

	//Подгружаю из env
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️ .env not found - using default DB connection")
	}

	// Проверяем, что пул соединений создан
	if db.Pool == nil {
		log.Fatalf("❌ DB pool is nil")
	}

	// Простая проверка подключения
	ctx := context.Background()
	conn, err := db.Pool.Acquire(ctx)
	if err != nil {
		log.Fatalf("❌ Failed to acquire DB connection: %v", err)
	}
	defer conn.Release()

	// Проверяем, что можем выполнить простой запрос
	var version string
	err = conn.QueryRow(ctx, "SELECT version()").Scan(&version)
	if err != nil {
		log.Fatalf("❌ Failed to query DB: %v", err)
	}
	log.Printf("✅ Database connected: %s", version)

	// Инициализируем репозитории и сервисы
	userRepo := repositories.NewUserRepository()
	taskRepo := repositories.NewTaskRepository()

	userService := services.NewUserService(userRepo)
	taskService := services.NewTaskService(taskRepo)

	// Инициализируем хендлеры с сервисами
	userHandler := handlers.NewUserHandler(userService)
	taskHandler := handlers.NewTaskHandler(taskService)

	app := fiber.New()
	app.Use(cors.New())

	//Запуск UI для openapi
	app.Static("/", "./docs")

	//Запросы по пользователям
	app.Get("/users", userHandler.GetUsers)
	app.Get("/users/:id", userHandler.GetUserByID)
	app.Post("/users", userHandler.PostUsers)
	app.Post("/users/upload", userHandler.UploadUsersXLS)
	app.Delete("/users/:id", userHandler.DeleteUsers)
	app.Patch("/users/:id", userHandler.PatchUsers)

	//Запросы по задачам
	app.Get("/tasks", taskHandler.GetTasks)
	app.Get("/tasks/available", taskHandler.GetAvailableUsers)
	app.Get("/tasks/:id", taskHandler.GetTaskByID)
	app.Post("/tasks", taskHandler.PostTasks)
	app.Delete("/tasks/:id", taskHandler.DeleteTasks)
	app.Patch("/tasks/:id", taskHandler.PatchTasks)
	app.Patch("/tasks/:id/users", taskHandler.PatchTaskUsers)

	// Запуск сервера
	log.Println("🚀 Server starting on :3000")
	if err := app.Listen("0.0.0.0:3000"); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}
