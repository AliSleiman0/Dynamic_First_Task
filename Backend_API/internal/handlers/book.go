package handlers

//parse request body to proper format
//call respective service
//return response
import (
	"first_task/go-fiber-api/internal/models"
	"first_task/go-fiber-api/internal/services"
	utils "first_task/go-fiber-api/pkg"
	"fmt"
	"html"
	"strconv"
	"strings"

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
//
//	func (B *BookHandler) GetAllBooks(c *fiber.Ctx) error {
//		books, err := B.Service.GetAllBooks()
//		if err != nil {
//			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
//				"error": "Failed to retrieve books",
//			})
//		}
//		return c.Status(fiber.StatusOK).JSON(fiber.Map{
//			"books": books,
//		})
//	}
func (B *BookHandler) GetAllBooks(c *fiber.Ctx) error {
	search := c.Query("search", "")

	books, err := B.Service.GetAllBooksFiltered(search)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve books",
		})
	}

	// if user explicitly wants JSON: /api/books?format=json
	if c.Query("format") == "json" {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"books": books,
		})
	}

	// build a simple Bootstrap grid of cards
	var sb strings.Builder
	sb.WriteString(`<div class="row g-3">`)

	if len(books) == 0 {
		sb.WriteString(`<div class="col-12"><div class="alert alert-warning mb-0">No books found.</div></div>`)
	} else {
		for _, book := range books {
			// fallback image
			img := book.Img_url
			if img == "" {
				img = "https://via.placeholder.com/150x220?text=No+Cover"
			}

			// escape user-provided strings to avoid injecting HTML
			title := html.EscapeString(book.Title)
			genre := html.EscapeString(book.Genre)

			sb.WriteString(fmt.Sprintf(`
                <div class="col-12 col-sm-6 col-md-4 col-lg-3">
                  <div class="card h-100">
                    <img src="%s" class="card-img-top" alt="%s cover" style="height:220px; object-fit:cover;">
                    <div class="card-body d-flex flex-column">
                      <h5 class="card-title">%s</h5>
                      <p class="card-text mb-1"><small class="text-muted">Genre: %s</small></p>
                      <p class="card-text mb-2"><small class="text-muted">Published: %d • Qty: %d</small></p>
                      <div class="mt-auto">
                        <button class="btn btn-sm btn-primary" 
								type="button" 
								hx-get="/books/%d/details" 
								hx-target="#modalBody" 
								hx-swap="innerHTML"
								data-bs-toggle="modal" 
								data-bs-target="#bookModal">
						Details
						</button>
                      </div>
                    </div>
                  </div>
                </div>
            `, img, title, title, genre, book.PublishedYear, book.Quantity, book.ID))
		}
	}

	sb.WriteString(`</div>`) // close row

	// Return HTML fragment (suitable for htmx hx-get)
	return c.Status(fiber.StatusOK).Type("html").SendString(sb.String())
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
func (B *BookHandler) GetBookDetails(c *fiber.Ctx) error {
	// Parse ID param
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid book ID")
	}

	// Call service
	book, err := B.Service.GetBookByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).SendString("Book not found")
	}

	// fallback image
	img := book.Img_url
	if img == "" {
		img = "https://via.placeholder.com/300x400?text=No+Cover"
	}

	// checkin/checkout logic (for button label)
	actionLabel := "Check Out"

	html := fmt.Sprintf(`
        <div class="d-flex flex-column align-items-center">
            <img src="%s" alt="%s cover" class="img-fluid mb-3" style="max-height:400px; object-fit:cover;">
            <h3>%s</h3>
            <p class="text-muted">Genre: %s</p>
            <p>Published: %d</p>
            <p>Quantity: %d</p>
            <button class="btn btn-success" 
                    hx-post="/books/%d/toggle" 
                    hx-target="#modalBody" 
                    hx-swap="innerHTML">
                %s
            </button>
        </div>
    `, img, book.Title, book.Title, book.Genre, book.PublishedYear, book.Quantity, book.ID, actionLabel)

	return c.Type("html").SendString(html)
}
