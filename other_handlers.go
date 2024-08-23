package main

import (
	"encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

// ----------------------------------------------------------------------------

func (a *App) testTings(w http.ResponseWriter, r *http.Request) {

	a.enableCors(&w)
	pr := WifiPlusResponse{
		Function:   "testTings",
		Cmd:        "whatevs",
		Action:     "testy testy test",
		StatusCode: 418,
		Message:    "testing whoami stuff",
	}
	wd := WifiDetails{}
	err := json.NewDecoder(r.Body).Decode(&wd)
	if err != nil {
		pr.StatusCode = 400
		pr.Message = "Incorrect input"
		pr.Data = Eek{Error: err.Error()}
		pr.ReturnResponse(w, err)
		return
	}

	hash := encryptPass(&wd, &err)

	wd.Password = hash
	pr.Data = wd
	pr.ReturnResponse(w, err)
}

func (a *App) wpSwitcher(w http.ResponseWriter, r *http.Request) {

	log.Debug("Before enabling CORS in wpSwitcher")
	a.enableCors(&w)

	log.Debug("After enabling CORS in wpSwitcher")

	pr := WifiPlusResponse{
		Function:   "wpSwitcher",
		Action:     "switcheroo",
		StatusCode: 200,
		Message:    "Attempting to switch modes",
	}
	si := SwitcherInfo{}
	wifiRunning := false
	wapRunning := false
	pcpInitdFileExists := false

	var err error
	var pc string
	var rc []byte
	ssi := ShellSwitchInfo{}

	log.Debug("[[[[[[[ z ]]]]]]]")
	a.sysPCPConfig(&pr, r.Method, &err, &pc)
	md := textToMap(pc)

	if _, err := os.Stat(os.Getenv("TCEPCPLOC")); errors.Is(err, os.ErrNotExist) {
		log.Debug("[[[[[[[ a ]]]]]]]")
		err = errors.New("Unable to switch. apmode is not installed")
		pr.ReturnResponse(w, err)
		return
	}

	log.Debug("[[[[[[[ y ]]]]]]]")
	err = a.wapConfig(&pr, r.Method, nil)
	if err != nil {
		log.Debug("[[[[[[[ b ]]]]]]]")
		pr.ReturnResponse(w, err)
		return
	}

	si.APStatus = 200
	si.APMode = md["APMODE"]
	si.Wifi = md["WIFI"]
	wc, ok := pr.Data.(WAPConfig)
	if !ok {
		log.Fatal("Unable to cast pr.Data to variable")
	}
	si.APAddress = wc.APIPAddress

	a.wifiStatus(&pr, &err)
	ws := pr.StatusCode

	if md["WIFI"] == "on" && ws == 200 {
		wifiRunning = true
	}
	if _, err := os.Stat(os.Getenv("PCPSH")); err == nil {
		log.Debug("[[[[[[[ c ]]]]]]]")
		pcpInitdFileExists = true
	}
	if si.APMode == "yes" && pcpInitdFileExists {
		wapRunning = true
	}

	if wifiRunning && wapRunning {
		err = errors.New("Both wifi and wap are running")
		log.Debug("[[[[[[[ d ]]]]]]]")
		pr.ReturnResponse(w, err)
		return
	} else if !wifiRunning && !wapRunning {
		err = errors.New("Both wifi and wap are not running")
		log.Debug("[[[[[[[ e ]]]]]]]")
		pr.ReturnResponse(w, err)
		return
	}

	if wifiRunning {
		// switch to wap
		log.Debug("[[[[[[[ SWITCHING TO WAP ]]]]]]]")
		pr.Message = "Switching to wap"
		pr.Cmd = "nohup ./wp-switcher.sh towap"
		rc, err = exec.Command("sh", "-c", "cd /mnt/UserData/ic-pcp-wifi-plus-server/pcp-scripts; nohup ./wp-switcher.sh towap").Output()
		//rc, err = exec.Command("sh", "-c", "cd /mnt/UserData/ic-pcp-wifi-plus-server/pcp-scripts; ./wp-switcher.sh towap").Output()
		if err != nil {
			log.Debug("[[[[[[[ f ]]]]]]]")
			log.Debug(err)
			pr.ReturnResponse(w, err)
			return
		}

		log.Debugf("RC is %s", string(rc))

	} else if wapRunning {
		// switch to wifi
		log.Debug("[[[[[[[ SWITCHING TO WIFI ]]]]]]]")
		pr.Message = "Switching to wifi"
		pr.Cmd = "nohup ./wp-switcher.sh towifi"
		rc, err = exec.Command("sh", "-c", "cd /mnt/UserData/ic-pcp-wifi-plus-server/pcp-scripts; nohup ./wp-switcher.sh towifi").Output()
		//rc, err = exec.Command("sh", "-c", "cd /mnt/UserData/ic-pcp-wifi-plus-server/pcp-scripts; ./wp-switcher.sh towifi").Output()
		if err != nil {
			pr.ReturnResponse(w, err)
			return
		}
	}

	err = json.Unmarshal(rc, &ssi)
	if err != nil {
		log.Fatal(err)
	}

	pr.Data = si
	pr.ReturnResponse(w, err)
}

func (a *App) getWPACliStatus(w http.ResponseWriter, _ *http.Request) {

	a.enableCors(&w)
	pr := WifiPlusResponse{
		Function: "getWPACliStatus",
		Action:   "wpa_cli",
		Cmd:      "wpa_cli status",
		Message:  "Getting wifi status from wpa_cli",
	}

	rc, err := exec.Command("sh", "-c", "wpa_cli status").Output()
	if err != nil {
		pr.ReturnResponse(w, err)
		return
	}
	lines := strings.Split(strings.TrimSpace(string(rc)), "\n")
	// remove first line
	lines = append(lines[:0], lines[1:]...)

	wpaData := WPACliResponse{}
	wpaData.OrganiseData(lines)

	pr.StatusCode = 200
	pr.Data = wpaData
	pr.ReturnResponse(w, err)

}

func (a *App) return404(w http.ResponseWriter, _ *http.Request) {
	a.enableCors(&w)
	pr := WifiPlusResponse{
		Function:   "return404",
		Action:     "rest",
		StatusCode: 404,
		Message:    "Nowt ere chap",
	}
	pr.ReturnResponse(w, nil)
}
