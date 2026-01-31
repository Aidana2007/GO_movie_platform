package models

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Genre struct {
	GenreID   int    `bson:"genre_id" json:"genre_id" validate:"required"`
	GenreName string `bson:"genre_name" json:"genre_name" validate:"required,min=2, max=100"`
}

type Movie struct {
	ID          bson.ObjectID `json:"id" bson:"_id,omitempty" `
	ImdbID      string        `json:"imdb_id" bson:"imdb_id" `
	Title       string        `json:"title" bson:"title" validate:"required, min=2, max=100"`
	Description string        `json:"description" bson:"description" validate:"required, min=2, max=1000"`
	Genre       []Genre       `json:"genre" bson:"genre" validate:"required, dive"`
	Year        int           `json:"year" bson:"year" validate:"required"`
	Ranking     int           `json:"ranking" bson:"ranking" validate:"required"`
	PosterPath  string        `json:"poster_path" bson:"poster_path" validate:"required, url"`
	VideoPath   string        `json:"video_path" bson:"video_path" validate:"required, url"`
	Review      string        `json:"review" bson:"review" validate:"required"`
}
