package database

type User struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
}

func (db *DB) CreateUser(email string) (User, error) {
	var usr User
	usr.Email = email

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
