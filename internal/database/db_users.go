package database

import "errors"

type User struct {
	Id       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (db *DB) CreateUser(email, password string) (User, error) {
	usr := User{
		Email:    email,
		Password: password,
	}

	database, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	usr.Id = len(database.Users) + 1
	database.Users[usr.Id] = usr

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
