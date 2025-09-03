package repositories

import (
	"backend/db"
	"backend/models"
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TaskRepository interface {
	GetAllTasks(ctx context.Context) ([]models.Task, error)
	GetTaskByID(ctx context.Context, id string) (*models.Task, error)
	CreateTask(ctx context.Context, task *models.Task) error
	BulkCreateOrUpdateTasks(ctx context.Context, tasks []models.Task) error
	UpdateTask(ctx context.Context, id string, task *models.Task) error
	DeleteTask(ctx context.Context, id string) error
	UpdateTaskUsers(ctx context.Context, taskID string, userIDs []int) error
	GetAvailableUsers(ctx context.Context) ([]models.UserWithLoad, error)
}

type TaskRepositoryImpl struct {
	pool *pgxpool.Pool
}

func NewTaskRepository() TaskRepository {
	return &TaskRepositoryImpl{pool: db.Pool}
}

func (r *TaskRepositoryImpl) GetAllTasks(ctx context.Context) ([]models.Task, error) {
	query := `
		SELECT 
			t.id, 
			t.title, 
			t.hours, 
			t.start_date, 
			t.deadline,
			COALESCE(ARRAY_AGG(tu.user_id) FILTER (WHERE tu.user_id IS NOT NULL), '{}') as user_ids
		FROM tasks t
		LEFT JOIN task_users tu ON t.id = tu.task_id
		GROUP BY t.id
		ORDER BY t.title
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		log.Printf("Error querying tasks: %v", err)
		return nil, fmt.Errorf("error querying tasks: %w", err)
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		if err := rows.Scan(&task.ID, &task.Title, &task.Hours, &task.StartDate, &task.Deadline, &task.UserIDs); err != nil {
			log.Printf("Error scanning task: %v", err)
			return nil, fmt.Errorf("error scanning task: %w", err)
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (r *TaskRepositoryImpl) GetTaskByID(ctx context.Context, id string) (*models.Task, error) {
	query := `
		SELECT 
			t.id, 
			t.title, 
			t.hours, 
			t.start_date, 
			t.deadline,
			COALESCE(ARRAY_AGG(tu.user_id) FILTER (WHERE tu.user_id IS NOT NULL), '{}') as user_ids
		FROM tasks t
		LEFT JOIN task_users tu ON t.id = tu.task_id
		WHERE t.id = $1
		GROUP BY t.id
	`

	row := r.pool.QueryRow(ctx, query, id)

	var task models.Task
	if err := row.Scan(&task.ID, &task.Title, &task.Hours, &task.StartDate, &task.Deadline, &task.UserIDs); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error scanning task: %w", err)
	}

	return &task, nil
}

func (r *TaskRepositoryImpl) CreateTask(ctx context.Context, task *models.Task) error {
	// Начинаем транзакцию
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	if task.ID == "" {
		task.ID = uuid.New().String()
	}

	// Создаем задачу с генерацией UUID
	query := `INSERT INTO tasks (id, title, hours, start_date, deadline) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err = tx.QueryRow(ctx, query, task.ID, task.Title, task.Hours, task.StartDate, task.Deadline).Scan(&task.ID)
	if err != nil {
		return fmt.Errorf("error creating task: %w", err)
	}

	// Добавляем связи с пользователями
	if len(task.UserIDs) > 0 {
		if err := r.updateTaskUsers(ctx, tx, task.ID, task.UserIDs); err != nil {
			return err
		}
	}

	// Коммитим транзакцию
	return tx.Commit(ctx)
}

func (r *TaskRepositoryImpl) BulkCreateOrUpdateTasks(ctx context.Context, tasks []models.Task) error {
	if len(tasks) == 0 {
		return nil
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	batch := &pgx.Batch{}
	hasUpserts := false
	for _, t := range tasks {
		if t.ID != "" {
			hasUpserts = true
			batch.Queue(
				`INSERT INTO tasks (id, title, hours, start_date, deadline)
				 VALUES ($1, $2, $3, $4, $5)
				 ON CONFLICT (id) DO UPDATE SET
				   title = EXCLUDED.title,
				   hours = EXCLUDED.hours,
				   start_date = EXCLUDED.start_date,
				   deadline = EXCLUDED.deadline,
				   updated_at = CURRENT_TIMESTAMP`,
				t.ID, t.Title, t.Hours, t.StartDate, t.Deadline,
			)
		}
	}

	if hasUpserts {
		br := tx.SendBatch(ctx, batch)
		for i := 0; i < batch.Len(); i++ {
			if _, err := br.Exec(); err != nil {
				br.Close()
				return fmt.Errorf("error executing upsert batch: %w", err)
			}
		}
		if err := br.Close(); err != nil {
			return fmt.Errorf("error closing upsert batch: %w", err)
		}
	}

	rows := make([][]interface{}, 0)
	for _, t := range tasks {
		if t.ID == "" {
			t.ID = uuid.New().String()
			rows = append(rows, []interface{}{t.Title, t.Hours, t.StartDate, t.Deadline})
		}
	}
	if len(rows) > 0 {
		_, err = tx.CopyFrom(
			ctx,
			pgx.Identifier{"tasks"},
			[]string{"title", "hours", "start_date", "deadline"},
			pgx.CopyFromRows(rows),
		)
		if err != nil {
			return fmt.Errorf("error bulk inserting tasks: %w", err)
		}
	}

	return tx.Commit(ctx)
}

func (r *TaskRepositoryImpl) UpdateTask(ctx context.Context, id string, task *models.Task) error {
	query := `UPDATE tasks SET title = $1, hours = $2, start_date = $3, deadline = $4 WHERE id = $5`
	result, err := r.pool.Exec(ctx, query, task.Title, task.Hours, task.StartDate, task.Deadline, id)
	if err != nil {
		return fmt.Errorf("error updating task: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("task with id %s not found", id)
	}

	return nil
}

func (r *TaskRepositoryImpl) DeleteTask(ctx context.Context, id string) error {
	query := `DELETE FROM tasks WHERE id = $1`
	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error deleting task: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("task with id %s not found", id)
	}

	return nil
}

func (r *TaskRepositoryImpl) UpdateTaskUsers(ctx context.Context, taskID string, userIDs []int) error {
	// Начинаем транзакцию
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Удаляем старые связи
	deleteQuery := `DELETE FROM task_users WHERE task_id = $1`
	if _, err := tx.Exec(ctx, deleteQuery, taskID); err != nil {
		return fmt.Errorf("error deleting task users: %w", err)
	}

	// Добавляем новые связи
	if len(userIDs) > 0 {
		if err := r.updateTaskUsers(ctx, tx, taskID, userIDs); err != nil {
			return err
		}
	}

	// Коммитим транзакцию
	return tx.Commit(ctx)
}

func (r *TaskRepositoryImpl) updateTaskUsers(ctx context.Context, tx pgx.Tx, taskID string, userIDs []int) error {
	valueStrings := make([]string, 0, len(userIDs))
	valueArgs := make([]interface{}, 0, len(userIDs)*2)

	for i, userID := range userIDs {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d)", i*2+1, i*2+2))
		valueArgs = append(valueArgs, taskID, userID)
	}

	query := fmt.Sprintf("INSERT INTO task_users (task_id, user_id) VALUES %s", strings.Join(valueStrings, ","))
	_, err := tx.Exec(ctx, query, valueArgs...)
	if err != nil {
		return fmt.Errorf("error inserting task users: %w", err)
	}

	return nil
}

func (r *TaskRepositoryImpl) GetAvailableUsers(ctx context.Context) ([]models.UserWithLoad, error) {
	query := `
		SELECT 
			u.id, 
			u.name, 
			u.role, 
			u.resource,
			COALESCE(SUM(t.hours), 0) as busy,
			u.resource - COALESCE(SUM(t.hours), 0) as free
		FROM users u
		LEFT JOIN task_users tu ON u.id = tu.user_id
		LEFT JOIN tasks t ON tu.task_id = t.id
		GROUP BY u.id, u.name, u.role, u.resource
		ORDER BY u.id
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error querying available users: %w", err)
	}
	defer rows.Close()

	var users []models.UserWithLoad
	for rows.Next() {
		var user models.UserWithLoad
		if err := rows.Scan(&user.ID, &user.Name, &user.Role, &user.Resource, &user.Busy, &user.Free); err != nil {
			return nil, fmt.Errorf("error scanning available user: %w", err)
		}
		users = append(users, user)
	}

	return users, nil
}
