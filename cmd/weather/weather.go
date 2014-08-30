// The weather command is an http server that provides a mobile website that can be used to control
// various home automation systems. It is designed to run on Raspberry Pi and provides controls for
// setting the volume of sound output, controlling Philips Hue light bulbs, and streaming audio
// from a microphone attached to the Pi.
package main

import (
	"flag"
	"fmt"
	"github.com/bklimt/hue"
	"github.com/bklimt/weather"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"net/url"
)

func main() {
	// Server flags
	port := flag.Int("port", 9696, "The port to listen on.")
	config := flag.String("config", "./config.json", "Database config file.")
	streamUrl := flag.String("audio_stream_url", "http://example.com:8000/stream.ogg", "Audio stream URL.")

	templateFlags()
	hue.Flags()
	flag.Parse()

	weather.LoadConfig(*config)

	r := mux.NewRouter()

	if err := loadTemplates(r); err != nil {
		log.Fatal("Unable to load static file.")
	}

	philipsHue = hue.FromFlags()

	u, err := url.Parse(*streamUrl)
	if err != nil {
		log.Fatal("Unable to parse stream url: %v", streamUrl)
	}
	streamHandler := streamGetHandler{u}

	r.HandleFunc("/", handleIndex).Methods("GET")
	r.HandleFunc("/session", handleSessionPost).Methods("POST")
	r.HandleFunc("/session", handleSessionDelete).Methods("DELETE")
	r.HandleFunc("/volume", handleVolumeGet).Methods("GET")
	r.HandleFunc("/volume", handleVolumePut).Methods("PUT")
	r.HandleFunc("/light", handleLightsGet).Methods("GET")
	r.HandleFunc("/light/{id}", handleLightPut).Methods("PUT")
	r.HandleFunc("/check_stream", handleCheckStreamPost).Methods("POST")
	r.Handle("/stream", &streamHandler).Methods("GET")

	http.Handle("/", r)

	log.Printf("Serving on port %d...", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}
