package weather

import (
	"crypto/rand"
	"fmt"
	"time"
)

// A unique string representing a user logged into the website.
// This token is passed as a cookie on all http requests for authentication purposes.
type SessionToken string

// All information about a user logged into the website.
type Session struct {
	Token   SessionToken
	Expires time.Time
	User    *User
}

func uuid() string {
	b := make([]byte, 16)
	rand.Read(b)
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

// Creates a new session in the website's database for the given user.
// If everything works as intended, but there is no user with the given username and passowrd,
// then nil is returned for both the session and the error.
func NewSession(username, password string) (*Session, error) {
	user, err := getUser(username, password)
	if user == nil {
		return nil, err
	}

	stmt, err := db.Prepare("insert into session (token, username, expires) values (?, ?, ?)")
	if err != nil {
		return nil, err
	}

	token := SessionToken(uuid())
	expires := time.Now().AddDate(0, 1, 0)

	stmt.Bind(token, username, expires)
	_, err = stmt.Run()
	if err != nil {
		return nil, err
	}

	return &Session{token, expires, user}, nil
}

// Deletes the session from the website's database so that it can no longer be used for requests.
// It's good to do this on logout, in addition to sending back an invalid cookie.
func (session SessionToken) Delete() error {
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

// Fetches all the information about this session.
func (token SessionToken) GetSession() (*Session, error) {
	stmt, err := db.Prepare("select username, expires from session where token = ? and expires > ? and deleted is null")
	if err != nil {
		return nil, err
	}

	stmt.Bind(token, time.Now())
	rows, res, err := stmt.Exec()
	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return nil, nil
	}

	name := rows[0].Str(res.Map("username"))
	user := &User{name}

	expires := rows[0].Localtime(res.Map("expires"))
	session := &Session{token, expires, user}

	return session, nil
}
