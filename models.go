package main

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
)

type WPACliResponse struct {
	BSSID     string `json:"bssid"`
	Freq      int    `json:"freq"`
	SSID      string `json:"ssid"`
	IPAddress string `json:"ip_address"`
	KeyMgmt   string `json:"key_mgmt"`
	Address   string `json:"mac_address"`
	UUID      string `json:"uuid"`
}

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
		w.WriteHeader(p.StatusCode)
		if _, err := io.WriteString(w, string(jsonData)); err != nil {
			log.Fatal(err)
		}
		return
	}

	var jsonStr string
	// check to avoid double brackets
	if substr(p.Data, 0, 1) == "{" {
		jsonStr = "{ \"command\": \"" + p.Cmd + "\", " +
			"\"message\": \"" + p.Message + "\"," +
			"\"data\": " + p.Data + " }"
	} else {
		jsonStr = "{ \"command\": \"" + p.Cmd + "\", " +
			"\"message\": \"" + p.Message + "\"," +
			"\"data\": {" + p.Data + "} }"
	}

	w.WriteHeader(p.StatusCode)
	if _, err := io.WriteString(w, jsonStr); err != nil {
		log.Fatal(err)
	}

}
