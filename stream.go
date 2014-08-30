package weather

import (
	"errors"
	"log"
	"net/url"
	"time"
)

func createStreamToken(session string) (string, error) {
	u, err := GetSession(session)
	if err != nil {
		return "", err
	}

	if u == nil {
		return "", errors.New("Not authorized.")
	}

	stmt, err := db.Prepare("insert into stream_token (token, username, created) values (?, ?, ?)")
	if err != nil {
		return "", err
	}

	token := uuid()
	created := time.Now()
	stmt.Bind(token, u.Name, created)
	_, err = stmt.Run()
	if err != nil {
		return "", err
	}

	return token, nil
}

func isRecentStreamToken(token string) bool {
	stmt, err := db.Prepare("select * from stream_token where token = ? and deleted > ?")
	if err != nil {
		log.Printf("Error creating recent stream token query: %v\n", err)
		return false
	}

	twoMinutesAgo := time.Now().Add(-2 * time.Minute)
	stmt.Bind(token, twoMinutesAgo)
	rows, _, err := stmt.Exec()
	if err != nil {
		log.Printf("Error looking up recent stream token: %v\n", err)
		return false
	}

	if len(rows) == 0 {
		log.Println("Recent stream token not found.")
		return false
	}

	log.Println("Recent stream token found.")
	return true
}

func redeemStreamToken(token string) bool {
	stmt, err := db.Prepare("select * from stream_token where token = ? and deleted is null")
	if err != nil {
		log.Printf("Error creating stream token query: %v\n", err)
		return false
	}

	stmt.Bind(token)
	rows, _, err := stmt.Exec()
	if err != nil {
		log.Printf("Error looking up stream token: %v\n", err)
		return false
	}

	if len(rows) == 0 {
		log.Println("Stream token not found.")
		return isRecentStreamToken(token)
	}

	stmt, err = db.Prepare("update stream_token set deleted = ? where token = ? and deleted is null")
	if err != nil {
		log.Printf("Error creating stream token update: %v\n", err)
		return false
	}

	stmt.Bind(time.Now(), token)
	_, err = stmt.Run()
	if err != nil {
		log.Printf("Unable to delete stream token: %v\n", err)
		return false
	}

	return true
}

type StreamRequest struct {
	Action string
	Server string
	Port   string
	Client string
	Mount  string
	User   string
	Pass   string
	Ip     string
	Agent  string
}

func CheckStream(req StreamRequest) bool {
	if req.Action != "listener_add" {
		return false
	}

	u, err := url.Parse(req.Mount)
	if err != nil {
		log.Printf("Error parsing mount %v: %v\n", req.Mount, err)
		return false
	}

	values := u.Query()
	tokens := values["token"]
	if len(tokens) != 1 {
		log.Printf("Incorrect token count: %v\n", len(tokens))
		return false
	}

	if !redeemStreamToken(tokens[0]) {
		return false
	}

	return true
}
