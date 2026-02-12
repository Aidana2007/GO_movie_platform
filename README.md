
# Movie Platform

A web platform for browsing and discussing movies built with Go (Gin framework) and HTML/JavaScript.

## Project Requirements

- Usage of middleware: AuthMiddleware, AdminMiddleware, and CORS middleware
- Frontend implemented with Go HTML templates (located in templates/)
- Usage of Go templating tags such as {{.TotalMovies}}, {{.TotalUsers}}, {{.TotalReviews}}
- Clean code and well-structured project architecture
- Bonus: Embedded JavaScript inside Go templates

## Project Structure

GO_movie_platform/
├── main.go
├── controllers/
├── middleware/
├── models/
├── routes/
├── database/
├── utils/
├── templates/
└── .env

## Installation

Requirements:
- Go 1.21+
- MongoDB

Setup:
1. go mod download
2. Create .env file with MONGODB_URI, MONGODB_DATABASE, SECRET_KEY, SECRET_REFRESH_KEY, PORT

## Run

go run main.go

Server runs at http://localhost:8080
