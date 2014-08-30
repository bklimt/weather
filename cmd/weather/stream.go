package main

import (
	"encoding/json"
	"fmt"
	"github.com/bklimt/weather"
	"log"
	"net/http"
	"net/url"
)

func handleCheckStreamPost(w http.ResponseWriter, r *http.Request) {
	req := weather.StreamRequest{
		r.FormValue("action"),
		r.FormValue("server"),
		r.FormValue("port"),
		r.FormValue("client"),
		r.FormValue("mount"),
		r.FormValue("user"),
		r.FormValue("pass"),
		r.FormValue("ip"),
		r.FormValue("agent"),
	}

	if j, err := json.Marshal(req); err != nil {
		log.Printf("Unable to create json for req %v: %v\n", req, err)
	} else {
		log.Printf("Received request to auth stream: %v\n", string(j))
	}

	if req.Check() {
		log.Println("Allowing stream request.")
		w.Header().Set("icecast-auth-user", "1")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Ok"))
	} else {
		log.Println("Denying stream request.")
		w.Header().Set("icecast-auth-message", "not authorized")
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("No"))
	}
}

type streamGetHandler struct {
	streamUrl *url.URL
}

func (h *streamGetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	session, ok := checkSession(w, r)
	if !ok {
		return
	}

	token, err := weather.NewStreamToken(session)
	if err != nil {
		writeJsonError(w, err)
		return
	}

	url := fmt.Sprintf("%v?token=%v", h.streamUrl, token)

	values := struct {
		StreamUrl string
	}{url}

	if err := templates.ExecuteTemplate(w, "stream.html", values); err != nil {
		writeJsonError(w, err)
		return
	}
}
