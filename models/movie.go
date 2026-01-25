package models

type Movie struct {
	ID          string `json:"id" bson:"_id,omitempty"`
	Title       string `json:"title" bson:"title"`
	Description string `json:"description" bson:"description"`
	Genre       string `json:"genre" bson:"genre"`
	Year        int    `json:"year" bson:"year"`
	VideoURL    string `json:"videoUrl" bson:"videoUrl"`
	Review      string `json:"review" bson:"review"`
}
