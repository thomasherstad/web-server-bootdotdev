package database

import (
	"errors"
	"time"
)

type RefreshToken struct {
	UserId       int       `json:"id"`
	RefreshToken string    `json:"refresh_token"`
	Expiry       time.Time `json:"-"`
}

func (db *DB) AddUserRefreshToken(id int, refreshToken string, expirationTime time.Time) (User, error) {
	database, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	usr, err := db.GetUserById(id)
	if err != nil {
		return User{}, err
	}

	database.Users[id] = usr

	database.RefreshTokens[refreshToken] = RefreshToken{
		UserId:       usr.ID,
		RefreshToken: refreshToken,
		Expiry:       expirationTime,
	}

	db.writeDB(database)

	return usr, nil
}

func (db *DB) GetUserByRefreshToken(refreshToken string) (User, error) {
	database, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	refTok, ok := database.RefreshTokens[refreshToken]
	if !ok {
		return User{}, errors.New("refresh token not found")
	}

	user := database.Users[refTok.UserId]

	return user, nil
}

func (db *DB) DeleteRefreshToken(refreshToken string) error {
	database, err := db.loadDB()
	if err != nil {
		return err
	}

	if _, ok := database.RefreshTokens[refreshToken]; !ok {
		return errors.New("refresh token not found")
	}

	delete(database.RefreshTokens, refreshToken)
	return nil
}

func (db *DB) GetRefreshToken(refreshToken string) (RefreshToken, error) {
	database, err := db.loadDB()
	if err != nil {
		return RefreshToken{}, err
	}

	refToken, ok := database.RefreshTokens[refreshToken]
	if !ok {
		return RefreshToken{}, errors.New("Refresh token not in database")
	}

	return refToken, nil
}
