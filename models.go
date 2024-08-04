package main

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
)

type WifiPlusResponse struct {
	Cmd        string `json:"cmd"`
	StatusCode int
	Message    string `json:"message"`
	Data       string `default0:"" json:"data"`
}

func (p *WifiPlusResponse) FormatResponse(w http.ResponseWriter, err error) {

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err != nil {
		log.Error("Error is %s", err)
		p.StatusCode = 500
		jsonData, _ := json.Marshal(p)
		/*
			jsonStr := "{ \"command\": \"" + p.Cmd + "\", " +
				"\"message\": \"error\"," +
				"\"data\": {" + err.Error() + "} }"
		*/

		w.WriteHeader(p.StatusCode)
		if _, err := io.WriteString(w, string(jsonData)); err != nil {
			log.Fatal(err)
		}
		return
	}

	jsonStr := "{ \"command\": \"" + p.Cmd + "\", " +
		"\"message\": \"" + p.Message + "\"," +
		"\"data\": {" + p.Data + "} }"

	w.WriteHeader(p.StatusCode)
	if _, err := io.WriteString(w, jsonStr); err != nil {
		log.Fatal(err)
	}

}
