package services

import (
	"first_task/go-fiber-api/internal/models"
	repo "first_task/go-fiber-api/internal/repository"
)

func NewBookService(r *repo.BookRepo) *BookService {
	return &BookService{Repo: r}
}

type BookService struct {
	Repo *repo.BookRepo
}

func (s *BookService) CreateBook(book *models.Book) error {
	return s.Repo.CreateBook(book)
}
func (s *BookService) GetAllBooksFiltered(search string) ([]models.Book, error) {
	return s.Repo.GetAllBooks(search)
}

func (s *BookService) GetBookByID(id int) (*models.Book, error) {
	return s.Repo.GetBookByID(id)
}

func (s *BookService) Checkin(id int) error {
	return s.Repo.Checkin(id)
}

func (s *BookService) Checkout(id int) error {
	return s.Repo.Checkout(id)
}
