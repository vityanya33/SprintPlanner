package repositories

import (
	"backend/db"
	"backend/models"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	GetAllUsers(ctx context.Context) ([]models.User, error)
	GetUserByID(ctx context.Context, id int) (*models.User, error)
	CreateUser(ctx context.Context, user *models.User) error
	UpdateUser(ctx context.Context, id int, user *models.User) error
	DeleteUser(ctx context.Context, id int) error
	GetUsersWithWorkload(ctx context.Context) ([]models.UserWithLoad, error)
}

type UserRepositoryImpl struct {
	pool *pgxpool.Pool
}

func NewUserRepository() UserRepository {
	return &UserRepositoryImpl{pool: db.Pool}
}

func (r *UserRepositoryImpl) GetAllUsers(ctx context.Context) ([]models.User, error) {
	query := `SELECT id, name, role, resource FROM users ORDER BY id`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error querying users: %w", err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Role, &user.Resource); err != nil {
			return nil, fmt.Errorf("error scanning user: %w", err)
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *UserRepositoryImpl) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	query := `SELECT id, name, role, resource FROM users WHERE id = $1`
	row := r.pool.QueryRow(ctx, query, id)

	var user models.User
	if err := row.Scan(&user.ID, &user.Name, &user.Role, &user.Resource); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error scanning user: %w", err)
	}

	return &user, nil
}

func (r *UserRepositoryImpl) CreateUser(ctx context.Context, user *models.User) error {
	query := `INSERT INTO users (name, role, resource) VALUES ($1, $2, $3) RETURNING id`
	err := r.pool.QueryRow(ctx, query, user.Name, user.Role, user.Resource).Scan(&user.ID)
	if err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}
	return nil
}

func (r *UserRepositoryImpl) UpdateUser(ctx context.Context, id int, user *models.User) error {
	query := `UPDATE users SET name = $1, role = $2, resource = $3 WHERE id = $4`
	result, err := r.pool.Exec(ctx, query, user.Name, user.Role, user.Resource, id)
	if err != nil {
		return fmt.Errorf("error updating user: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("user with id %d not found", id)
	}

	return nil
}

func (r *UserRepositoryImpl) DeleteUser(ctx context.Context, id int) error {
	query := `DELETE FROM users WHERE id = $1`
	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error deleting user: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("user with id %d not found", id)
	}

	return nil
}

func (r *UserRepositoryImpl) GetUsersWithWorkload(ctx context.Context) ([]models.UserWithLoad, error) {
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
		return nil, fmt.Errorf("error querying users with workload: %w", err)
	}
	defer rows.Close()

	var users []models.UserWithLoad
	for rows.Next() {
		var user models.UserWithLoad
		if err := rows.Scan(&user.ID, &user.Name, &user.Role, &user.Resource, &user.Busy, &user.Free); err != nil {
			return nil, fmt.Errorf("error scanning user with workload: %w", err)
		}
		users = append(users, user)
	}

	return users, nil
}
