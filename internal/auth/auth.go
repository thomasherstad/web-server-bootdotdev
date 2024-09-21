package auth

import "golang.org/x/crypto/bcrypt"

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
