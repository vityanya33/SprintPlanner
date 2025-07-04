package handlers

import (
	"backend/db"
	"backend/models"
	"context"
	"strconv"
	"time"

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
		var (
			id       int
			title    string
			userId   int
			start    time.Time
			deadline time.Time
		)

		if err := rows.Scan(&id, &title, &userId, &start, &deadline); err != nil {
			return err
		}

		tasks = append(tasks, models.Task{
			ID:        id,
			Title:     title,
			UserID:    userId,
			StartDate: start.Format("2006-01-02"),
			Deadline:  deadline.Format("2006-01-02"),
		})
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

	var (
		task     models.Task
		start    time.Time
		deadline time.Time
	)

	err = db.Pool.QueryRow(
		context.Background(),
		"SELECT id, title, user_id, start_date, deadline FROM tasks WHERE id = $1",
		id,
	).Scan(&task.ID, &task.Title, &task.UserID, &start, &deadline)

	if err != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}

	task.StartDate = start.Format("2006-01-02")
	task.Deadline = deadline.Format("2006-01-02")

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
		task.Title, task.UserID, task.StartDate, task.Deadline,
	).Scan(&task.ID); err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(task)
}

// PATCH task
func PatchTasks(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid task ID")
	}

	var task models.Task
	if err := c.BodyParser(&task); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid body")
	}

	cmdTag, err := db.Pool.Exec(
		context.Background(),
		`UPDATE tasks
		 SET title = $1, user_id = $2, start_date = $3, deadline = $4
		 WHERE id = $5`,
		task.Title, task.UserID, task.StartDate, task.Deadline, id,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Update failed")
	}

	if cmdTag.RowsAffected() == 0 {
		return c.Status(fiber.StatusNotFound).SendString("Task not found")
	}

	// Возвращаем обратно то, что прислали
	task.ID = int(id)
	return c.Status(fiber.StatusOK).JSON(task)
}


// DELETE task
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
		id,
	).Scan(&deletedID)

	if err != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}

	return c.SendStatus(fiber.StatusNoContent)
}
