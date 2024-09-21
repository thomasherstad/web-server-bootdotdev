package auth

import (
	"log"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func CreateJWTToken(expireInSeconds, userID int, jwtSecret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Second * time.Duration(expireInSeconds))),
		Subject:   strconv.Itoa(userID),
	})
	signedToken, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		log.Println(err)
		return "", err
	}

	return signedToken, nil
}
