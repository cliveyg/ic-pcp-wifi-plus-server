package main

import (
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"regexp"
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

type WifiDetails struct {
	BSSID    string `json:"bssid"`
	SSID     string `json:"ssid"`
	Password string `json:"password"`
}

type ShellSwitchInfo struct {
	Res interface{} `json:"data,omitempty"`
}

type SysStatus struct {
	Wifi string `json:"wifi,omitempty"`
	WAP  string `json:"wap,omitempty"`
	Ping int    `json:"ping"`
}

type ShellResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type SwitcherInfo struct {
	APMode     string `json:"ap_mode"`
	APAddress  string `json:"ap_address"`
	APStatus   int    `json:"ap_status"`
	Wifi       string `json:"wifi"`
	WifiStatus int    `json:"wifi_status"`
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
		log.WithFields(log.Fields{"err": err}).Error("Something went bang!")
		p.StatusCode = 500
		p.Message = "Server error"
		p.Data = Eek{Error: err.Error()}

	}

	var jba []byte
	jba, err = json.Marshal(p)
	if err != nil {
		log.Fatal(err)
	}
	w.WriteHeader(p.StatusCode)
	if _, err = io.WriteString(w, string(jba)); err != nil {
		log.Fatal(err)
	}
}

type WAPConfig struct {
	SSID        string `json:"ssid"`
	APIPAddress string `json:"ap_ip_address"`
	Password    string `json:"password"`
	CountryCode string `json:"country_code"`
	Channel     int    `json:"channel"`
}

func (c *WAPConfig) ValidateInput(err *error) {
	if len(c.SSID) > 31 {
		*err = errors.New("SSID is too long")
		return
	}
	if len(c.Password) < 8 || len(c.Password) > 63 {
		*err = errors.New("Password must be between 8 and 63 characters")
		return
	}
	//TODO - could do more in depth country code validation
	match, rErr := regexp.MatchString("[A-Z][A-Z]", c.CountryCode)
	if !match || rErr != nil {
		*err = errors.New("Country code not valid")
		return
	}
	// must match 10.nnn.nnn.1 where nnn is between 1 and 255 inclusive
	match, rErr = regexp.MatchString("[10].[1-255].[1-255].1", c.APIPAddress)
	if !match || rErr != nil {
		*err = errors.New("Access point IP address not valid")
		return
	}
	if c.Channel < 1 || c.Channel > 140 {
		*err = errors.New("Channel is not valid")
		return
	}
}

func (c *WAPConfig) Stringify() string {
	rs := "SSID=" + c.SSID + "&"
	rs = rs + "IP=" + c.APIPAddress + "&"
	rs = rs + "Ch=" + fmt.Sprint(c.Channel) + "&"
	rs = rs + "CC=" + c.CountryCode + "&"
	rs = rs + "Pass=" + c.Password
	return rs
}
