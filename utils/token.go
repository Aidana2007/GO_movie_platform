package utils

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/Aidana2007/GO_movie_platform/database"
	"github.com/gin-gonic/gin"
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

var SECRET_KEY string = os.Getenv("SECRET_KEY")
var SECRET_REFRESH_KEY string = os.Getenv("SECRET_REFRESH_KEY")
var userCollection = database.GetCollection("users")

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
	signedToken, err := token.SignedString([]byte(SECRET_KEY))

	if err != nil {
		return "", "", err
	}

	refreshClaims := &SignedDetails{
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
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	resfreshSignedToken, err := refreshToken.SignedString([]byte(SECRET_REFRESH_KEY))

	if err != nil {
		return "", "", err
	}

	return signedToken, resfreshSignedToken, nil
}

func UpdateAllTokens(userId, token, refreshToken string) (err error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	updateAt, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	updatedata := bson.M{
		"$set": bson.M{
			"token":         token,
			"refresh_token": refreshToken,
			"updated_at":    updateAt,
		},
	}
	_, err = userCollection.UpdateOne(ctx, bson.M{"userId": userId}, updatedata)

	if err != nil {
		return err
	}
	return nil
}

func GetAccessToken(c *gin.Context) (string, error) {
	authHeader := c.Request.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("no auth header")
	}
	tokenString := authHeader[len("Bearer "):]

	if tokenString == "" {
		return "", errors.New("bearer token not found")
	}
	return tokenString, nil
}
func ValidateToken(tokenString string) (*SignedDetails, error) {
	claims := &SignedDetails{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRET_KEY), nil
	})
	if err != nil {
		return nil, err
	}

	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, err
	}

	if claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, errors.New("token is expired")
	}
	return claims, nil
}
