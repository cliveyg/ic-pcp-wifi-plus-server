package main

import (
	log "github.com/sirupsen/logrus"
	"net/http"
	"os/exec"
	"strings"
)

// ----------------------------------------------------------------------------

func (a *App) testTings(w http.ResponseWriter, _ *http.Request) {

	log.Debug("In testTings")
	rc, err := exec.Command("sh", "-c", "cd cgi-bin && sudo ./wifi-plus.sh wp_test").Output()
	a.FormatResponse(w, "testTings", 200, strings.TrimSpace(string(rc)), "", err)

}

func (a *App) getSystemStatus(w http.ResponseWriter, _ *http.Request) {

	log.Debug("In getSystemStatus")
	rc, err := exec.Command("sh", "-c", "cd cgi-bin && sudo ./wifi-plus.sh wp_status 200").Output()

	a.FormatResponse(w, "getSystemStatus", 200, strings.TrimSpace(string(rc)), "", err)

}

func (a *App) getWifiStatus(w http.ResponseWriter, _ *http.Request) {

	log.Debug("In getWifiStatus")
	var message string
	args := []string{"wlan0", "status"}
	var sc int
	rc, err := a.ExecCmd("/usr/local/etc/init.d/wifi", args)
	if strings.Contains(rc, "wpa_supplicant running") {
		message = "wpa_supplicant running"
		sc = 200
	} else {
		message = "wpa_supplicant not running"
		sc = 404
	}

	a.FormatResponse(w, "getWifiStatus", sc, message, "", err)

}

func (a *App) getWifiSSID(w http.ResponseWriter, _ *http.Request) {

	log.Debug("In getWifiSSID")
	var message string
	args := []string{"-r"}
	SSID, err := a.ExecCmd("iwgetid", args)
	var sc int
	data := ""

	if SSID == "" {
		sc = 404
		message = "No SSID found"
	} else {
		sc = 200
		//message = `{ "SSID": "` + strings.TrimSpace(SSID) + `" }`
		message = "SSID found"
		data = `"SSID": "` + strings.TrimSpace(SSID) + `"`
	}

	a.FormatResponse(w, "getWifiSSID", sc, message, data, err)

}
