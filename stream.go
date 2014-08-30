package weather

import (
	"errors"
	"log"
	"net/url"
	"time"
)

// A one-time-use token for getting access to an audio stream.
// This server trades a session token in for a stream token and then gives the stream token back to
// the client. The stream token can then be given to a separate streaming server. The streaming
// server checks with this server to make sure the stream token is valid. If it is, the token is
// invalidated so that it can't be used again.
type StreamToken string

// Creates a new StreamToken in the server's database for the given session.
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

// Returns true if the token has been used in the last 30 seconds.
// This function isn't currently used but could be to allow a window of use for a token.
// This can be useful for example when navigating directly to a stream in Chrome, because
// Chrome actually makes two requests for the content.
func (token StreamToken) isRecent() bool {
	stmt, err := db.Prepare("select * from stream_token where token = ? and deleted > ?")
	if err != nil {
		log.Printf("Error creating recent stream token query: %v\n", err)
		return false
	}

	thirtySecondsAgo := time.Now().Add(-30 * time.Second)
	stmt.Bind(token, thirtySecondsAgo)
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

// Returns true if this stream token is still good, and deletes it.
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
		// return token.isRecent()
    return false
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

// The information in a request from a separate streaming server, such as IceCast2, used for asking
// this server whether access to a particular stream should be allowed.
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

// Returns whether the given request from a separate streaming server should be allowed.
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
