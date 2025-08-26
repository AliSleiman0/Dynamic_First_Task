package models

type Book struct {
	ID            int    `gorm:"primaryKey;autoIncrement" json:"id"`
	Title         string `gorm:"size:255;not null" json:"title"`
	PublishedYear int    `json:"published_year"`
	Quantity      int    `json:"quantity"`
	Genre         string `gorm:"size:100" json:"genre"`
	Img_url       string `json:"img_url"`
	
	PublisherID   int    `json:"publisher_id"` // Foreign Key
	Publisher     User   `gorm:"foreignKey:PublisherID" json:"publisher"`
}
