package weather

import (
	"crypto/rand"
	"fmt"
	"time"
)

type SessionToken string

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
