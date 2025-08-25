package utils

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// ParseID extracts an integer "id" from the URL params.
// Returns an error if conversion fails.
func ParseID(c *fiber.Ctx) (int, error) {
	id := c.Params("id")
	return strconv.Atoi(id)
}
