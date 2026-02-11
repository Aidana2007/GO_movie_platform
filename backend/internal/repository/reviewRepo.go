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

type ReviewRepository struct {
	collection *mongo.Collection
}

func NewReviewRepository(db *mongo.Database) *ReviewRepository {
	return &ReviewRepository{
		collection: db.Collection("reviews"),
	}
}

func (r *ReviewRepository) Create(review *model.Review) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	review.CreatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, review)
	if err != nil {
		return err
	}

	review.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *ReviewRepository) FindByMovieID(movieID primitive.ObjectID) ([]*model.Review, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	opts := options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}})
	cursor, err := r.collection.Find(ctx, bson.M{"movie": movieID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var reviews []*model.Review
	if err = cursor.All(ctx, &reviews); err != nil {
		return nil, err
	}

	return reviews, nil
}

func (r *ReviewRepository) FindByID(id primitive.ObjectID) (*model.Review, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var review model.Review
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&review)
	if err != nil {
		return nil, err
	}

	return &review, nil
}

func (r *ReviewRepository) Delete(id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *ReviewRepository) GetAverageRating(movieID primitive.ObjectID) (float64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pipeline := []bson.M{
		{"$match": bson.M{"movie": movieID}},
		{"$group": bson.M{
			"_id": nil,
			"avg": bson.M{"$avg": "$rating"},
		}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)

	var result []bson.M
	if err = cursor.All(ctx, &result); err != nil {
		return 0, err
	}

	if len(result) == 0 {
		return 0, nil
	}

	avg, ok := result[0]["avg"].(float64)
	if !ok {
		return 0, nil
	}

	return avg, nil
}
