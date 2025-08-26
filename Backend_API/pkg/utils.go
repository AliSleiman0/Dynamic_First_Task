package utils

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

// ParseID extracts an integer "id" from the URL params.
// Returns an error if conversion fails.
func ParseID(c *fiber.Ctx) (int, error) {
	id := c.Params("id")
	return strconv.Atoi(id)
}
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}
func CheckPassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
