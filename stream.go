package weather

import (
	"errors"
	"log"
	"net/url"
	"time"
)

type StreamToken string

func NewStreamToken(sessionToken SessionToken) (StreamToken, error) {
	session, err := sessionToken.GetSession()
	if err != nil {
		return StreamToken(""), err
	}

	if session == nil {
		return StreamToken(""), errors.New("Not authorized.")
	}

	stmt, err := db.Prepare("insert into stream_token (token, username, created) values (?, ?, ?)")
	if err != nil {
		return StreamToken(""), err
	}

	streamToken := StreamToken(uuid())
	created := time.Now()
	stmt.Bind(streamToken, session.User.Name, created)
	_, err = stmt.Run()
	if err != nil {
		return StreamToken(""), err
	}

	return streamToken, nil
}

func (token StreamToken) isRecent() bool {
	stmt, err := db.Prepare("select * from stream_token where token = ? and deleted > ?")
	if err != nil {
		log.Printf("Error creating recent stream token query: %v\n", err)
		return false
	}

	thirtySeconsAgo := time.Now().Add(-30 * time.Second)
	stmt.Bind(token, thirtySeconsAgo)
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

func (token StreamToken) redeem() bool {
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
		return token.isRecent()
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

func (req *StreamRequest) Check() bool {
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

	token := StreamToken(tokens[0])
	if !token.redeem() {
		return false
	}

	return true
}
