package repo

import (
	"errors"
	"first_task/go-fiber-api/internal/models"

	"gorm.io/gorm"
)

type BookRepo struct {
	DB *gorm.DB
}

func (r *BookRepo) CreateBook(book *models.Book) error {
	return r.DB.Create(book).Error
}
func (r *BookRepo) GetAllBooks() ([]models.Book, error) {
	var books []models.Book
	result := r.DB.Find(&books)
	return books, result.Error
}
func (r *BookRepo) GetBookByID(id int) (*models.Book, error) {
	var book models.Book
	result := r.DB.First(&book, id)
	return &book, result.Error
}
func (r *BookRepo) Checkin(id int) error {
	var book models.Book
	results := r.DB.First(&book, id)
	if results.Error != nil {
		return results.Error
	}
	book.Quantity = book.Quantity + 1
	return r.DB.Save(&book).Error
}
func (r *BookRepo) Checkout(id int) error {
	var book models.Book
	results := r.DB.First(&book, id)
	if results.Error != nil {
		return results.Error
	}
	if book.Quantity == 0 {
		return errors.New("book not available for checkout")
	}
	book.Quantity = book.Quantity - 1
	return r.DB.Save(&book).Error
}
