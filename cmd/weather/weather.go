package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/bklimt/hue"
	"github.com/bklimt/volume"
	"github.com/gorilla/mux"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var philipsHue *hue.Hue
var templates *template.Template

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

type light struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	On   bool   `json:"on"`
	Hue  int    `json:"hue"`
	Sat  int    `json:"sat"`
	Bri  int    `json:"bri"`
}

func handleLightsGet(w http.ResponseWriter, r *http.Request) {
	lights := &hue.GetLightsResponse{}
	if err := philipsHue.GetLights(lights); err != nil {
		writeJsonError(w, err)
		return
	}

	var result []light
	for id, _ := range *lights {
		l := &hue.GetLightResponse{}
		if err := philipsHue.GetLight(id, l); err != nil {
			writeJsonError(w, err)
			return
		}
		s := l.State
		result = append(result, light{id, l.Name, s.On, s.Hue, s.Sat, s.Bri})
	}

	writeJsonResult(w, &result)
}

func handleLightPut(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeJsonError(w, err)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	req := struct {
		Hue *int
		Sat *int
		Bri *int
		On  *bool
	}{nil, nil, nil, nil}
	err = json.Unmarshal(body, &req)
	if err != nil {
		writeJsonError(w, err)
		return
	}

	l := &hue.PutLightRequest{}
	if req.Hue != nil {
		l.Hue = req.Hue
	}
	if req.Sat != nil {
		l.Sat = req.Sat
	}
	if req.Bri != nil {
		l.Bri = req.Bri
	}
	if req.On != nil {
		l.On = req.On
	}
	if err := philipsHue.PutLight(id, l); err != nil {
		writeJsonError(w, err)
		return
	}

	res := struct{}{}
	writeJsonResult(w, res)
}

func handleVolumeGet(w http.ResponseWriter, r *http.Request) {
	if vol, err := volume.GetVolume(); err != nil {
		writeJsonError(w, err)
	} else {
		result := struct {
			Volume int `json:"volume"`
		}{vol}
		writeJsonResult(w, &result)
	}
}

func handleVolumePut(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeJsonError(w, err)
		return
	}

	req := struct {
		Volume int
	}{-1}
	err = json.Unmarshal(body, &req)
	if err != nil {
		writeJsonError(w, err)
		return
	}

	if req.Volume >= 0 {
		volume.SetVolume(req.Volume)
	}

	res := struct{}{}
	writeJsonResult(w, res)
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
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
	templateDir := flag.String("templates", "./templates", "The directory to search for template files.")
	staticDir := flag.String("static", "./static", "The directory to search for static files.")

	// Hue flags
	ip := flag.String("ip", "192.168.1.3", "IP Address of Philips Hue hub.")
	userName := flag.String("username", "HueGoRaspberryPiUser", "Username for Hue hub.")
	deviceType := flag.String("device_type", "HueGoRaspberryPi", "Device type for Hue hub.")

	flag.Parse()

	r := mux.NewRouter()

	log.Println("Loading templates:")
	templates = template.New("templates")
	err := filepath.Walk(*templateDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(*templateDir, path)
		if err != nil {
			return err
		}

		t := template.Must(template.ParseFiles(path))
		log.Printf("  %v\n", rel)
		templates = template.Must(templates.AddParseTree(rel, t.Tree))
		return nil
	})
	if err != nil {
		log.Fatal("Unable to load templates.")
	}

	log.Println("Loading static files:")
	err = filepath.Walk(*staticDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(*staticDir, path)
		if err != nil {
			return err
		}
		log.Printf("  %v\n", rel)

		r.HandleFunc("/"+rel, func(w http.ResponseWriter, r *http.Request) {
			b, err := ioutil.ReadFile(path)
			if err != nil {
				writeJsonError(w, err)
				return
			}
			w.Write(b)
		}).Methods("GET")

		return nil
	})
	if err != nil {
		log.Fatal("Unable to load static file.")
	}

	philipsHue = &hue.Hue{*ip, *userName, *deviceType}

	r.HandleFunc("/", handleIndex).Methods("GET")
	r.HandleFunc("/volume", handleVolumeGet).Methods("GET")
	r.HandleFunc("/volume", handleVolumePut).Methods("PUT")
	r.HandleFunc("/light", handleLightsGet).Methods("GET")
	r.HandleFunc("/light/{id}", handleLightPut).Methods("PUT")

	http.Handle("/", r)

	log.Printf("Serving on port %d...", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}
