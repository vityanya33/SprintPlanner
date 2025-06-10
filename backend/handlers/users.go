package handlers

import (
	"backend/db"
	"backend/models"
	"context"
	"strconv"

	//"strconv"

	"github.com/gofiber/fiber/v2"
)

// Get users
func GetUsers(c *fiber.Ctx) error {
	rows, err := db.Pool.Query(context.Background(), "SELECT id, name, role from users")
	if err != nil {
		return err
	}
	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Role); err != nil {
			return err
		}
		users = append(users, user)
	}

	return c.JSON(users)
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
		"SELECT id, name, role FROM users WHERE id = $1",
		id,
	).Scan(&user.ID, &user.Name, &user.Role)

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
		"UPDATE users SET name = $1, role = $2 WHERE id = $3 RETURNING id, name, role",
		user.Name, user.Role, id).Scan(&user.ID, &user.Name, &user.Role)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("DB update failed")
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

// POST users
func PostUsers(c *fiber.Ctx) error {
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return err
	}
	if err := db.Pool.QueryRow(
		context.Background(),
		"INSERT INTO users (name, role) VALUES ($1, $2) RETURNING id",
		user.Name, user.Role).Scan(&user.ID); err != nil {
		return err
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
