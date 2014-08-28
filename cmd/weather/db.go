package main

import (
	"code.google.com/p/go.crypto/bcrypt"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/native"
	"io/ioutil"
	"log"
	"time"
)

var db mysql.Conn

func loadConfig(path string) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal("Unable to open config: ", err)
	}

	var info struct {
		Database string
		Username string
		Password string
	}

	if err = json.Unmarshal(b, &info); err != nil {
		log.Fatal("Unable to parse config: ", err)
	}

	db = mysql.New("tcp", "", "127.0.0.1:3306", info.Username, info.Password, info.Database)
	if err := db.Connect(); err != nil {
		log.Fatal("Unable to open database: ", err)
	}

	log.Printf("Opened database %s", info.Database)
}

func uuid() string {
	b := make([]byte, 16)
	rand.Read(b)
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

func createSession(username, password string) (string, time.Time, error) {
	q, err := db.Prepare("select salt, bcrypt from user where name=?")
	if err != nil {
		return "", time.Time{}, err
	}

	q.Bind(username)
	rows, res, err := q.Exec()
	if err != nil {
		return "", time.Time{}, err
	}

	if len(rows) == 0 {
		log.Printf("User %v does not exist.\n", username)
		return "", time.Time{}, nil
	}

	if len(rows) > 1 {
		log.Fatal("Too many rows in database for user ", username)
		return "", time.Time{}, errors.New("Database is corrupted.")
	}

	row := rows[0]
	saltCol := res.Map("salt")
	bcryptCol := res.Map("bcrypt")
	salt, bcpw := row.Bin(saltCol), row.Bin(bcryptCol)

	b := append(salt, password...)
	if err := bcrypt.CompareHashAndPassword(bcpw, b); err != nil {
		log.Printf("Password does not match. Error: %v\n", err)
		return "", time.Time{}, nil
	}

	// Okay, it's actually a valid user.

	stmt, err := db.Prepare("insert into session (token, username, expires) values (?, ?, ?)")
	if err != nil {
		return "", time.Time{}, err
	}

	token := uuid()
	expires := time.Now().AddDate(0, 1, 0)

	stmt.Bind(token, username, expires)
	_, err = stmt.Run()
	if err != nil {
		return "", time.Time{}, err
	}

	return token, expires, nil
}

func deleteSession(session string) error {
	stmt, err := db.Prepare("update session set deleted=? where token=?")
	if err != nil {
		return err
	}

	stmt.Bind(time.Now(), session)
	_, err = stmt.Run()
	if err != nil {
		return err
	}

	return nil
}

func validSession(session string) (bool, error) {
	stmt, err := db.Prepare("select * from session where token = ? and expires > ? and deleted is null")
	if err != nil {
		return false, err
	}

	stmt.Bind(session, time.Now())
	rows, _, err := stmt.Exec()
	if err != nil {
		return false, err
	}

	return len(rows) > 0, nil
}
