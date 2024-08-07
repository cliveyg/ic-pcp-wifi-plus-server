package main

import (
	"encoding/json"
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
		Function:   "testTings",
		Cmd:        "whatevs",
		Action:     "testy testy test",
		StatusCode: 200,
		Message:    "tings",
	}

	r := `{"boopy": "beep"}`
	var b map[string]interface{}
	err := json.Unmarshal([]byte(r), &b)
	if err != nil {
		log.Fatal()
	}
	pr.Data = b
	pr.ReturnResponse(w, nil)
}

func (a *App) systemAction(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	sa := vars["action"]

	//TODO: Check input string more thoroughly

	var err error
	pr := WifiPlusResponse{
		Function: "sysAction",
		Action:   sa,
	}

	switch sa {
	case "picore":
		a.sysPiCoreDetails(w, &pr)
	case "reboot":
		a.sysReboot(w, &pr)
		return
	case "shutdown":
		a.sysShutdown(w, &pr)
		return
	case "status":
		a.sysStatus(w, &pr)
	default:
		// do nowt
		pr.StatusCode = 400
		pr.Message = "Action does not exist"
	}

	pr.ReturnResponse(w, err)
}

func (a *App) wapAction(w http.ResponseWriter, r *http.Request) {

	log.Debug("-----------------------------")
	log.Debug("In wapAction")
	vars := mux.Vars(r)
	wa := vars["action"]

	//TODO: Check input string more thoroughly
	log.Debug(r.Method)
	var err error
	pr := WifiPlusResponse{
		Function: "wapAction",
		Action:   wa,
	}

	switch wa {
	case "stop", "start":
		if r.Method == http.MethodGet {
			a.wapStopStart(w, &pr, wa)
		} else {
			pr.StatusCode = 405
			pr.Message = "Incorrect method for action"
		}
	default:
		// do nowt
		pr.StatusCode = 400
		pr.Message = "Action does not exist"
	}

	log.WithFields(log.Fields{"Full response is ": pr}).Debug()
	pr.ReturnResponse(w, err)

}

func (a *App) wifiAction(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	wa := vars["action"]

	//TODO: Check input string more thoroughly

	var err error
	pr := WifiPlusResponse{
		Function: "wifiAction",
		Action:   wa,
	}

	switch wa {
	case "restart":
		a.wifiRestart(w, &pr)
		return
	case "scan":
		a.wifiScan(w, &pr)
	case "ssid":
		a.wifiSSID(w, &pr)
	case "status":
		a.wifiStatus(w, &pr)
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
		Function:   "getWPACliStatus",
		Action:     "wpa_cli",
		StatusCode: 200,
		Message:    "wpa_cli status",
		Data:       wpaData}
	pr.ReturnResponse(w, err)

}

func (a *App) return404(w http.ResponseWriter, _ *http.Request) {
	pr := WifiPlusResponse{
		Function:   "return404",
		Action:     "rest",
		StatusCode: 404,
		Message:    "Nowt ere chap",
	}
	pr.ReturnResponse(w, nil)
}
