package models

type Book struct {
	ID            int    `gorm:"primaryKey;autoIncrement" json:"id"`
	Title         string `gorm:"size:255;not null" json:"title"`
	Author        string `gorm:"size:255;not null" json:"author"`
	PublishedYear int    `json:"published_year"`
	Quantity      int    `json:"quantity"`
	Genre         string `gorm:"size:100" json:"genre"`
}
