package auth

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func CreateJWTToken(userID int, jwtSecret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Second * time.Duration(3600))), //expires in 1 hour
		Subject:   strconv.Itoa(userID),
	})
	signedToken, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		log.Printf("Problem signing token. Error: %v\n", err)
		return "", err
	}

	return signedToken, nil
}

func ParseJWTToken(tokenString, tokenSecret string) (int, error) {
	claims := &jwt.RegisteredClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return 0, err
	}

	if !token.Valid {
		err = errors.New("Invalid token")
		return 0, err
	}

	idString, err := token.Claims.GetSubject()
	if err != nil {
		log.Println("Couldn't find userID in token")
		return 0, err
	}

	id, err := strconv.Atoi(idString)
	if err != nil {
		log.Println("Couldn't convert id from string to int")
		return 0, err
	}

	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		return 0, err
	}

	if issuer != "chirpy" {
		return 0, errors.New("invalid issuer")
	}

	return id, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("no authorization header included")
	}

	splitAuthHeader := strings.Split(authHeader, " ")
	if len(splitAuthHeader) < 2 || splitAuthHeader[0] != "Bearer" {
		return "", errors.New("malformed authorization header")
	}

	return splitAuthHeader[1], nil
}
