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
func (r *UserRepo) GetAllPublishersWithBookCount() ([]models.PublisherWithCount, error) {
	var publishers []models.PublisherWithCount

	query := `
		SELECT DISTINCT u.id, u.first_name, u.last_name, u.email, u.created_at , u.updated_at , u.img_src,
		       COUNT(b.id) as book_count
		FROM users u
		INNER JOIN books b ON u.id = b.publisher_id
		GROUP BY u.id, u.first_name, u.last_name, u.email, u.created_at, u.updated_at, u.img_src
		ORDER BY u.id
	`

	result := r.DB.Raw(query).Scan(&publishers)

	return publishers, result.Error
}
