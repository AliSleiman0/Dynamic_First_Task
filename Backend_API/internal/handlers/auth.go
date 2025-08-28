package handlers

import (
	"first_task/go-fiber-api/internal/services"
	utils "first_task/go-fiber-api/pkg"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
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
func (h *AuthHandler) Profile(c *fiber.Ctx) error {
	// --- 1. Authorization ---
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).SendString("Missing Authorization header")
	}

	tokenString := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).SendString("Invalid token")
	}

	// --- 2. Load secret ---
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatalf("JWT_SECRET not set in .env")
	}

	// --- 3. Parse JWT ---
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).SendString("Invalid token")
	}

	// --- 4. Extract user ID from claims ---
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid token claims")
	}

	userIDClaim, exists := claims["sub"]
	if !exists {
		return c.Status(fiber.StatusBadRequest).SendString("Token missing sub claim")
	}

	var userID int
	switch v := userIDClaim.(type) {
	case float64:
		userID = int(v)
	case string:
		id, err := strconv.Atoi(v)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid sub claim")
		}
		userID = id
	default:
		return c.Status(fiber.StatusBadRequest).SendString("Unknown sub claim type")
	}

	// --- 5. Fetch user ---
	user, err := h.UserService.GetUserByID(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).SendString("User not found")
	}

	// --- 6. Build HTML fragment for HTMX ---
	profileHTML := fmt.Sprintf(`
<div class="card shadow-sm" style="max-width: 600px; margin: auto;">
  <div class="row g-0 align-items-center">
    <div class="col-md-4 text-center p-3">
      <img src="%s" class="img-fluid rounded-circle" alt="Profile Image">
    </div>
    <div class="col-md-8">
      <div class="card-body">
        <h4 class="card-title">%s %s</h4>
        <p class="card-text"><i class="bi bi-envelope me-1"></i>%s</p>
        <p class="card-text"><small class="text-muted">Joined: %s</small></p>
        <p class="card-text"><small class="text-muted">Last updated: %s</small></p>
      </div>
    </div>
  </div>
</div>
`, user.ImgSrc, user.FirstName, user.LastName, user.Email,
		user.CreatedAt.Format("Jan 2, 2006"), user.UpdatedAt.Format("Jan 2, 2006"))

	// --- 7. Return HTML ---
	return c.Type("html").SendString(profileHTML)
}
