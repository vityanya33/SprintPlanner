package services

import (
	"backend/models"
	"backend/repositories"
	"context"

	"github.com/google/uuid"
)

type TaskService interface {
	GetAllTasks(ctx context.Context) ([]models.Task, error)
	GetTaskByID(ctx context.Context, id uuid.UUID) (*models.Task, error)
	CreateTask(ctx context.Context, task *models.Task) error
	UpdateTask(ctx context.Context, id uuid.UUID, task *models.Task) error
	DeleteTask(ctx context.Context, id uuid.UUID) error
	UpdateTaskUsers(ctx context.Context, taskID uuid.UUID, userIDs []int) error
	GetAvailableUsers(ctx context.Context, hours int) ([]models.UserWithLoad, error)
}

// Исправьте на TaskServiceImpl вместо taskService
type TaskServiceImpl struct {
	taskRepo repositories.TaskRepository
}

// И возвращайте интерфейс TaskService
func NewTaskService(taskRepo repositories.TaskRepository) TaskService {
	return &TaskServiceImpl{taskRepo: taskRepo}
}

func (s *TaskServiceImpl) GetAllTasks(ctx context.Context) ([]models.Task, error) {
	return s.taskRepo.GetAllTasks(ctx)
}

func (s *TaskServiceImpl) GetTaskByID(ctx context.Context, id uuid.UUID) (*models.Task, error) {
	return s.taskRepo.GetTaskByID(ctx, id)
}

func (s *TaskServiceImpl) CreateTask(ctx context.Context, task *models.Task) error {
	return s.taskRepo.CreateTask(ctx, task)
}

func (s *TaskServiceImpl) UpdateTask(ctx context.Context, id uuid.UUID, task *models.Task) error {
	return s.taskRepo.UpdateTask(ctx, id, task)
}

func (s *TaskServiceImpl) DeleteTask(ctx context.Context, id uuid.UUID) error {
	return s.taskRepo.DeleteTask(ctx, id)
}

func (s *TaskServiceImpl) UpdateTaskUsers(ctx context.Context, taskID uuid.UUID, userIDs []int) error {
	return s.taskRepo.UpdateTaskUsers(ctx, taskID, userIDs)
}

func (s *TaskServiceImpl) GetAvailableUsers(ctx context.Context, hours int) ([]models.UserWithLoad, error) {
	users, err := s.taskRepo.GetAvailableUsers(ctx)
	if err != nil {
		return nil, err
	}

	var available []models.UserWithLoad
	for _, user := range users {
		// Простой расчет: свободный ресурс = общий ресурс - занятость
		// Фильтруем тех, у кого свободного ресурса >= требуемого ресурса задачи
		if user.Free >= hours {
			available = append(available, user)
		}
	}

	return available, nil
}
