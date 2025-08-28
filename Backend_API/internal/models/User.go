package models

import (
	"time"
)

// snake case for database and json data
type User struct {
	ID             int       `gorm:"primaryKey;autoIncrement" json:"id"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	Email          string    `json:"email"`
	Password       string    `json:"-"` // hidden from JSON
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	ImgSrc         string    `json:"img_src"`
	BooksPublished []Book    `gorm:"foreignKey:PublisherID" json:"books"`
}

type PublisherWithCount struct {
	ID        int       `json:"id" gorm:"column:id"`
	FirstName string    `json:"first_name" gorm:"column:first_name"`
	LastName  string    `json:"last_name" gorm:"column:last_name"`
	Email     string    `json:"email" gorm:"column:email"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at"`
	ImgSrc    string    `json:"img_src" gorm:"column:img_src"`
	BookCount int       `json:"book_count" gorm:"column:book_count"`
}
