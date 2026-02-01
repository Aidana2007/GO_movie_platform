package models

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Movie struct {
	ID          bson.ObjectID `json:"id" bson:"_id,omitempty" `
	Title       string        `json:"title" bson:"title" validate:"required,min=2,max=100"`
	Description string        `json:"description" bson:"description" validate:"required,min=2,max=1000"`
	GenreIDs    []int         `bson:"genre_ids" json:"genre_ids" validate:"required"`
	Year        int           `json:"year" bson:"year" validate:"required"`
	PosterPath  string        `json:"poster_path" bson:"poster_path" validate:"required,url"`
	VideoPath   string        `json:"video_path" bson:"video_path" validate:"required,url"`
}
