package handler

import (
	"net/http"

	"github.com/Aidana2007/GO_movie_platform/internal/model"
	"github.com/Aidana2007/GO_movie_platform/internal/service"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FriendHandler struct {
	friendService *service.FriendService
	userService   *service.UserService
}

func NewFriendHandler(friendService *service.FriendService, userService *service.UserService) *FriendHandler {
	return &FriendHandler{
		friendService: friendService,
		userService:   userService,
	}
}

func (h *FriendHandler) SendFriendRequest(c *gin.Context) {
	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "Unauthorized"})
		return
	}

	var req model.SendFriendRequestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid request"})
		return
	}

	toUserID, err := primitive.ObjectIDFromHex(req.TargetUserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid user ID"})
		return
	}

	err = h.friendService.SendFriendRequest(userID.(primitive.ObjectID), toUserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Friend request sent"})
}

func (h *FriendHandler) GetFriendRequests(c *gin.Context) {
	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "Unauthorized"})
		return
	}

	requests, err := h.friendService.GetPendingRequests(userID.(primitive.ObjectID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Failed to get requests"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": requests})
}

func (h *FriendHandler) AcceptFriendRequest(c *gin.Context) {
	requestID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid request ID"})
		return
	}

	err = h.friendService.AcceptFriendRequest(requestID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Friend request accepted"})
}

func (h *FriendHandler) RejectFriendRequest(c *gin.Context) {
	requestID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid request ID"})
		return
	}

	err = h.friendService.RejectFriendRequest(requestID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Friend request rejected"})
}

func (h *FriendHandler) GetFriends(c *gin.Context) {
	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "Unauthorized"})
		return
	}

	friends, err := h.userService.GetFriends(userID.(primitive.ObjectID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Failed to get friends"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": friends})
}

func (h *FriendHandler) RemoveFriend(c *gin.Context) {
	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "Unauthorized"})
		return
	}

	friendID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid friend ID"})
		return
	}

	err = h.userService.RemoveFriend(userID.(primitive.ObjectID), friendID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Failed to remove friend"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Friend removed"})
}

func (h *FriendHandler) GetSentRequests(c *gin.Context) {
	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "Unauthorized"})
		return
	}

	requests, err := h.friendService.GetSentRequests(userID.(primitive.ObjectID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Failed to get sent requests"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": requests})
}

func (h *FriendHandler) CancelFriendRequest(c *gin.Context) {
	requestID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid request ID"})
		return
	}

	err = h.friendService.CancelFriendRequest(requestID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Friend request cancelled"})
}
