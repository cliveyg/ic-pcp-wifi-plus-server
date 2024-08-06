package main

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type WPACliResponse struct {
	BSSID     string `json:"bssid"`
	Freq      int    `json:"freq"`
	SSID      string `json:"ssid"`
	IPAddress string `json:"ip_address"`
	KeyMgmt   string `json:"key_mgmt"`
	Address   string `json:"address"`
	UUID      string `json:"uuid"`
}

func (p *WPACliResponse) OrganiseData(lines []string) {

	for _, line := range lines {
		kv := strings.Split(line, "=")
		statusKey := kv[0]
		switch statusKey {
		case "bssid":
			p.BSSID = kv[1]
		case "freq":
			frq, err := strconv.Atoi(kv[1])
			if err != nil {
				log.Fatal(err)
			}
			p.Freq = frq
		case "ip_address":
			p.IPAddress = kv[1]
		case "ssid":
			p.SSID = kv[1]
		case "key_mgmt":
			p.KeyMgmt = kv[1]
		case "address":
			p.Address = kv[1]
		case "uuid":
			p.UUID = kv[1]
		default:
			// do nowt
		}
	}

}

type WifiPlusResponse struct {
	Cmd        string `json:"cmd"`
	Action     string `json:"action"`
	StatusCode int
	Message    string `json:"message"`
	Data       string `default0:"" json:"data"`
}

func (p *WifiPlusResponse) ReturnResponse(w http.ResponseWriter, err error) {

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Debug("Something went bang")
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
			"\"action\": \"" + p.Action + "\"," +
			"\"message\": \"" + p.Message + "\"," +
			"\"data\": " + p.Data + " }"
	} else {
		jsonStr = "{ \"command\": \"" + p.Cmd + "\", " +
			"\"action\": \"" + p.Action + "\"," +
			"\"message\": \"" + p.Message + "\"," +
			"\"data\": {" + p.Data + "} }"
	}

	w.WriteHeader(p.StatusCode)
	if _, err := io.WriteString(w, jsonStr); err != nil {
		log.Fatal(err)
	}

}
