package services

import (
	"first_task/go-fiber-api/internal/models"
	repo "first_task/go-fiber-api/internal/repository"
)

type UserService struct {
	Repo *repo.UserRepo
}

func NewUserService(r *repo.UserRepo) *UserService {
	return &UserService{Repo: r}
}
func (s *UserService) CreateUser(user *models.User) error {
	return s.Repo.CreateUser(user)
}

func (s *UserService) GetAllUsers() ([]models.User, error) {
	return s.Repo.GetAllUsers()
}

func (s *UserService) GetUserByID(id int) (*models.User, error) {
	return s.Repo.GetUserByID(id)
}
func (s *UserService) UpdateUser(user *models.User) error {
	return s.Repo.UpdateUser(user)
}
func (s *UserService) GetAllPublishersWithBookCount() ([]models.PublisherWithCount, error) {
	return s.Repo.GetAllPublishersWithBookCount()
}
