package utils

import "golang.org/x/crypto/bcrypt"

func HashRefreshToken(refreshToken string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(refreshToken), 10)
	return string(bytes), err
}

func CheckRefreshTokenHash(refreshToken, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(refreshToken))
	return err == nil
}
