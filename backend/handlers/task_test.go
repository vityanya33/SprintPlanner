package handlers

import (
	"backend/models"
	"bytes"
	"context"
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type mockTaskService struct{}

func (m *mockTaskService) GetAllTasks(ctx context.Context) ([]models.Task, error) {
	return []models.Task{}, nil
}

func (m *mockTaskService) GetTaskByID(ctx context.Context, id string) (*models.Task, error) {
	return &models.Task{}, nil
}

func (m *mockTaskService) CreateTask(ctx context.Context, task *models.Task) error {
	return nil
}

func (m *mockTaskService) UpdateTask(ctx context.Context, id string, task *models.Task) error {
	return nil
}

func (m *mockTaskService) DeleteTask(ctx context.Context, id string) error {
	return nil
}

func (m *mockTaskService) GetAvailableUsers(ctx context.Context, hours int) ([]models.UserWithLoad, error) {
	return []models.UserWithLoad{}, nil
}

func (m *mockTaskService) UpdateTaskUsers(ctx context.Context, taskID string, userIDs []int) error {
	return nil
}

func (m *mockTaskService) SyncTasksWithJira(ctx context.Context, jql string) error {
	return nil
}

func TestGetTasks(t *testing.T) {
	app := fiber.New()

	h := NewTaskHandler(&mockTaskService{})

	app.Get("/tasks", h.GetTasks)

	req := httptest.NewRequest("GET", "/tasks", nil)
	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestGetTaskByID(t *testing.T) {
	app := fiber.New()

	h := NewTaskHandler(&mockTaskService{})
	app.Get("/tasks/:id", h.GetTaskByID)

	req := httptest.NewRequest("GET", "/tasks/1", nil)
	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestPostTasks(t *testing.T) {
	app := fiber.New()

	h := NewTaskHandler(&mockTaskService{})

	app.Post("/tasks", h.PostTasks)

	body := `{
	"title": "Creating app",
	"hours": 10,
	"startDate": "2025-01-01",
	"deadline": "2025-01-02",
	"user_ids": [1, 2, 3]
	}`

	req := httptest.NewRequest("POST", "/tasks", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, 201, resp.StatusCode, func() string {
		body, _ := io.ReadAll(resp.Body)
		return string(body)
	}())
}

func TestPatchTask(t *testing.T) {
	app := fiber.New()

	h := NewTaskHandler(&mockTaskService{})

	app.Patch("/tasks/:id", h.PatchTasks)

	body := `{
	"title": "Download scripts",
	"hours": 5,
	"user_ids": [4,5]
	}`

	req := httptest.NewRequest("PATCH", "/tasks/1", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestDeleteTask(t *testing.T) {
	app := fiber.New()

	h := NewTaskHandler(&mockTaskService{})

	app.Delete("/tasks/:id", h.DeleteTasks)

	req := httptest.NewRequest("DELETE", "/tasks/1", nil)
	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, 204, resp.StatusCode)
}

func TestGetAvailableUsers(t *testing.T) {
	app := fiber.New()

	h := NewTaskHandler(&mockTaskService{})

	app.Get("/available-users", h.GetAvailableUsers)

	req := httptest.NewRequest("GET", "/available-users?hours=10", nil)
	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestPatchTaskUsers(t *testing.T) {
	app := fiber.New()

	h := NewTaskHandler(&mockTaskService{})

	app.Patch("/tasks/:id/users", h.PatchTaskUsers)

	taskID := uuid.New()
	body := bytes.NewBufferString(`{"user_ids": [1, 2, 3]}`)
	req := httptest.NewRequest("PATCH", "/tasks/"+taskID.String()+"/users", body)
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}
