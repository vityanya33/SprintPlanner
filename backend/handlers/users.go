package handlers

import (
	"backend/db"
	"backend/models"
	"context"
	//"strconv"

	"github.com/gofiber/fiber/v2"
)

// Get /users
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
