package models

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Review struct {
	ID        bson.ObjectID `json:"id" bson:"_id,omitempty"`
	MovieID   bson.ObjectID `json:"movieId" bson:"movie_id" validate:"required"`
	UserID    string        `json:"userId" bson:"user_id" validate:"required"`
	Rating    int           `json:"rating" bson:"rating" validate:"required,min=1,max=5"`
	Comment   string        `json:"comment" bson:"comment" validate:"required,min=2,max=1000"`
	CreatedAt string        `json:"createdAt" bson:"created_at"`
}
