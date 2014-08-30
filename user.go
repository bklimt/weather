package weather

import (
	"code.google.com/p/go.crypto/bcrypt"
	"errors"
	"log"
)

// Information about a user who can log into this website.
type User struct {
	Name string
}

func getUser(username, password string) (*User, error) {
	q, err := db.Prepare("select salt, bcrypt from user where name=?")
	if err != nil {
		return nil, err
	}

	q.Bind(username)
	rows, res, err := q.Exec()
	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		log.Printf("User %v does not exist.\n", username)
		return nil, nil
	}

	if len(rows) > 1 {
		log.Fatal("Too many rows in database for user ", username)
		return nil, errors.New("Database is corrupted.")
	}

	row := rows[0]
	saltCol := res.Map("salt")
	bcryptCol := res.Map("bcrypt")
	salt, bcpw := row.Bin(saltCol), row.Bin(bcryptCol)

	b := append(salt, password...)
	if err := bcrypt.CompareHashAndPassword(bcpw, b); err != nil {
		log.Printf("Password does not match. Error: %v\n", err)
		return nil, nil
	}

	return &User{username}, nil
}
