package db

import (
	"fmt"
	"log"
	"os"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	name := os.Getenv("DB_NAME")
	charset := os.Getenv("DB_CHARSET")

	//dsn yaane data source name
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local", user, pass, host, port, name, charset)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{}) //&gorm.Config{} is just default config of gorm
	if err != nil {
		return nil, err
	}
	return db, nil
}
