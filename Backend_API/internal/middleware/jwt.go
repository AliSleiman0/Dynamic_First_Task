// internal/middleware/jwt.go
package middleware

import (
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/contrib/jwt"
)

// NewJWT returns a Fiber middleware configured for HS256
//fiber.Handler function b red a middleware function ta aamallu attach to routes.
func NewJWT(secret string) fiber.Handler { //jwtware.New is the function that creates the middleware
	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(secret)}, //middleware byestaamel this to verify incoming access tokens.
		SuccessHandler: func(c *fiber.Ctx) error {
			// Optionally do extra validation (aud/iss) after signature valid
			// the middleware already sets c.Locals("user") 
			return c.Next() //continue to the requested route
		},
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			// Return JSON error instead of default
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid or missing token",
			})
		},
		
	})
}
