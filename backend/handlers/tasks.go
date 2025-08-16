package handlers

import (
	"backend/db"
	"backend/models"
	"context"
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

// GET tasks
func GetTasks(c *fiber.Ctx) error {
	rows, err := db.Pool.Query(context.Background(), `
        SELECT
            t.id, t.title, t.hours, t.start_date, t.deadline,
            COALESCE(ARRAY_AGG(tu.user_id) FILTER (WHERE tu.user_id IS NOT NULL), '{}') AS user_ids
        FROM tasks t
        LEFT JOIN task_users tu ON t.id = tu.task_id
        GROUP BY t.id
    `)
	if err != nil {
		return err
	}
	defer rows.Close()

	var tasks []models.Task

//новая структура для правильного отображения пользователей задачи
	for rows.Next() {
		var (
			id       int
			title    string
			hours    int
			start    time.Time
			deadline time.Time
			userIDs  []int
		)

		if err := rows.Scan(&id, &title, &hours, &start, &deadline, &userIDs); err != nil {
			return err
		}

		tasks = append(tasks, models.Task{
			ID:        id,
			Title:     title,
			Hours:     hours,
			UserIDs:   userIDs, // массив ID
			StartDate: start.Format("2006-01-02"),
			Deadline:  deadline.Format("2006-01-02"),
		})
	}

	return c.JSON(tasks)
}

// GET one task
func GetTaskByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	var (
		task     models.Task
		start    time.Time
		deadline time.Time
	)

	// Запрос задачи для всех пользователей, связанных с ней
	err = db.Pool.QueryRow(
        context.Background(),
        `SELECT
            t.id, t.hours, t.title, t.start_date, t.deadline,
            COALESCE(ARRAY_AGG(tu.user_id) FILTER (WHERE tu.user_id IS NOT NULL), '{}') AS user_ids
        FROM tasks t
        LEFT JOIN task_users tu ON t.id = tu.task_id
        WHERE t.id = $1
        GROUP BY t.id`,
        id,
    ).Scan(&task.ID, &task.Title, &task.Hours, &start, &deadline, &task.UserIDs)


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
		return c.Status(fiber.StatusBadRequest).SendString("Invalid body")
	}

	// Валидация
	if task.Title == "" || len(task.UserIDs) == 0 || task.StartDate == "" || task.Deadline == "" || task.Hours <= 0 {
		return c.Status(fiber.StatusBadRequest).SendString("Missing required fields")
	}

	// Создание задачи
	err := db.Pool.QueryRow(
		context.Background(),
		"INSERT INTO tasks (title, hours, start_date, deadline) VALUES ($1, $2, $3, $4) RETURNING id",
		task.Title, task.Hours, task.StartDate, task.Deadline,
	).Scan(&task.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to insert task")
	}

	// Добавление связи с пользователями
	for _, userID := range task.UserIDs {
		_, err := db.Pool.Exec(
			context.Background(),
			"INSERT INTO task_users (task_id, user_id) VALUES ($1, $2)",
			task.ID, userID,
		)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to link users to task")
		}
	}

	return c.Status(fiber.StatusCreated).JSON(task)
}


// PATCH task
func PatchTasks(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid task ID")
	}

	var task models.Task
	if err := c.BodyParser(&task); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid body")
	}

	// Обновление задачи
	cmdTag, err := db.Pool.Exec(
		context.Background(),
		`UPDATE tasks SET title = $1, hours= $2, start_date = $3, deadline = $4 WHERE id = $5`,
		task.Title, task.Hours, task.StartDate, task.Deadline, id,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Update failed")
	}
	if cmdTag.RowsAffected() == 0 {
		return c.Status(fiber.StatusNotFound).SendString("Task not found")
	}

	// Очистка старых связей
	_, err = db.Pool.Exec(
		context.Background(),
		"DELETE FROM task_users WHERE task_id = $1",
		id,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to clear task users")
	}

	// Добавление новых связей
	for _, userID := range task.UserIDs {
		_, err := db.Pool.Exec(
			context.Background(),
			"INSERT INTO task_users (task_id, user_id) VALUES ($1, $2)",
			id, userID,
		)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to insert new task users")
		}
	}

	task.ID = id
	return c.Status(fiber.StatusOK).JSON(task)
}

// DELETE task
func DeleteTasks(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	// Удаление сначала связи с пользователями
	_, err = db.Pool.Exec(
		context.Background(),
		"DELETE FROM task_users WHERE task_id = $1",
		id,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to delete task-user links")
	}

	// Удаление самой задачи
	var deletedID int
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

//GET задачи со свободными пользователями
func GetAvailableUsers(c *fiber.Ctx) error {
	start := c.Query("start_date")
	deadlineStr := c.Query("deadline")
	hoursStr := c.Query("hours")
    log.Println("Cырые параметры: ", start, deadlineStr, hoursStr)

	// Преобразуем даты
	layout := "2006-01-02"
	startDate, err := time.Parse(layout, start)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid start_date format"})
	}
	deadlineParsed, err := time.Parse(layout, deadlineStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid deadline format"})
	}

	// Преобразуем часы
	hours, err := strconv.Atoi(hoursStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid hours format"})
	}

	// Считаем количество дней
	days := int(deadlineParsed.Sub(startDate).Hours()/24) + 1

	// Логируем значения
	log.Println("start =", startDate)
	log.Println("deadline =", deadlineParsed)
	log.Println("hours =", hours)
	log.Println("days =", days)

	// SQL-запрос
	query := `
		SELECT u.id, u.name, u.role, u.resource,
			COALESCE(SUM(t.hours), 0) as busy_hours
		FROM users u
		LEFT JOIN task_users tu ON u.id = tu.user_id
		LEFT JOIN tasks t ON tu.task_id = t.id
			AND t.start_date <= $2 AND t.deadline >= $1
		GROUP BY u.id, u.name, u.role, u.resource
	`

	// Проверка запроса
	log.Println("🟢 Выполняем запрос к БД")

	rows, err := db.Pool.Query(context.Background(), query, startDate, deadlineParsed)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	defer rows.Close()

	var available []models.UserWithLoad
	for rows.Next() {
		var u models.UserWithLoad
		var busy int
		if err := rows.Scan(&u.ID, &u.Name, &u.Role, &u.Resource, &busy); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		u.Busy = busy
		totalCapacity := u.Resource * days
		u.Free = totalCapacity - busy
		if u.Free >= hours {
			available = append(available, u)
		}
	}

	return c.JSON(available)
}

//получение всех пользователей
func GetAllUsers(c *fiber.Ctx) error {
	rows, err := db.Pool.Query(context.Background(), `
		SELECT id, name, role, resource FROM users
	`)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Role, &u.Resource); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		users = append(users, u)
	}

	return c.JSON(users)
}

// PATCH для задач с drag & drop
func PatchTaskUsers(c *fiber.Ctx) error {
    idStr := c.Params("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).SendString("Invalid task ID")
    }

    var payload struct {
        UserIDs []int `json:"user_ids"`
    }
    if err := c.BodyParser(&payload); err != nil {
        return c.Status(fiber.StatusBadRequest).SendString("Invalid body")
    }

    // Чистим старые связи
    _, err = db.Pool.Exec(
        context.Background(),
        "DELETE FROM task_users WHERE task_id = $1",
        id,
    )
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).SendString("Failed to clear task users")
    }

    // Вставляем новые связи
    for _, userID := range payload.UserIDs {
        _, err := db.Pool.Exec(
            context.Background(),
            "INSERT INTO task_users (task_id, user_id) VALUES ($1, $2)",
            id, userID,
        )
        if err != nil {
            return c.Status(fiber.StatusInternalServerError).SendString("Failed to insert task user")
        }
    }

    return c.SendStatus(fiber.StatusOK)
}
