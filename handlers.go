package main

import (
	"encoding/json"
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
	wpaData := WPACliResponse{}
	for _, line := range lines {
		kv := strings.Split(line, "=")
		statusKey := kv[0]
		switch statusKey {
		case "bssid":
			wpaData.BSSID = kv[1]
		case "freq":
			frq, err := strconv.Atoi(kv[1])
			if err != nil {
				log.Fatal(err)
			}
			wpaData.Freq = frq
		case "ip_address":
			wpaData.IPAddress = kv[1]
		case "ssid":
			wpaData.SSID = kv[1]
		case "key_mgmt":
			wpaData.KeyMgmt = kv[1]
		case "mac_address":
			wpaData.MACAddress = kv[1]
		case "uuid":
			wpaData.UUID = kv[1]
		default:
			// do nowt
		}
	}

	jsonData, _ := json.Marshal(wpaData)

	pr := WifiPlusResponse{
		Cmd:        "testTings",
		StatusCode: 200,
		Message:    "wpa_cli test",
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
		StatusCode: 200,
		Message:    "piCore details",
		Data:       string(retData)}
	pr.FormatResponse(w, err)

}

func (a *App) RebootSystem(w http.ResponseWriter, _ *http.Request) {
	log.Debug("In RebootSystem")
	pr := WifiPlusResponse{
		Cmd:        "RebootSystem",
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
		StatusCode: rcInt,
		Message:    "System running"}
	pr.FormatResponse(w, err)

}

func (a *App) getWifiStatus(w http.ResponseWriter, _ *http.Request) {

	log.Debug("In getWifiStatus")
	args := []string{"wlan0", "status"}
	rc, err := a.ExecCmd("/usr/local/etc/init.d/wifi", args)
	pr := WifiPlusResponse{Cmd: "getWifiStatus"}

	if strings.Contains(rc, "wpa_supplicant running") {
		pr.Message = "wpa_supplicant running"
		pr.StatusCode = 200
	} else {
		pr.Message = "wpa_supplicant not running"
		pr.StatusCode = 404
	}
	pr.FormatResponse(w, err)
}

func (a *App) getWifiSSID(w http.ResponseWriter, _ *http.Request) {

	log.Debug("In getWifiSSID")
	args := []string{"-r"}
	SSID, err := a.ExecCmd("iwgetid", args)
	pr := WifiPlusResponse{Cmd: "getWifiSSID"}

	if SSID == "" {
		pr.StatusCode = 404
		pr.Message = "No SSID found"
	} else {
		pr.StatusCode = 200
		pr.Message = "SSID found"
		pr.Data = `"SSID": "` + SSID + `"`
	}
	pr.FormatResponse(w, err)
}
