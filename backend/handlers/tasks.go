package handlers

import (
	"backend/db"
	"backend/models"
	"context"
	"fmt"
	"strconv"

	//"strconv"

	"github.com/gofiber/fiber/v2"
)

// GET tasks
func GetTasks(c *fiber.Ctx) error {
	rows, err := db.Pool.Query(context.Background(), "SELECT id, title, user_id, start_date, deadline from tasks")
	if err != nil {
		return err
	}
	defer rows.Close()

	var tasks []models.Task

	for rows.Next() {
		var task models.Task
		if err := rows.Scan(&task.ID, &task.Title, &task.UserID, &task.StartDate, &task.Deadline); err != nil {
			return err
		}
		tasks = append(tasks, task)
	}

	return c.JSON(tasks)
}

// GET one task
func GetTaskByID(c *fiber.Ctx) error {
	idStr := c.Params("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	var task models.Task
	err = db.Pool.QueryRow(
		context.Background(),
		"SELECT id, title, user_id, start_date, deadline from tasks WHERE id = $1",
		id,
	).Scan(&task.ID, &task.Title, &task.UserID, &task.StartDate, &task.Deadline)

	if err != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}

	return c.JSON(task)
}

// POST task
func PostTasks(c *fiber.Ctx) error {
	var task models.Task
	if err := c.BodyParser(&task); err != nil {
		return err
	}
	if err := db.Pool.QueryRow(
		context.Background(),
		"INSERT INTO tasks (title, user_id, start_date, deadline) VALUES ($1, $2, $3, $4) RETURNING id",
		task.Title, task.UserID, task.StartDate, task.Deadline).Scan(&task.ID); err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(task)
}

// PATCH users
func PatchTasks(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	var task models.Task
	if err := c.BodyParser(&task); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	fmt.Println(task)
	err = db.Pool.QueryRow(
		context.Background(),
		"UPDATE tasks set title = $1, user_id = $2, start_date = $3, deadline = $4 WHERE id = $5 RETURNING id, title, user_id, start_date, deadline",
		task.Title, task.UserID, task.StartDate, task.Deadline, id).Scan(&task.ID, &task.Title, &task.UserID, &task.StartDate, &task.Deadline)

	if err != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}

	return c.Status(fiber.StatusOK).JSON(task)
}

// DELETE users
func DeleteTasks(c *fiber.Ctx) error {
	idStr := c.Params("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	var deletedID int64
	err = db.Pool.QueryRow(
		context.Background(),
		"DELETE FROM tasks WHERE id = $1 RETURNING id",
		id).Scan(&deletedID)
	if err != nil {
		return c.SendStatus(404)
	}
	return c.SendStatus(fiber.StatusNoContent)
}
