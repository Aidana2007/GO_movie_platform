Here is the revised README file in English, without emojis or decorative elements:

Movie Platform

A web platform for browsing and discussing movies built with Go (Gin framework) and HTML/JavaScript.


Project Structure
GO_movie_platform/
├── main.go                 # Application entry point
├── controllers/            # Request handlers
│   ├── home.go             # Home page (Go template)
│   ├── movies.go           # HTML + JS pages
│   ├── movie.go            # Movies API
│   ├── review.go           # Reviews API
│   └── user.go             # Users API
├── middleware/             # Middleware layer
│   ├── auth.go             # Authentication middleware
│   └── admin.go            # Admin authorization middleware
├── models/                 # Data models
├── routes/                 # Route definitions
├── database/               # Database connection
├── utils/                  # Utility functions
├── templates/              # Go HTML templates (HTML + JS)
│   ├── index.html          # Home page (Go template)
│   ├── movies.html         # Movies list page
│   └── movie-details.html  # Movie details page
└── .env                    # Environment variables
Installation

Go 1.21 or higher

MongoDB

Setup Steps

Install Go dependencies:

go mod download

Create a .env file in the project root:

MONGODB_URI=mongodb://localhost:27017
MONGODB_DATABASE=movie_platform
SECRET_KEY=your-secret-key-here
SECRET_REFRESH_KEY=your-refresh-secret-key-here
PORT=8080
Running the Application
go run main.go

The server will start at:

http://localhost:8080
Available URLs

Home (Go Template): http://localhost:8080/

Movies Page: http://localhost:8080/movies

API Base: http://localhost:8080/api

API Endpoints
Public Endpoints

GET /api/movies — Get list of movies

GET /api/movies/:id — Get movie details

GET /api/movies/:id/reviews — Get reviews for a movie

POST /api/register — User registration

POST /api/login — User login

Requires Authentication

POST /api/movies/:id/reviews — Add a review

DELETE /api/reviews/:id — Delete a review

Requires Admin Role

POST /api/movies — Create a new movie

PUT /api/movies/:id — Update a movie

DELETE /api/movies/:id — Delete a movie

Technologies Used

Backend:

Go 1.21+

Gin (web framework)

MongoDB

JWT for authentication

Frontend:

Go HTML Templates

Vanilla JavaScript (Fetch API)

Vanilla CSS


Server-side rendering with Go templates

JWT-based authentication and role-based authorization
