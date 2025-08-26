package handlers

import (
	"first_task/go-fiber-api/internal/services"
	utils "first_task/go-fiber-api/pkg"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	TokenSvc    *services.TokenService
	UserService *services.UserService
}

func NewAuthHandler(ts *services.TokenService, us *services.UserService) *AuthHandler {
	return &AuthHandler{
		TokenSvc:    ts,
		UserService: us,
	}
}

type LoginRequest struct {
	ID   int    `json:"id"`
	Pass string `json:"pass"`
}

// Login godoc
// @Summary Login a user
// @Description Authenticate a user by ID and password, returning an access token
// @Tags auth
// @Accept  json
// @Produce  json
// @Param   credentials  body  LoginRequest  true  "User ID and Password"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} string
// @Router /login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	// check if user exists using UserHandler logic
	user, err := h.UserService.GetUserByID(req.ID)
	if err != nil || user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "user not found",
		})
	}

	if !utils.CheckPassword(user.Password, req.Pass) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid password",
		})
	}

	// generate access token
	token, err := h.TokenSvc.CreateAccessToken(user.ID, map[string]any{"role": "admin"})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("could not create token")
	}

	return c.JSON(fiber.Map{
		"token":      token,
		"expires_in": h.TokenSvc.ExpiresInSeconds(),
	})
}
