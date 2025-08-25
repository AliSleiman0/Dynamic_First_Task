package repo

import (
	"errors"
	"first_task/go-fiber-api/internal/models"

	"gorm.io/gorm"
)

type UserRepo struct {
	DB *gorm.DB
}

func (r *UserRepo) CreateUser(user *models.User) error {
	return r.DB.Create(user).Error
}
func (r *UserRepo) GetAllUsers() ([]models.User, error) {
	var users []models.User
	result := r.DB.Find(&users)
	return users, result.Error
}
func (r *UserRepo) GetUserByID(id int) (*models.User, error) {
	var user models.User
	result := r.DB.First(&user, id)
	return &user, result.Error
}
func (r *UserRepo) UpdateUser(user *models.User) error {
	if user.ID == 0 {
		return errors.New("invalid user id")
	}
	res := r.DB.Model(&models.User{}).
		Where("id = ?", user.ID).
		Select("first_name", "last_name", "email").
		Updates(user)

	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("user not found")
	}
	return nil
}
