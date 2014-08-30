package main

import (
	"flag"
	"fmt"
	"github.com/bklimt/hue"
	"github.com/bklimt/weather"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	// Server flags
	port := flag.Int("port", 9696, "The port to listen on.")
	config := flag.String("config", "./config.json", "Database config file.")
	templateFlags()
	hue.Flags()
	flag.Parse()

	weather.LoadConfig(*config)

	r := mux.NewRouter()

	if err := loadTemplates(r); err != nil {
		log.Fatal("Unable to load static file.")
	}

	philipsHue = hue.FromFlags()

	r.HandleFunc("/", handleIndex).Methods("GET")
	r.HandleFunc("/session", handleSessionPost).Methods("POST")
	r.HandleFunc("/session", handleSessionDelete).Methods("DELETE")
	r.HandleFunc("/volume", handleVolumeGet).Methods("GET")
	r.HandleFunc("/volume", handleVolumePut).Methods("PUT")
	r.HandleFunc("/light", handleLightsGet).Methods("GET")
	r.HandleFunc("/light/{id}", handleLightPut).Methods("PUT")
	r.HandleFunc("/check_stream", handleCheckStreamPost).Methods("POST")

	http.Handle("/", r)

	log.Printf("Serving on port %d...", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}
