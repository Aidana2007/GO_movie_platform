package repository

import (
	"context"

	"github.com/Aidana2007/GO_movie_platform/internal/model"

	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type FriendRepository struct {
	collection *mongo.Collection
}

func NewFriendRepository(db *mongo.Database) *FriendRepository {
	return &FriendRepository{
		collection: db.Collection("friendRequests"),
	}
}

func (r *FriendRepository) CreateRequest(fromUser, toUser primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	existing := r.collection.FindOne(ctx, bson.M{
		"fromUser": fromUser,
		"toUser":   toUser,
		"status":   "pending",
	})
	if existing.Err() == nil {
		return nil // Request already exists
	}

	request := &model.FriendRequest{
		FromUser:  fromUser,
		ToUser:    toUser,
		Status:    "pending",
		CreatedAt: time.Now(),
	}

	_, err := r.collection.InsertOne(ctx, request)
	return err
}

func (r *FriendRepository) GetPendingRequests(userID primitive.ObjectID) ([]*model.FriendRequest, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.M{
		"toUser": userID,
		"status": "pending",
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var requests []*model.FriendRequest
	if err = cursor.All(ctx, &requests); err != nil {
		return nil, err
	}

	return requests, nil
}

func (r *FriendRepository) UpdateRequestStatus(requestID primitive.ObjectID, status string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": requestID},
		bson.M{"$set": bson.M{"status": status}},
	)
	return err
}

func (r *FriendRepository) GetRequestByID(requestID primitive.ObjectID) (*model.FriendRequest, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var request model.FriendRequest
	err := r.collection.FindOne(ctx, bson.M{"_id": requestID}).Decode(&request)
	if err != nil {
		return nil, err
	}

	return &request, nil
}

func (r *FriendRepository) CheckPendingRequest(fromUser, toUser primitive.ObjectID) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	count, err := r.collection.CountDocuments(ctx, bson.M{
		"$or": []bson.M{
			{"fromUser": fromUser, "toUser": toUser, "status": "pending"},
			{"fromUser": toUser, "toUser": fromUser, "status": "pending"},
		},
	})

	return count > 0, err
}

func (r *FriendRepository) GetSentRequests(userID primitive.ObjectID) ([]*model.FriendRequest, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.M{
		"fromUser": userID,
		"status":   "pending",
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var requests []*model.FriendRequest
	if err = cursor.All(ctx, &requests); err != nil {
		return nil, err
	}

	return requests, nil
}
