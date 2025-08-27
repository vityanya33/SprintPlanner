package handlers

import (
	"backend/models"
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xuri/excelize/v2"
)

type mockUserService struct{}

func (m *mockUserService) GetAllUsers(ctx context.Context) ([]models.User, error) {
	return []models.User{}, nil
}

func (m *mockUserService) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	return &models.User{}, nil
}

func (m *mockUserService) CreateUser(ctx context.Context, user *models.User) error {
	return nil
}

func (m *mockUserService) UpdateUser(ctx context.Context, id int, user *models.User) error {
	return nil
}

func (m *mockUserService) DeleteUser(ctx context.Context, id int) error {
	return nil
}

func (m *mockUserService) GetUsersWithWorkload(ctx context.Context) ([]models.UserWithLoad, error) {
	return []models.UserWithLoad{}, nil
}

// Unit - тесты
func TestGetUsers(t *testing.T) {
	app := fiber.New()

	h := NewUserHandler(&mockUserService{})

	app.Get("/users", h.GetUsers)

	req := httptest.NewRequest("GET", "/users", nil)
	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestGetUserByID(t *testing.T) {
	app := fiber.New()

	h := NewUserHandler(&mockUserService{})

	app.Get("/users/:id", h.GetUserByID)

	req := httptest.NewRequest("GET", "/users/1", nil)
	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestPostUsers(t *testing.T) {
	app := fiber.New()

	h := NewUserHandler(&mockUserService{})

	app.Post("/users", h.PostUsers)

	body := `{
		"name": "Ivan",
		"role": "developer",
		"resource": 100
	}`

	req := httptest.NewRequest("POST", "/users", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, 201, resp.StatusCode)
}

func TestPatchUser(t *testing.T) {
	app := fiber.New()

	h := NewUserHandler(&mockUserService{})

	app.Patch("/users/:id", h.PatchUsers)

	body := `{
		"name": "Ivan",
		"role": "HR",
		"resource": 50
	}`

	req := httptest.NewRequest("PATCH", "/users/1", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestDeleteUser(t *testing.T) {
	app := fiber.New()

	h := NewUserHandler(&mockUserService{})
	app.Delete("/users/:id", h.DeleteUsers)

	req := httptest.NewRequest("DELETE", "/users/1", nil)
	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, 204, resp.StatusCode)
}

func TestPostUsersXLS(t *testing.T) {
	app := fiber.New()

	h := NewUserHandler(&mockUserService{})
	app.Post("/upload", h.UploadUsersXLS)

	f := excelize.NewFile()
	sheet := f.GetSheetName(0)
	_ = f.SetSheetRow(sheet, "A1", &[]string{"Name", "Role", "Resource"})
	_ = f.SetSheetRow(sheet, "A2", &[]interface{}{"Ivan", "Developer", 100})

	tmpfile, err := os.CreateTemp("", "*.xlsx")
	require.NoError(t, err)

	defer os.Remove(tmpfile.Name())

	require.NoError(t, f.SaveAs(tmpfile.Name()))
	tmpfile.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(tmpfile.Name()))
	require.NoError(t, err)

	file, err := os.Open(tmpfile.Name())
	require.NoError(t, err)
	defer file.Close()

	_, err = io.Copy(part, file)
	require.NoError(t, err)
	writer.Close()

	req := httptest.NewRequest("POST", "/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}
