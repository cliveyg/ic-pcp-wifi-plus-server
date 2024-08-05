package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// ----------------------------------------------------------------------------

func (a *App) testTings(w http.ResponseWriter, _ *http.Request) {

	log.Debug("In testTings")
	rc, err := exec.Command("sh", "-c", "wpa_cli status").Output()
	if err != nil {
		log.Fatal(err)
	}
	lines := strings.Split(strings.TrimSpace(string(rc)), "\n")
	lines = append(lines[:0], lines[1:]...)

	pr := WifiPlusResponse{
		Cmd:        "testTings",
		StatusCode: 200,
		Action:     "Testing stuff",
		Message:    "testy test"}
	pr.FormatResponse(w, err)

}

func (a *App) wifiAction(w http.ResponseWriter, r *http.Request) {

	log.Debug("In wifiAction")
	vars := mux.Vars(r)
	wifiAction := vars["action"]

	//TODO: Check input string more thoroughly

	var rc []byte
	var sr string
	var err error
	var args []string
	pr := WifiPlusResponse{
		Cmd:    "wifiAction",
		Action: wifiAction,
	}

	switch wifiAction {
	case "restart":
		// send message before enacting command
		pr.StatusCode = 202
		pr.FormatResponse(w, nil)
		time.Sleep(2 * time.Second)
		rc, err = exec.Command("sh", "-c", "/usr/local/etc/init.d/wifi wlan0 stop && /usr/local/etc/init.d/wifi wlan0 start").Output()
		log.Debug(rc)
		return
	case "status":
		args = []string{"wlan0", "status"}
		sr, err = a.ExecCmd("/usr/local/etc/init.d/wifi", args)

		if strings.Contains(sr, "wpa_supplicant running") {
			pr.Message = "wpa_supplicant running"
			pr.StatusCode = 200
		} else {
			pr.Message = "wpa_supplicant not running"
			pr.StatusCode = 404
		}
	case "ssid":
		args = []string{"-r"}
		sr, err = a.ExecCmd("iwgetid", args)

		if sr == "" {
			pr.StatusCode = 404
			pr.Message = "No SSID found"
		} else {
			pr.StatusCode = 200
			pr.Message = "SSID found"
			pr.Data = `"SSID": "` + sr + `"`
		}
	default:
		// do nowt
		pr.StatusCode = 400
	}

	pr.FormatResponse(w, err)
}

func (a *App) getWPACliStatus(w http.ResponseWriter, _ *http.Request) {

	log.Debug("In getWPACliStatus")
	rc, err := exec.Command("sh", "-c", "wpa_cli status").Output()
	if err != nil {
		log.Fatal(err)
	}
	lines := strings.Split(strings.TrimSpace(string(rc)), "\n")
	// remove first line
	lines = append(lines[:0], lines[1:]...)

	wpaData := WPACliResponse{}
	wpaData.OrganiseData(lines)

	jsonData, _ := json.Marshal(wpaData)

	pr := WifiPlusResponse{
		Cmd:        "getWPACliStatus",
		Action:     "wpa_cli",
		StatusCode: 200,
		Message:    "wpa_cli status",
		Data:       string(jsonData)}
	pr.FormatResponse(w, err)

}

func (a *App) getPiCoreDetails(w http.ResponseWriter, _ *http.Request) {

	log.Debug("In getPiCoreDetails")
	retData, err := exec.Command("sh", "-c", "cd cgi-bin && sudo ./wifi-plus.sh wp_picore_details").Output()
	if err != nil {
		log.Fatal(err)
	}
	pr := WifiPlusResponse{
		Cmd:        "getPiCoreDetails",
		Action:     "wifi-plus.sh",
		StatusCode: 200,
		Message:    "piCore details",
		Data:       string(retData)}
	pr.FormatResponse(w, err)

}

func (a *App) RebootSystem(w http.ResponseWriter, _ *http.Request) {
	log.Debug("In RebootSystem")
	pr := WifiPlusResponse{
		Cmd:        "RebootSystem",
		Action:     "pcp rb",
		StatusCode: 202,
		Message:    "System rebooting"}
	pr.FormatResponse(w, nil)
	time.Sleep(2 * time.Second)
	_, err := exec.Command("sh", "-c", "sudo pcp rb").Output()
	if err != nil {
		log.Fatal(err)
	}
}

func (a *App) getSystemStatus(w http.ResponseWriter, _ *http.Request) {

	log.Debug("In getSystemStatus")
	rc, err := exec.Command("sh", "-c", "cd cgi-bin && ./wifi-plus.sh wp_status 200").Output()
	if err != nil {
		log.Fatal(err)
	}
	rcInt, err := strconv.Atoi(strings.TrimSpace(string(rc)))
	if err != nil {
		log.Fatal(err)
	}

	pr := WifiPlusResponse{
		Cmd:        "getSystemStatus",
		Action:     "wifi-plus.sh",
		StatusCode: rcInt,
		Message:    "System running"}
	pr.FormatResponse(w, err)

}
