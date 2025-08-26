package handlers

//parse request body to proper format
//call respective service
//return response
import (
	"first_task/go-fiber-api/internal/models"
	"first_task/go-fiber-api/internal/services"
	utils "first_task/go-fiber-api/pkg"

	"github.com/gofiber/fiber/v2"
)

func NewBookHandler(s *services.BookService) *BookHandler {
	return &BookHandler{Service: s}
}

type BookHandler struct {
	Service *services.BookService
}

// @Summary Create a new book
// @Description Add a new book to the store
// @Tags books
// @Accept  json
// @Produce  json
// @Param   book  body  models.Book  true  "Book Data"
// @Success 201 {object} models.Book
// @Router /books [post]
func (B *BookHandler) CreateBook(c *fiber.Ctx) error {
	//fiber.ctx corresponds the the request
	var book models.Book
	//step one:parsing
	if err := c.BodyParser(&book); err != nil { //Take the incoming request body and fill in the fields of my book struct.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{ //.JSON() sends a JSON response back to the client.
			"error": "Failed to parse request body",
		}) //fiber.Map is just a shorthand for map[string]interface{} map[KeyType]ValueType
	}
	//2 service call
	if err := B.Service.CreateBook(&book); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create book",
		})
	}
	//time for a response
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Book created successfully",
		"book":    book,
	})
}

// @Summary Get all books
// @Description Retrieve all books from the store
// @Tags books
// @Produce  json
// @Success 200 {array} models.Book
// @Router /books [get]
func (B *BookHandler) GetAllBooks(c *fiber.Ctx) error {
	books, err := B.Service.GetAllBooks()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve books",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"books": books,
	})
}

// @Summary Get book by ID
// @Description Retrieve a book by its ID
// @Tags books
// @Produce  json
// @Param   id  path  int  true  "Book ID"
// @Success 200 {object} models.Book
// @Failure 404 {object} map[string]string
// @Router /books/{id} [get]
func (B *BookHandler) GetBookByID(c *fiber.Ctx) error {
	id, err := utils.ParseID(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid Book",
		})
	}
	book, err := B.Service.GetBookByID(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve book",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"book": book,
	})
}

// @Summary Checkin a book
// @Description Increase the quantity of a book by 1
// @Tags books
// @Produce  json
// @Param   id  path  int  true  "Book ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /books/{id}/checkin [post]
func (B *BookHandler) Checkin(c *fiber.Ctx) error {
	id, err := utils.ParseID(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid Book",
		})
	}
	if err := B.Service.Checkin(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to checkin book",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Book checked in successfully",
	})
}

// @Summary Checkout a book
// @Description Decrease the quantity of a book by 1
// @Tags books
// @Produce  json
// @Param   id  path  int  true  "Book ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /books/{id}/checkout [post]
func (B *BookHandler) Checkout(c *fiber.Ctx) error {
	id, err := utils.ParseID(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid Book",
		})
	}
	if err := B.Service.Checkout(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to checkout book",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Book checked out successfully",
	})
}
