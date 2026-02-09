package utils

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/Aidana2007/GO_movie_platform/database"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type SignedDetails struct {
	Email     string
	FirstName string
	LastName  string
	Role      string
	UserId    string
	jwt.RegisteredClaims
}

func getSecretKey() string {
	key := os.Getenv("SECRET_KEY")
	if key == "" {
		panic("SECRET_KEY not set")
	}
	return key
}

func getRefreshSecretKey() string {
	key := os.Getenv("SECRET_REFRESH_KEY")
	if key == "" {
		panic("SECRET_REFRESH_KEY not set")
	}
	return key
}

func GenerateAllTokens(email, firstname, lastname, role, userId string) (string, string, error) {
	claims := &SignedDetails{
		Email:     email,
		FirstName: firstname,
		LastName:  lastname,
		Role:      role,
		UserId:    userId,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "movie_land",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(getSecretKey()))
	if err != nil {
		return "", "", err
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	refreshSignedToken, err := refreshToken.SignedString([]byte(getRefreshSecretKey()))
	if err != nil {
		return "", "", err
	}

	return signedToken, refreshSignedToken, nil
}

func UpdateAllTokens(userId, token, refreshToken string) error {
	userCollection := database.GetCollection("users")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	updateData := bson.M{
		"$set": bson.M{
			"token":         token,
			"refreshToken":  refreshToken,
			"updated_at":    time.Now(),
		},
	}

	_, err := userCollection.UpdateOne(ctx, bson.M{"userId": userId}, updateData)
	return err
}

func ValidateToken(tokenString string) (*SignedDetails, error) {
	claims := &SignedDetails{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(getSecretKey()), nil
	})
	if err != nil {
		return nil, err
	}

	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, errors.New("unexpected signing method")
	}

	if claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, errors.New("token is expired")
	}

	return claims, nil
}
