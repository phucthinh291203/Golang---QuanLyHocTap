package auth

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var SecretKey = []byte(os.Getenv("secretKey"))

type BaseClaims struct {
	Username string             `json:"username"`
	Title    string             `json:"title"`
	Name     string             `json:"name"`
	UserID   primitive.ObjectID `json:"user_id"`
	jwt.StandardClaims
}

func GenerateJWT(claims BaseClaims) (string, error) {
	expirationTime := time.Now().Add(5 * time.Hour)
	claims.ExpiresAt = expirationTime.Unix() // Thiết lập thời gian hết hạn

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(SecretKey)

}

func ParseJWT(tokenString string) (*BaseClaims, error) {
	claims := &BaseClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return SecretKey, nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}
	return claims, nil
}
