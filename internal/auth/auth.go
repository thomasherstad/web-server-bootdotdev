package auth

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(p string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func ComparePasswords(hashedPassword, attemptedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(attemptedPassword))
	if err != nil {
		return false
	}
	return true
}

func GenerateRefreshToken() (string, time.Time, error) {
	num := make([]byte, 32)
	_, err := rand.Read(num)
	if err != nil {
		return "", time.Time{}, err
	}

	//60 days
	expiration := time.Now().UTC().AddDate(0, 0, 60)

	return hex.EncodeToString(num), expiration, nil
}
