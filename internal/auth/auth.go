package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/http"
	"strings"
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

func GetApiKey(req http.Header) (string, error) {
	keyString := req.Get("Authorization")

	keys := strings.Fields(keyString)
	if len(keys) != 2 {
		return "", errors.New("malformed authorization header")
	}

	if keys[0] != "ApiKey" {
		return "", errors.New("no api key in header")
	}

	return keys[1], nil
}
