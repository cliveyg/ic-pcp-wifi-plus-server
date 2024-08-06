package main

import (
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os/exec"
	"strings"
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

	var err error
	pr := WifiPlusResponse{
		Method: "sysAction",
		Action: sysAction,
	}

	switch sysAction {
	case "status":
		a.sysStatus(w, &pr, err)
	case "picore":
		a.sysPiCoreDetails(w, &pr, err)
	case "reboot":
		a.sysReboot(w, &pr)
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

	var err error
	pr := WifiPlusResponse{
		Method: "wifiAction",
		Action: wifiAction,
	}

	switch wifiAction {
	case "restart":
		a.wifiRestart(w, &pr)
		return
	case "scan":
		a.wifiScan(w, &pr, err)
	case "status":
		a.wifiStatus(w, &pr, err)
	case "ssid":
		a.wifiSSID(w, &pr, err)
	default:
		// do nowt
		pr.StatusCode = 400
		pr.Message = "Action does not exist"
	}

	log.WithFields(log.Fields{"Full response is ": pr}).Debug()
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

	pr := WifiPlusResponse{
		Method:     "getWPACliStatus",
		Action:     "wpa_cli",
		StatusCode: 200,
		Message:    "wpa_cli status",
		Data:       wpaData}
	pr.ReturnResponse(w, err)

}
