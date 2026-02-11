package repository

import (
	"context"

	"github.com/Aidana2007/GO_movie_platform/internal/model"

	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MovieRepository struct {
	collection *mongo.Collection
}

func NewMovieRepository(db *mongo.Database) *MovieRepository {
	return &MovieRepository{
		collection: db.Collection("movies"),
	}
}

func (r *MovieRepository) Create(movie *model.Movie) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	movie.CreatedAt = time.Now()
	movie.UpdatedAt = time.Now()
	movie.Reviews = []primitive.ObjectID{}
	movie.Ranking = 0

	result, err := r.collection.InsertOne(ctx, movie)
	if err != nil {
		return err
	}

	movie.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *MovieRepository) FindAll() ([]*model.Movie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	opts := options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}})
	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var movies []*model.Movie
	if err = cursor.All(ctx, &movies); err != nil {
		return nil, err
	}

	return movies, nil
}

func (r *MovieRepository) FindByID(id primitive.ObjectID) (*model.Movie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var movie model.Movie
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&movie)
	if err != nil {
		return nil, err
	}

	return &movie, nil
}

func (r *MovieRepository) FindByIDs(ids []primitive.ObjectID) ([]*model.Movie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.M{"_id": bson.M{"$in": ids}})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var movies []*model.Movie
	if err = cursor.All(ctx, &movies); err != nil {
		return nil, err
	}

	return movies, nil
}

func (r *MovieRepository) Update(id primitive.ObjectID, movie *model.Movie) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	movie.UpdatedAt = time.Now()

	update := bson.M{
		"$set": bson.M{
			"title":       movie.Title,
			"description": movie.Description,
			"year":        movie.Year,
			"director":    movie.Director,
			"cast":        movie.Cast,
			"genre":       movie.Genre,
			"posterUrl":   movie.PosterURL,
			"trailerUrl":  movie.TrailerURL,
			"updatedAt":   movie.UpdatedAt,
		},
	}

	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

func (r *MovieRepository) Delete(id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *MovieRepository) AddReview(movieID, reviewID primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": movieID},
		bson.M{
			"$addToSet": bson.M{"reviews": reviewID},
			"$set":      bson.M{"updatedAt": time.Now()},
		},
	)
	return err
}

func (r *MovieRepository) RemoveReview(movieID, reviewID primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": movieID},
		bson.M{
			"$pull": bson.M{"reviews": reviewID},
			"$set":  bson.M{"updatedAt": time.Now()},
		},
	)
	return err
}

func (r *MovieRepository) UpdateRanking(movieID primitive.ObjectID, ranking float64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": movieID},
		bson.M{
			"$set": bson.M{
				"ranking":   ranking,
				"updatedAt": time.Now(),
			},
		},
	)
	return err
}

func (r *MovieRepository) Search(query string, genre string, sort string, minRating float64) ([]*model.Movie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{}

	if query != "" {
		filter["$or"] = []bson.M{
			{"title": bson.M{"$regex": query, "$options": "i"}},
			{"director": bson.M{"$regex": query, "$options": "i"}},
			{"cast": bson.M{"$regex": query, "$options": "i"}},
		}
	}

	if genre != "" {
		filter["genre"] = genre
	}

	if minRating > 0 {
		filter["ranking"] = bson.M{"$gte": minRating}
	}

	// Determine sort order
	var sortField string
	var sortOrder int = -1 // descending by default

	switch sort {
	case "rating":
		sortField = "ranking"
	case "year":
		sortField = "year"
	case "latest":
		sortField = "createdAt"
	default:
		sortField = "createdAt"
	}

	opts := options.Find().SetSort(bson.D{{Key: sortField, Value: sortOrder}})
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var movies []*model.Movie
	if err = cursor.All(ctx, &movies); err != nil {
		return nil, err
	}

	return movies, nil
}

func (r *MovieRepository) GetTopRated(limit int) ([]*model.Movie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	opts := options.Find().
		SetSort(bson.D{{Key: "ranking", Value: -1}}).
		SetLimit(int64(limit))

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var movies []*model.Movie
	if err = cursor.All(ctx, &movies); err != nil {
		return nil, err
	}

	return movies, nil
}
