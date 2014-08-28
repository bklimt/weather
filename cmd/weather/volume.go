package main

import (
	"encoding/json"
	"github.com/bklimt/volume"
	"io/ioutil"
	"net/http"
)

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
