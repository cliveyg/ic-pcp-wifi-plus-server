package main

import (
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os/exec"
	"strings"
)

// ----------------------------------------------------------------------------

func (a *App) testTings(w http.ResponseWriter, _ *http.Request) {

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	rc, err := exec.Command("sh", "-c", "cd cgi-bin && sudo ./wifi-plus.sh wp_test").Output()
	if err != nil {
		log.Error("Error is %s", err)
		mess := `{"error": "` + err.Error() + `", "rc": "` + strings.TrimSpace(string(rc)) + `"}`
		w.WriteHeader(500)
		if _, err := io.WriteString(w, mess); err != nil {
			log.Fatal(err)
		}
		return
	}

	mess := `{"message": "` + string(rc) + `"}`
	if _, err := io.WriteString(w, mess); err != nil {
		log.Fatal(err)
	}
}

func (a *App) getSystemStatus(w http.ResponseWriter, _ *http.Request) {

	log.Debug("In getSystemStatus")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	rc, err := exec.Command("sh", "-c", "cd cgi-bin && sudo ./wifi-plus.sh wp_status 200").Output()
	if err != nil {
		log.Error("Error is %s", err)
		mess := `{"error": "` + err.Error() + `", "rc": "` + strings.TrimSpace(string(rc)) + `"}`
		w.WriteHeader(500)
		if _, err := io.WriteString(w, mess); err != nil {
			log.Fatal(err)
		}
		return
	}

	mess := `{"message": "System running...", "return_code": "` + strings.TrimSpace(string(rc)) + `"}`
	if _, err := io.WriteString(w, mess); err != nil {
		log.Fatal(err)
	}

}

func (a *App) getWifiStatus(w http.ResponseWriter, _ *http.Request) {

	log.Debug("In getWifiStatus")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var message string
	args := []string{"wlan0", "status"}
	rc, err := a.ExecCmd("/usr/local/etc/init.d/wifi", args)
	if err != nil {
		w.WriteHeader(500)
		mess := `{"error": "` + err.Error() + `"}`
		if _, err := io.WriteString(w, mess); err != nil {
			log.Fatal(err)
		}
		return
	}
	if strings.Contains(rc, "wpa_supplicant running") {
		message = `{"command": "wifi status", "message": "wpa_supplicant running" }`
	} else {
		message = `{"command": "wifi status", "message": "wpa_supplicant not running"}`
	}
	if _, err := io.WriteString(w, message); err != nil {
		log.Fatal(err)
	}
}

func (a *App) getWifiSSID(w http.ResponseWriter, _ *http.Request) {

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var message string
	args := []string{"-r"}
	SSID, err := a.ExecCmd("iwgetid", args)
	if err != nil {
		w.WriteHeader(500)
		mess := `{"error": "` + err.Error() + `"}`
		if _, err := io.WriteString(w, mess); err != nil {
			log.Fatal(err)
		}
		return
	}
	if SSID == "" {
		w.WriteHeader(404)
		message = `{, "message": "No SSID found" }`
	} else {
		message = `{ "SSID": "` + strings.TrimSpace(SSID) + `" }`
	}
	if _, err := io.WriteString(w, message); err != nil {
		log.Fatal(err)
	}
}
