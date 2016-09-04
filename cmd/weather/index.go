package main

import (
	"github.com/bklimt/hue"
	"github.com/bklimt/volume"
	"net/http"
)

type indexHandler struct {
	Card   string
	Hue    *hue.Hue
	Volume bool
}

func (h *indexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if _, ok := checkSession(w, r); !ok {
		return
	}

	values := struct {
		ShowVolume bool
		Volume     int
		Lights     map[string]hue.GetLightResponse
	}{false, 0, map[string]hue.GetLightResponse{}}

	if h.Volume {
		if vol, err := volume.GetVolume(h.Card); err != nil {
			writeJsonError(w, err)
			return
		} else {
			values.ShowVolume = true
			values.Volume = vol
		}
	}

	if h.Hue != nil {
		lightNames := &hue.GetLightsResponse{}
		if err := h.Hue.GetLights(lightNames); err != nil {
			writeJsonError(w, err)
			return
		}

		lights := make(map[string]hue.GetLightResponse)
		for id, _ := range *lightNames {
			l := &hue.GetLightResponse{}
			if err := h.Hue.GetLight(id, l); err != nil {
				writeJsonError(w, err)
				return
			}
			lights[id] = *l
		}

		values.Lights = lights
	}

	if err := templates.ExecuteTemplate(w, "index.html", values); err != nil {
		writeJsonError(w, err)
		return
	}
}
