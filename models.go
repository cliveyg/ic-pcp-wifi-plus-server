package main

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type WifiStatus struct {
	WPASupplicantStatus string `json:"wpa_supplicant_status"`
	UDHCPStatus         string `json:"udhcp_status"`
}

type SSID struct {
	SSID string `json:"ssid"`
}

type Eek struct {
	Error string `json:"error"`
}

type WifiNetwork struct {
	BSSID string `json:"bssid"`
	SSID  string `json:"ssid"`
	Flags string `json:"flags"`
}

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
	Function   string      `json:"function"`
	Action     string      `json:"action"`
	Cmd        string      `json:"cmd,omitempty"`
	StatusCode int         `json:"-"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"`
}

type PiCoreSystemData struct {
	PiCoreVersion       string `json:"picore_version"`
	PiCorePlayerVersion string `json:"picoreplayer_version"`
	SqueezeliteVersion  string `json:"squeezelite_version"`
	LinuxVersion        string `json:"linux_release"`
}

func (p *WifiPlusResponse) ReturnResponse(w http.ResponseWriter, err error) {

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err != nil {
		log.Debug("(((((( 0 ))))))")
		log.WithFields(log.Fields{"err": err}).Error("Something went bang")
		p.StatusCode = 500
		p.Message = "Server error"
		p.Data = Eek{Error: err.Error()}

		//jsonData, err := json.Marshal(p)
		//if err != nil {
		//	log.Debug("=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-[[]]")
		//	log.Fatal(err)
		//}
		//w.WriteHeader(p.StatusCode)
		//if _, err := io.WriteString(w, string(jsonData)); err != nil {
		//	log.Fatal(err)
		//}
		//return
	}
	log.Debug("(((((( 5 ))))))")
	//var jba []byte
	jba, er2 := json.Marshal(p)
	if er2 != nil {
		log.Debug("(((((( 6 ))))))")
		log.Fatal(err)
	}
	w.WriteHeader(p.StatusCode)
	if _, er3 := io.WriteString(w, string(jba)); er3 != nil {
		log.Fatal(er3)
	}

}
