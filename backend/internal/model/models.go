package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID   `bson:"_id,omitempty" json:"_id,omitempty"`
	Username  string               `bson:"username" json:"username"`
	Email     string               `bson:"email" json:"email"`
	Password  string               `bson:"password" json:"password,omitempty"`
	Role      string               `bson:"role" json:"role"`
	Watchlist []primitive.ObjectID `bson:"watchlist" json:"watchlist"`
	Friends   []primitive.ObjectID `bson:"friends" json:"friends"`
	CreatedAt time.Time            `bson:"createdAt" json:"createdAt"`
}

type Movie struct {
	ID          primitive.ObjectID   `bson:"_id,omitempty" json:"_id,omitempty"`
	Title       string               `bson:"title" json:"title"`
	Description string               `bson:"description" json:"description"`
	Year        int                  `bson:"year" json:"year"`
	Director    string               `bson:"director" json:"director"`
	Cast        []string             `bson:"cast" json:"cast"`
	Ranking     float64              `bson:"ranking" json:"ranking"`
	Genre       []string             `bson:"genre" json:"genre"`
	PosterURL   string               `bson:"posterUrl" json:"posterUrl"`
	TrailerURL  string               `bson:"trailerUrl" json:"trailerUrl"`
	CreatedBy   primitive.ObjectID   `bson:"createdBy" json:"createdBy"`
	Reviews     []primitive.ObjectID `bson:"reviews" json:"reviews"`
	CreatedAt   time.Time            `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time            `bson:"updatedAt" json:"updatedAt"`
}

type Review struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	User      primitive.ObjectID `bson:"user" json:"user"`
	Movie     primitive.ObjectID `bson:"movie" json:"movie"`
	Username  string             `bson:"username" json:"username"`
	Rating    int                `bson:"rating" json:"rating"`
	Comment   string             `bson:"comment" json:"comment"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
}

type FriendRequest struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	FromUser  primitive.ObjectID `bson:"fromUser" json:"fromUser"`
	ToUser    primitive.ObjectID `bson:"toUser" json:"toUser"`
	Status    string             `bson:"status" json:"status"` // pending, accepted, rejected
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type CreateMovieRequest struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Year        int      `json:"year"`
	Director    string   `json:"director"`
	Cast        []string `json:"cast"`
	Genre       []string `json:"genre"`
	PosterURL   string   `json:"posterUrl"`
	TrailerURL  string   `json:"trailerUrl"`
}

type CreateReviewRequest struct {
	Rating  int    `json:"rating"`
	Comment string `json:"comment"`
}

type SendFriendRequestRequest struct {
	TargetUserID string `json:"targetUserId"`
}
