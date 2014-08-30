package main

import (
	"encoding/json"
	"github.com/bklimt/weather"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const sessionCookieName = "session"

func handleSessionPost(w http.ResponseWriter, r *http.Request) {
	req := struct {
		Username string
		Password string
	}{"", ""}

	if r.Header.Get("Content-Type") == "application/json" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			writeJsonError(w, err)
			return
		}

		err = json.Unmarshal(body, &req)
		if err != nil {
			writeJsonError(w, err)
			return
		}
	} else {
		req.Username = r.PostFormValue("username")
		req.Password = r.PostFormValue("password")
	}

	session, err := weather.NewSession(req.Username, req.Password)
	if err != nil {
		writeJsonError(w, err)
		return
	}

	if session == nil {
		http.Redirect(w, r, "/login?failed=true", http.StatusSeeOther)
		return
	}

	http.SetCookie(w, &http.Cookie{
		sessionCookieName,     // Name
		string(session.Token), // Value
		"",              // Path
		"",              // Domain
		session.Expires, // Expires
		"",              // RawExpires
		0,               // MaxAge
		false,           // true, // Secure
		false,           // HttpOnly
		"",              // Raw
		nil,             // Unparsed
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func handleSessionDelete(w http.ResponseWriter, r *http.Request) {
	if sessionCookie, err := r.Cookie(sessionCookieName); err == nil {
		token := weather.SessionToken(sessionCookie.Value)
		if err = token.Delete(); err != nil {
			log.Printf("Unable to delete session: %v\n", err)
		}

		expires := time.Now().AddDate(-1, 0, 0)

		http.SetCookie(w, &http.Cookie{
			sessionCookieName, // Name
			"",                // Value
			"",                // Path
			"",                // Domain
			expires,           // Expires
			"",                // RawExpires
			0,                 // MaxAge
			false,             // true, // Secure
			false,             // HttpOnly
			"",                // Raw
			nil,               // Unparsed
		})
	}

	http.Redirect(w, r, "/login.html", http.StatusSeeOther)
}

func checkSession(w http.ResponseWriter, r *http.Request) (weather.SessionToken, bool) {
	if sessionCookie, err := r.Cookie(sessionCookieName); err != nil {
		if err != http.ErrNoCookie {
			log.Printf("Error reading cookie: %v\n", err)
		}
	} else {
		token := weather.SessionToken(sessionCookie.Value)
		if user, err := token.GetSession(); err != nil {
			log.Printf("Error validating cookie: %v\n", err)
		} else {
			if user != nil {
				return token, true
			}
		}
	}

	http.Redirect(w, r, "/login.html", http.StatusSeeOther)
	return weather.SessionToken(""), false
}
