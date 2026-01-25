package models

type User struct {
	ID           string `json:"id" bson:"_id,omitempty"`
	Email        string `json:"email" bson:"email"`
	PasswordHash string `json:"passwordHash" bson:"passwordHash"`
	Role         string `json:"role" bson:"role"`
}
