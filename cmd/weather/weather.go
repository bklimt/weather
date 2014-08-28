package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/bklimt/hue"
	"github.com/bklimt/volume"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func writeJsonError(w http.ResponseWriter, err error) {
	var j struct {
		Err string `json:"error"`
	}
	j.Err = fmt.Sprintf("%v", err)
	s, err2 := json.Marshal(j)
	if err2 != nil {
		// Well, we did our best.
		fmt.Fprintf(w, "Unable to generate json for error:\n%v\n%v", err, err2)
		return
	}
	fmt.Fprintf(w, "%s", s)
}

func writeJsonResult(w http.ResponseWriter, result interface{}) {
	s, err := json.Marshal(result)
	if err != nil {
		writeJsonError(w, err)
		return
	}
	fmt.Fprintf(w, "%s", s)
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	if !checkSession(w, r) {
		return
	}

	vol := 0
	if vol2, err := volume.GetVolume(); err != nil {
		writeJsonError(w, err)
		return
	} else {
		vol = vol2
	}

	lightNames := &hue.GetLightsResponse{}
	if err := philipsHue.GetLights(lightNames); err != nil {
		writeJsonError(w, err)
		return
	}

	lights := make(map[string]hue.GetLightResponse)
	for id, _ := range *lightNames {
		l := &hue.GetLightResponse{}
		if err := philipsHue.GetLight(id, l); err != nil {
			writeJsonError(w, err)
			return
		}
		lights[id] = *l
	}

	values := struct {
		Volume int
		Lights map[string]hue.GetLightResponse
	}{vol, lights}

	if err := templates.ExecuteTemplate(w, "index.html", values); err != nil {
		writeJsonError(w, err)
		return
	}
}

func main() {
	// Server flags
	port := flag.Int("port", 9696, "The port to listen on.")
	config := flag.String("config", "./config.json", "Database config file.")

	templateFlags()

	// Hue flags
	ip := flag.String("ip", "192.168.1.3", "IP Address of Philips Hue hub.")
	userName := flag.String("username", "HueGoRaspberryPiUser", "Username for Hue hub.")
	deviceType := flag.String("device_type", "HueGoRaspberryPi", "Device type for Hue hub.")

	flag.Parse()
	loadConfig(*config)

	r := mux.NewRouter()

	if err := loadTemplates(r); err != nil {
		log.Fatal("Unable to load static file.")
	}

	philipsHue = &hue.Hue{*ip, *userName, *deviceType}

	r.HandleFunc("/", handleIndex).Methods("GET")
	r.HandleFunc("/session", handleSessionPost).Methods("POST")
	r.HandleFunc("/session", handleSessionDelete).Methods("DELETE")
	r.HandleFunc("/volume", handleVolumeGet).Methods("GET")
	r.HandleFunc("/volume", handleVolumePut).Methods("PUT")
	r.HandleFunc("/light", handleLightsGet).Methods("GET")
	r.HandleFunc("/light/{id}", handleLightPut).Methods("PUT")

	http.Handle("/", r)

	log.Printf("Serving on port %d...", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}
