package weather

import (
	"code.google.com/p/go.crypto/bcrypt"
	"crypto/rand"
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

// Saves all of the information for this user to the database.
// This function does not check that the password is correct or any permissions.
// If a user with that name already exists, this function fails.
func (user User) Save(password string) error {
	salt := make([]byte, 20)
	_, err := rand.Read(salt)
	if err != nil {
		log.Printf("Unable to generate salt: %v\n", err)
		return err
	}
	// log.Printf("salt: %v\n", salt)

	b := append(salt, password...)
	bcpw, err := bcrypt.GenerateFromPassword(b, 0)
	if err != nil {
		log.Printf("Unable to bcrypt password: %v\n", err)
		return err
	}

	stmt, err := db.Prepare("INSERT INTO user (name, salt, bcrypt) values (?, ?, ?)")
	if err != nil {
		log.Printf("Unable to prepare statement: %v\n", err)
		return err
	}

	stmt.Bind(user.Name, salt, bcpw)
	_, err = stmt.Run()
	if err != nil {
		log.Printf("Unable to insert row into database: %v\n", err)
		return err
	}

	return nil
}
