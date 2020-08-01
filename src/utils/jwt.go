package utils

import (
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

var jwtKey = []byte("my_secret_key")

type Claims struct {
	UserId string `json:"userid"`
	// время истечения токена
	jwt.StandardClaims
}

func CreateTokenAndRefreshToken(c *Claims, key []byte) (string, string, error) {
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS512, c)
	gen := uuid.New()
	refreshToken := gen.String()
	tokenString, err := newToken.SignedString(key)

	return tokenString, refreshToken, err
}

func DecodeTokenAndGetUserId(reqToken string) (string, error) {
	claims := Claims{}
	_, err := jwt.ParseWithClaims(reqToken, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	})
	return claims.UserId, err
}

func GetHeaderToken(w http.ResponseWriter, r *http.Request) string {
	w.Header().Set("Content-Type", "application/json")
	reqToken := r.Header.Get("Authorization")
	token := strings.Split(reqToken, " ")[1]
	return token
}

func VerifyToken(reqToken string, key []byte) bool {
	token, err := jwt.Parse(reqToken, func(t *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})
	if err == nil && token.Valid {
		return true
	} else {
		return false
	}
}
