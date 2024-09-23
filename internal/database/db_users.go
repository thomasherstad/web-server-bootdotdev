package database

import (
	"errors"
)

type User struct {
	ID          int    `json:"id"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	IsChirpyRed bool   `json:"is_chirpy_red"`
}

var ErrNotExists = errors.New("user doesn't exist")

func (db *DB) CreateUser(email, password string) (User, error) {

	database, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	usr := User{
		ID:          len(database.Users) + 1,
		Email:       email,
		Password:    password,
		IsChirpyRed: false,
	}
	database.Users[usr.ID] = usr

	//Save to the disk
	err = db.writeDB(database)
	if err != nil {
		return User{}, err
	}

	return usr, nil

}

func (db *DB) GetUserByEmail(email string) (User, error) {
	database, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	for _, usr := range database.Users {
		if usr.Email == email {
			return usr, nil
		}
	}

	notFound := errors.New("user not found")
	return User{}, notFound
}

func (db *DB) GetUserById(id int) (User, error) {
	database, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	user, ok := database.Users[id]
	if !ok {
		return User{}, ErrNotExists
	}

	return user, nil
}

func (db *DB) UpdateUser(id int, newEmail, newPassword string) (User, error) {
	database, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	usr, err := db.GetUserById(id)
	if err != nil {
		return User{}, err
	}

	usr.Email = newEmail
	usr.Password = newPassword
	database.Users[id] = usr

	db.writeDB(database)

	return usr, nil
}

func (db *DB) UpgradeUser(id int) error {
	database, err := db.loadDB()
	if err != nil {
		return err
	}

	user, ok := database.Users[id]
	if !ok {
		return ErrNotExists
	}

	user.IsChirpyRed = true
	database.Users[user.ID] = user

	db.writeDB(database)
	return nil
}
