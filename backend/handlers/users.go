package handlers

import (
	"backend/db"
	"backend/models"
	"context"
	"strconv"

	//"strconv"

	"github.com/gofiber/fiber/v2"
)

//Структура необходимая для обработки рабочих часов каждого пользователя
type UserWithLoad struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Role     string `json:"role"`
	Resource int    `json:"resource"`
	Busy     int    `json:"busy"`
	Free     int    `json:"free"`
}

// Get users
func GetUsers(c *fiber.Ctx) error {
	// Шаг 1: Получаем всех пользователей
	rows, err := db.Pool.Query(context.Background(), `
		SELECT id, name, role, resource
		FROM users
	`)
	if err != nil {
		return err
	}
	defer rows.Close()

	var usersWithLoad []UserWithLoad

	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Role, &user.Resource); err != nil {
			return err
		}

		// Шаг 2: Получаем суммарную загрузку (busy) по задачам
		var busy int
		err = db.Pool.QueryRow(
			context.Background(),
			`SELECT COALESCE(SUM(tasks.hours), 0)
			 FROM task_users
			 JOIN tasks ON tasks.id = task_users.task_id
			 WHERE task_users.user_id = $1`, user.ID,
		).Scan(&busy)
		if err != nil {
			return err
		}

		// Шаг 3: Считаем доступные часы
		free := user.Resource - busy
		if free < 0 {
			free = 0 // чтобы не уйти в минус
		}

		// Шаг 4: Добавляем в массив
		usersWithLoad = append(usersWithLoad, UserWithLoad{
			ID:       user.ID,
			Name:     user.Name,
			Role:     user.Role,
			Resource: user.Resource,
			Busy:     busy,
			Free:     free,
		})
	}

	return c.JSON(usersWithLoad)
}

// GETUSER
func GetUserByID(c *fiber.Ctx) error {
	idStr := c.Params("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	var user models.User
	err = db.Pool.QueryRow(
		context.Background(),
		"SELECT id, name, role, resource FROM users WHERE id = $1",
		id,
	).Scan(&user.ID, &user.Name, &user.Role, &user.Resource)

	if err != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}

	return c.JSON(user)
}

// PATCH users
func PatchUsers(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("invalid id")
	}
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("invalid body")
	}

	err = db.Pool.QueryRow(
		context.Background(),
		"UPDATE users SET name = $1, role = $2, resource = $3 WHERE id = $4 RETURNING id, name, role, resource",
		user.Name, user.Role, user.Resource, id).Scan(&user.ID, &user.Name, &user.Role, &user.Resource)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("DB update failed")
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

// POST users
func PostUsers(c *fiber.Ctx) error {
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	// Валидация
	if user.Name == "" || user.Role == "" || user.Resource <= 0 {
		return c.Status(fiber.StatusBadRequest).SendString("Missing required fields")
	}

	err := db.Pool.QueryRow(
		context.Background(),
		"INSERT INTO users (name, role, resource) VALUES ($1, $2, $3) RETURNING id",
		user.Name, user.Role, user.Resource,
	).Scan(&user.ID)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "database error",
			"info": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(user)
}

// DELETE users
func DeleteUsers(c *fiber.Ctx) error {
	idStr := c.Params("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	var deletedID int64
	err = db.Pool.QueryRow(
		context.Background(),
		"DELETE FROM users WHERE id = $1 RETURNING id",
		id).Scan(&deletedID)
	if err != nil {
		return c.SendStatus(404)
	}
	return c.SendStatus(fiber.StatusNoContent)
}
