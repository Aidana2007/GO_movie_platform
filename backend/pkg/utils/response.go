package utils

import (
	"encoding/json"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ContextKey string

const (
	UserIDKey ContextKey = "user_id"
	EmailKey  ContextKey = "email"
)

type Claims struct {
	UserID primitive.ObjectID
	Email  string
}

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func RespondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func RespondError(w http.ResponseWriter, status int, message string) {
	RespondJSON(w, status, Response{
		Success: false,
		Error:   message,
	})
}

func RespondSuccess(w http.ResponseWriter, status int, data interface{}) {
	RespondJSON(w, status, Response{
		Success: true,
		Data:    data,
	})
}

func GetUserIDFromContext(r *http.Request) (primitive.ObjectID, bool) {
	userID, ok := r.Context().Value(UserIDKey).(primitive.ObjectID)
	return userID, ok
}
