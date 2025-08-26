package handlers

import (
	"first_task/go-fiber-api/internal/models"
	"first_task/go-fiber-api/internal/services"
	utils "first_task/go-fiber-api/pkg"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type UserHandler struct {
	Service *services.UserService
}

func NewUserHandler(s *services.UserService) *UserHandler {
	return &UserHandler{Service: s}
}

type CreateUserRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var req CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request body",
		})
	}

	user := models.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
	}

	if err := h.Service.CreateUser(&user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User created successfully",
		"user":    user,
	})
}

// GetAllUsers godoc
// @Summary Get all users
// @Description Retrieve all users from the system
// @Tags users
// @Produce  json
// @Success 200 {array} models.User
// @Router /users [get]
func (h *UserHandler) GetAllUsers(c *fiber.Ctx) error {
	users, err := h.Service.GetAllUsers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve users",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"users": users,
	})
}

// GetUserByID godoc
// @Summary Get user by ID
// @Description Retrieve a user by their ID
// @Tags users
// @Produce  json
// @Param   id  path  int  true  "User ID"
// @Success 200 {object} models.User
// @Failure 404 {object} map[string]string
// @Router /users/{id} [get]
func (h *UserHandler) GetUserByID(c *fiber.Ctx) error {
	id, err := utils.ParseID(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid User",
		})
	}
	user, err := h.Service.GetUserByID(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve user",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"user": user,
	})
}

// UpdateUser godoc
// @Summary Update a user
// @Description Update an existing user's information
// @Tags users
// @Accept  json
// @Produce  json
// @Param   user  body  models.User  true  "User Data"
// @Success 200 {object} models.User
// @Failure 404 {object} map[string]string
// @Router /users [put]
func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request body",
		})
	}
	if err := h.Service.UpdateUser(&user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User updated successfully",
		"user":    user,
	})
}

func (h *UserHandler) Protected(c *fiber.Ctx) error {
	userToken := c.Locals("user")
	if userToken == nil {
		return c.Status(fiber.StatusUnauthorized).SendString("missing token")
	}
	token := userToken.(*jwt.Token) // from jwt middleware
	claims := token.Claims.(jwt.MapClaims)
	sub := claims["sub"]
	name, _ := claims["name"].(string)
	return c.SendString(fmt.Sprintf("hello user %v (name: %s)", sub, name))
}

type SignupRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

// Signup godoc
// @Summary Register a new user
// @Description Create a new user account, hash the password, and return access and refresh tokens
// @Tags auth
// @Accept  json
// @Produce  json
// @Param   user  body  SignupRequest  true  "User Signup Data"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /signup [post]
func (h *UserHandler) Signup(c *fiber.Ctx) error {
	var req SignupRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	// hash password
	hashed, err := utils.HashPassword(req.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to hash password"})
	}

	user := models.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  hashed,
	}

	if err := h.Service.CreateUser(&user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to create user"})
	}

	// generate tokens
	accessToken, refreshToken, err := utils.GenerateTokens(user.ID, user.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to generate tokens"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"user ID":       user.ID,
		"message":       "user created successfully",
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}
