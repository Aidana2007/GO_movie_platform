package models

import "go.mongodb.org/mongo-driver/v2/bson"

type User struct {
	ID           bson.ObjectID `json:"id" bson:"_id,omitempty"`
	UserId       string `json:"userId" bson:"userId"`
	FirstName    string `json:"firstName" bson:"firstName" validate:"required,min=2,max=100"`
	LastName     string `json:"lastName" bson:"lastName" validate:"required,min=2,max=100"`
	Email        string `json:"email" bson:"email" validate:"required,email"`
	Password     string `json:"password" bson:"password" validate:"required,min=6,max=32"`
	Role         string `json:"role" bson:"role" validate:"omitempty,oneof=ADMIN USER"`
	Token        string `json:"token" bson:"token"`
	RefreshToken string `json:"refreshToken" bson:"refreshToken"`
}

type UserLogin struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=32"`
}
type UserResponse struct {
	UserID       string `json:"userId"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	Email        string `json:"email"`
	Role         string `json:"role"`
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
}
