package main

import (
	"github.com/bklimt/hue"
	"github.com/bklimt/volume"
	"net/http"
)

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
