package main

import (
	"encoding/json"
	"github.com/bklimt/weather"
	"log"
	"net/http"
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

	if weather.CheckStream(req) {
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
