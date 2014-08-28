package main

import (
	"encoding/json"
	"github.com/bklimt/hue"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

var philipsHue *hue.Hue

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
