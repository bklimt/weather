package main

import (
	"encoding/json"
	"fmt"
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
