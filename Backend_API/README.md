# Bookstore Library API

A RESTful API built with **Go** and **Fiber** to manage a bookstore library and users, featuring JWT-based authentication, Swagger documentation, and MySQL database integration.

---

## Features

- **Books**
  - Create, Read, Update, Delete (CRUD) operations
  - Check-in and Check-out functionality
- **Users**
  - Signup and Login
  - Password hashing with bcrypt
  - JWT-based authentication (access + refresh tokens)
  - Update and retrieve users
- **Authentication**
  - JWT middleware to protect routes
  - Token generation and validation
- **Documentation**
  - Swagger integration for API documentation

---

## Tech Stack

- **Backend:** Go, Fiber
- **Database:** MySQL (via GORM)
- **Authentication:** JWT (HS256)
- **Password Hashing:** bcrypt
- **Documentation:** Swagger (Swaggo)

---

## Installation

1. Clone the repository:

```bash
git clone https://github.com/AliSleiman0/your-project.git
cd your-project

```
Install dependencies:
```bash
go mod tidy
```
Set up your .env file with the following:
```bash
JWT_SECRET=your_access_secret
JWT_REFRESH_SECRET=your_refresh_secret
JWT_TTL_HOURS=72
DB_USER=root
DB_PASS=password
DB_HOST=localhost
DB_PORT=3306
DB_NAME=bookstore
DB_CHARSET=utf8mb4
APP_PORT=3000
APP_ENV=development
```
Initialize the database:
```bash
go run main.go
```
This will automatically run GORM migrations for Book and User models.

## Running the API
```bash
go run main.go
```
Access the API at: http://localhost:3000

Swagger docs: http://localhost:3000/swagger/index.html

## API Endpoints
### Public

- POST /api/signup – Create a new user

- POST /api/login – Login and receive JWT tokens

- GET /health – Health check

### Protected (JWT required)
#### Books

- POST /api/books – Add a new book

- GET /api/books – List all books

- GET /api/books/:id – Get book by ID

- POST /api/books/:id/checkin – Check in a book

- POST /api/books/:id/checkout – Check out a book

#### Users

- GET /api/users – List all users

- GET /api/users/:id – Get user by ID

- PUT /api/users/:id – Update user info

## Authentication
- Use Bearer JWT tokens for protected endpoints.

- Access tokens expire according to JWT_TTL_HOURS.

- Refresh tokens can be used to generate new access tokens (future implementation).

## Password Security
- User passwords are hashed using bcrypt.

- Plain-text passwords are never stored.

## Swagger Documentation
- Access via /swagger/index.html after running the API.

- Provides interactive testing for all endpoints.

## Contributing
- Contributions are welcome! Feel free to open issues or submit pull requests.

## License
- Ali License