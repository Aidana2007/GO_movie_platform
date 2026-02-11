package service

import (
	"errors"

	"github.com/Aidana2007/GO_movie_platform/internal/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FriendService struct {
	friendRepo *repository.FriendRepository
	userRepo   *repository.UserRepository
}

func NewFriendService(friendRepo *repository.FriendRepository, userRepo *repository.UserRepository) *FriendService {
	return &FriendService{
		friendRepo: friendRepo,
		userRepo:   userRepo,
	}
}

func (s *FriendService) SendFriendRequest(fromUserID, toUserID primitive.ObjectID) error {
	if fromUserID == toUserID {
		return errors.New("cannot send friend request to yourself")
	}

	isFriend, err := s.userRepo.IsFriend(fromUserID, toUserID)
	if err != nil {
		return err
	}
	if isFriend {
		return errors.New("already friends")
	}

	hasPending, err := s.friendRepo.CheckPendingRequest(fromUserID, toUserID)
	if err != nil {
		return err
	}
	if hasPending {
		return errors.New("friend request already pending")
	}

	return s.friendRepo.CreateRequest(fromUserID, toUserID)
}

func (s *FriendService) GetPendingRequests(userID primitive.ObjectID) ([]map[string]interface{}, error) {
	requests, err := s.friendRepo.GetPendingRequests(userID)
	if err != nil {
		return nil, err
	}

	var result []map[string]interface{}
	for _, req := range requests {
		user, err := s.userRepo.FindByID(req.FromUser)
		if err != nil {
			continue
		}

		result = append(result, map[string]interface{}{
			"requestId": req.ID.Hex(),
			"username":  user.Username,
			"userId":    user.ID.Hex(),
		})
	}

	return result, nil
}

func (s *FriendService) AcceptFriendRequest(requestID primitive.ObjectID) error {
	request, err := s.friendRepo.GetRequestByID(requestID)
	if err != nil {
		return err
	}

	if request.Status != "pending" {
		return errors.New("request is not pending")
	}

	err = s.userRepo.AddFriend(request.FromUser, request.ToUser)
	if err != nil {
		return err
	}

	return s.friendRepo.UpdateRequestStatus(requestID, "accepted")
}

func (s *FriendService) RejectFriendRequest(requestID primitive.ObjectID) error {
	request, err := s.friendRepo.GetRequestByID(requestID)
	if err != nil {
		return err
	}

	if request.Status != "pending" {
		return errors.New("request is not pending")
	}

	return s.friendRepo.UpdateRequestStatus(requestID, "rejected")
}

func (s *FriendService) CheckPendingRequest(fromUser, toUser primitive.ObjectID) (bool, error) {
	return s.friendRepo.CheckPendingRequest(fromUser, toUser)
}

func (s *FriendService) GetSentRequests(userID primitive.ObjectID) ([]map[string]interface{}, error) {
	requests, err := s.friendRepo.GetSentRequests(userID)
	if err != nil {
		return nil, err
	}

	var result []map[string]interface{}
	for _, req := range requests {
		user, err := s.userRepo.FindByID(req.ToUser)
		if err != nil {
			continue
		}

		result = append(result, map[string]interface{}{
			"requestId": req.ID.Hex(),
			"username":  user.Username,
			"userId":    user.ID.Hex(),
		})
	}

	return result, nil
}

func (s *FriendService) CancelFriendRequest(requestID primitive.ObjectID) error {
	request, err := s.friendRepo.GetRequestByID(requestID)
	if err != nil {
		return err
	}

	if request.Status != "pending" {
		return errors.New("request is not pending")
	}

	return s.friendRepo.UpdateRequestStatus(requestID, "cancelled")
}
