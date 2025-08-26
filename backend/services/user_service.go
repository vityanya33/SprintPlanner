package services

import (
	"backend/models"
	"backend/repositories"
	"context"
)

type UserService interface {
	GetAllUsers(ctx context.Context) ([]models.User, error)
	GetUserByID(ctx context.Context, id int) (*models.User, error)
	CreateUser(ctx context.Context, user *models.User) error
	UpdateUser(ctx context.Context, id int, user *models.User) error
	DeleteUser(ctx context.Context, id int) error
	GetUsersWithWorkload(ctx context.Context) ([]models.UserWithLoad, error)
}

type UserServiceImpl struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) UserService {
	return &UserServiceImpl{userRepo: userRepo}
}

func (s *UserServiceImpl) GetAllUsers(ctx context.Context) ([]models.User, error) {
	return s.userRepo.GetAllUsers(ctx)
}

func (s *UserServiceImpl) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	return s.userRepo.GetUserByID(ctx, id)
}

func (s *UserServiceImpl) CreateUser(ctx context.Context, user *models.User) error {
	return s.userRepo.CreateUser(ctx, user)
}

func (s *UserServiceImpl) UpdateUser(ctx context.Context, id int, user *models.User) error {
	return s.userRepo.UpdateUser(ctx, id, user)
}

func (s *UserServiceImpl) DeleteUser(ctx context.Context, id int) error {
	return s.userRepo.DeleteUser(ctx, id)
}

func (s *UserServiceImpl) GetUsersWithWorkload(ctx context.Context) ([]models.UserWithLoad, error) {
	return s.userRepo.GetUsersWithWorkload(ctx)
}
