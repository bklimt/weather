package main

import (
	"encoding/json"
	"github.com/bklimt/volume"
	"io/ioutil"
	"net/http"
)

type volumeGetHandler struct {
	Card string
}

func (h *volumeGetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if vol, err := volume.GetVolume(h.Card); err != nil {
		writeJsonError(w, err)
	} else {
		result := struct {
			Volume int `json:"volume"`
		}{vol}
		writeJsonResult(w, &result)
	}
}

type volumePutHandler struct {
	Card string
}

func (h *volumePutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
		volume.SetVolume(h.Card, req.Volume)
	}

	res := struct{}{}
	writeJsonResult(w, res)
}
