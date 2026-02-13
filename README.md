# MovieLand – Movie Platform

MovieLand is a full-stack movie platform built with Go (Golang) for the backend and HTML/CSS/JavaScript for the frontend.  
The platform allows users to browse movies, view details, leave reviews, manage watchlists, add friends, and provides an admin panel for management.

---

## Tech Stack

### Backend
- Go (Golang)
- net/http
- MongoDB
- JWT authentication
- Middleware (Auth + RBAC)
- Clean Architecture (Handler → Service → Repository)

### Frontend
- HTML templates
- CSS
- Vanilla JavaScript
- Server-side rendering

---

## Project Structure

```
GO_movie_platform/
│
├── backend/
│   │
│   ├── cmd/
│   │   └── server/
│   │       └── main.go                 # Application entry point
│   │
│   ├── internal/
│   │   ├── config/
│   │   │   └── config.go               # Environment configuration
│   │   │
│   │   ├── handler/                    # HTTP handlers (controllers)
│   │   │   ├── adminHandler.go
│   │   │   ├── authHandler.go
│   │   │   ├── friendHandler.go
│   │   │   ├── movieHandler.go
│   │   │   ├── pageHandler.go
│   │   │   ├── reviewHandler.go
│   │   │   └── userHandler.go
│   │   │
│   │   ├── middleware/
│   │   │   ├── authMiddleware.go
│   │   │   └── rbaciddleware.go
│   │   │
│   │   ├── model/                      # Database models
│   │   │   └── models.go
│   │   │
│   │   ├── repository/                 # Database access layer
│   │   │   ├── friendRepo.go
│   │   │   ├── movieRepo.go
│   │   │   ├── reviewRepo.go
│   │   │   └── userRepo.go
│   │   │
│   │   └── service/                    # Business logic layer
│   │       ├── authService.go
│   │       ├── friendService.go
│   │       ├── movieService.go
│   │       ├── reviewService.go
│   │       └── userService.go
│   │   
│   ├── pkg/
│   │   └── response.go             
│   │
│   ├── go.mod
│   ├── go.sum
│   └── .env                            # Environment variables 
│
├── frontend/
│   │
│   ├── static/
│   │   ├── css/
│   │   │   ├── admin.css
│   │   │   └── style.css
│   │   │
│   │   └── js/
│   │       ├── admin.js
│   │       ├── auth.js
│   │       ├── main.js
│   │       ├── modal.js
│   │       ├── movie-details.js
│   │       ├── profile.js
│   │       ├── users.js
│   │       └── watchlist.js
│   │
│   └── templates/
│       ├── admin_panel.html
│       ├── app.html
│       ├── base.html
│       ├── error.html
│       ├── home.html
│       ├── login.html
│       ├── movie_details.html
│       ├── movies.html
│       ├── profile.html
│       ├── register.html
│       ├── user-profile.html
│       ├── users.html
│       └── watchlist.html
│
└── README.md
```

---

## Features

### Authentication
- User registration
- User login
- JWT-based authentication
- Protected routes
- Role-based access control (RBAC)

### User Features
- View all movies
- View movie details
- Add/remove movies to watchlist
- Leave reviews
- Add friends
- View user profiles
- Personal profile management

### Admin Features
- Admin panel
- Manage users
- Manage movies
- Moderate content

---

## Architecture

The backend follows a layered architecture:

```
Handler → Service → Repository → Database
```

- **Handler**: Handles HTTP requests and responses
- **Service**: Contains business logic
- **Repository**: Communicates with MongoDB
- **Middleware**: Handles authentication and authorization
- **Model**: Defines database schemas

This separation improves maintainability, scalability, and testing.

---

## Database

The project uses MongoDB.

Main collections:
- users
- movies
- reviews
- friends
- watchlists

---

## Environment Variables

Create a `.env` file inside the `backend` directory or configure environment variables manually.

Example:

```
PORT=8080
MONGO_URI=mongodb://localhost:27017
DB_NAME=movieland
JWT_SECRET=your_secret_key
SECRET_REFRESH_KEY=your_secret_refresh_key
```

Adjust values according to your MongoDB configuration.

---

## Installation and Setup

### 1. Clone the repository

```
git clone https://github.com/Aidana2007/GO_movie_platform.git
cd GO_movie_platform/backend
```

### 2. Install dependencies

```
go mod tidy
```

### 3. Configure environment variables

Set up your MongoDB connection and JWT secret.

### 4. Run the server

```
go run cmd/server/main.go
```

Server will start on:

```
http://localhost:8080
```

---

## Frontend

The frontend uses server-side rendered HTML templates located in:

```
frontend/templates
```

Static files:
- CSS → `frontend/static/css`
- JS → `frontend/static/js`

---

## Main Endpoints (Overview)

### Authentication
- `POST /register`
- `POST /login`

### Movies
- `GET /movies`
- `GET /movies/{id}`
- `POST /movies` (admin only)

### Reviews
- `POST /reviews`
- `DELETE /reviews/{id}`

### Users
- `GET /profile`
- `GET /users`
- `GET /users/{id}`

### Watchlist
- Add to watchlist
- Remove from watchlist
- View watchlist

Access to certain routes is restricted by JWT and role middleware.

---

## Middleware

- **AuthMiddleware** – Validates JWT token
- **RBACMiddleware** – Restricts access based on user role

---

## Security

- Password hashing
- JWT authentication
- Role-based access control
- Protected routes

---
