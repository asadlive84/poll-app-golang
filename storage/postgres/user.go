package postgres

import (
	"log"
	"poll-app/storage"
)

const createUserQuery = `
	INSERT INTO users 
	(name, username, email, password) 
	VALUES (:name, :username, :email, :password) 
	RETURNING id, created_at, updated_at
`

func (s *Storage) SaveUserDB(user storage.User) (int32, error) {
	stmt, err := s.db.PrepareNamed(createUserQuery)

	if err != nil {
		log.Println(err)
		return 0, err
	}
	var id int32
	if err := stmt.Get(&id, user); err != nil {
		return 0, err
	}

	return id, nil

}

const getUserQuery = `
	SELECT * from users 
	where email=$1 and password=$2
`

func (s *Storage) GetUser(email string, password string) (*storage.User, error) {
	user := storage.User{}
	if err := s.db.Get(&user, getUserQuery, email, password); err != nil {
		return nil, err
	}
	return &user, nil
}

const getUserID = `SELECT id FROM users WHERE email=$1 LIMIT 1`

func (s *Storage) GetUserDBID(email interface{}) (int32, error) {
	user := storage.User{}

	if _, ok := email.(string); ok {
		if err := s.db.Get(&user, getUserID, email); err != nil {
			return 0, err
		}
	}

	return user.ID, nil
}
