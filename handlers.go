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

	log.Debug("-----------------------------")
	log.Debug("In testTings")

	pr := WifiPlusResponse{
		Method:     "testTings",
		Cmd:        "whatevs",
		Action:     "testy testy test",
		StatusCode: 200,
		Message:    "tings"}

	/*
		_, err := exec.Command("sh", "-c", "cd /mnt/UserData/industrialcool-pcp-wifi-plus/pcp-scripts; nohup ./wifi-plus.sh wp_wifi_restart > /dev/null 2>&1 &").Output()
		if err != nil {
			log.Fatal(err)
		}

	*/

	pr.Data = "\"boopy\": \"beep\""
	pr.ReturnResponse(w, nil)
}

func (a *App) systemAction(w http.ResponseWriter, r *http.Request) {

	log.Debug("-----------------------------")
	log.Debug("In systemAction")
	vars := mux.Vars(r)
	sysAction := vars["action"]

	//TODO: Check input string more thoroughly

	var rc []byte
	var rcInt int
	var err error
	pr := WifiPlusResponse{
		Method: "sysAction",
		Action: sysAction,
	}

	switch sysAction {
	case "status":
		pr.Cmd = "wifi-plus.sh wp_status 200"
		rc, err = exec.Command("sh", "-c", "cd cgi-bin && ./wifi-plus.sh wp_status 200").Output()
		if err != nil {
			pr.ReturnResponse(w, err)
		}
		rcInt, err = strconv.Atoi(strings.TrimSpace(string(rc)))
		if err != nil {
			pr.ReturnResponse(w, err)
		}
		pr.StatusCode = rcInt
		pr.Message = "System running"

	case "picore":
		pr.Cmd = "wifi-plus.sh wp_picore_details"
		rc, err = exec.Command("sh", "-c", "cd cgi-bin && sudo ./wifi-plus.sh wp_picore_details").Output()
		if err != nil {
			pr.ReturnResponse(w, err)
		}
		pr.StatusCode = 200
		pr.Message = "piCore details"
		pr.Data = string(rc)

	case "reboot":
		pr.StatusCode = 202
		pr.Message = "System rebooting"
		pr.Cmd = "sudo pcp rb"
		pr.ReturnResponse(w, nil)
		time.Sleep(2 * time.Second)
		rc, err := exec.Command("sh", "-c", "sudo pcp rb").Output()
		log.Debug(rc)
		if err != nil {
			pr.ReturnResponse(w, err)
		}
		return

	default:
		// do nowt
		pr.StatusCode = 400
		pr.Message = "Action does not exist"
	}

	pr.ReturnResponse(w, err)
}

func (a *App) wifiAction(w http.ResponseWriter, r *http.Request) {

	log.Debug("-----------------------------")
	log.Debug("In wifiAction")
	vars := mux.Vars(r)
	wifiAction := vars["action"]

	//TODO: Check input string more thoroughly

	var sr string
	var err error
	var args []string
	pr := WifiPlusResponse{
		Method: "wifiAction",
		Action: wifiAction,
	}

	switch wifiAction {
	case "restart":
		pr.StatusCode = 202
		pr.Message = "Now we wait..."
		pr.Cmd = "nohup ./wp-wifi-refresh.sh"
		_, err := exec.Command("sh", "-c", "cd /mnt/UserData/industrialcool-pcp-wifi-plus/pcp-scripts; nohup ./wp-wifi-refresh.sh > /dev/null 2>&1 &").Output()
		if err != nil {
			pr.ReturnResponse(w, err)
		}
		pr.Data = `"script called": "wp-wifi-refresh.sh"`
		pr.ReturnResponse(w, nil)
		return

	case "scan":
		pr.StatusCode = 200
		pr.Message = "Searching for networks..."
		pr.Cmd = "wpa_cli scan wlan0; wpa_cli scan_results"
		rc, err := exec.Command("sh", "-c", "wpa_cli scan wlan0; wpa_cli scan_results").Output()
		if err != nil {
			pr.ReturnResponse(w, err)
		}
		lines := strings.Split(strings.TrimSpace(string(rc)), "\n")
		// remove first 4 lines
		lines = append(lines[:0], lines[4:]...)
		log.WithFields(log.Fields{"no of wifi networks": len(lines)}).Debug()

		jsonStr := "["
		for i := 0; i < len(lines); i++ {
			wifiDetails := strings.Split(lines[0], "\t")
			wifiNum := string(rune(i))
			log.WithFields(log.Fields{"WIFI NO:": wifiNum}).Debug()
			jsonStr = jsonStr + `{"wifi ` + wifiNum + `": { "ssid": "` + wifiDetails[4] + `",` +
				`"bssid": "` + wifiDetails[0] + `",` +
				`"flags": "` + wifiDetails[3] + `"}},`
		}
		pr.Data = jsonStr + "]"
		log.WithFields(log.Fields{"pr.Data": pr.Data}).Debug()

	case "status":
		args = []string{"wlan0", "status"}
		pr.Cmd = "/usr/local/etc/init.d/wifi"
		statret, err := a.ExecCmd("/usr/local/etc/init.d/wifi", args)
		if err != nil {
			pr.ReturnResponse(w, err)
		}
		statuses := strings.Split(statret, "\n")
		pr.Message = "init.d/wifi wlan0 status"
		pr.Data = `"wpa_supplicant status": "` + statuses[0] + `", "udhcpc status" : "` + statuses[0] + `"`
		if strings.Contains(statret, "not running") {
			pr.StatusCode = 404
		} else {
			pr.StatusCode = 200
		}

	case "ssid":
		args = []string{"-r"}
		pr.Cmd = "iwgetid"
		sr, err = a.ExecCmd("iwgetid", args)
		if err != nil {
			pr.ReturnResponse(w, err)
		}
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
		pr.Message = "Action does not exist"
	}

	pr.ReturnResponse(w, err)

}

func (a *App) getWPACliStatus(w http.ResponseWriter, _ *http.Request) {

	log.Debug("-----------------------------")
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
		Method:     "getWPACliStatus",
		Action:     "wpa_cli",
		StatusCode: 200,
		Message:    "wpa_cli status",
		Data:       string(jsonData)}
	pr.ReturnResponse(w, err)

}
