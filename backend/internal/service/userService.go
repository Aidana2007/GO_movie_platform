package service

import (
	"github.com/Aidana2007/GO_movie_platform/internal/model"
	"github.com/Aidana2007/GO_movie_platform/internal/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) GetUserByID(userID primitive.ObjectID) (*model.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	user.Password = ""
	return user, nil
}

func (s *UserService) AddToWatchlist(userID, movieID primitive.ObjectID) error {
	return s.userRepo.AddToWatchlist(userID, movieID)
}

func (s *UserService) RemoveFromWatchlist(userID, movieID primitive.ObjectID) error {
	return s.userRepo.RemoveFromWatchlist(userID, movieID)
}

func (s *UserService) GetWatchlist(userID primitive.ObjectID) ([]primitive.ObjectID, error) {
	return s.userRepo.GetWatchlist(userID)
}

func (s *UserService) GetFriends(userID primitive.ObjectID) ([]*model.User, error) {
	return s.userRepo.GetFriends(userID)
}

func (s *UserService) SearchUsers(query string, limit int) ([]*model.User, error) {
	return s.userRepo.SearchUsers(query, limit)
}

func (s *UserService) RemoveFriend(userID, friendID primitive.ObjectID) error {
	return s.userRepo.RemoveFriend(userID, friendID)
}

func (s *UserService) IsFriend(userID, friendID primitive.ObjectID) (bool, error) {
	return s.userRepo.IsFriend(userID, friendID)
}
