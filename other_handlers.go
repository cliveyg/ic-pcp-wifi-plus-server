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

func (a *App) testTings(w http.ResponseWriter, _ *http.Request) {

	log.Debug("-----------------------------")
	log.Debug("In testTings")

	pr := WifiPlusResponse{
		Function:   "testTings",
		Cmd:        "whatevs",
		Action:     "testy testy test",
		StatusCode: 418,
		Message:    "testing whoami stuff",
	}
	/*
		pr.Cmd = "nohup ./wp-switcher.sh"
		//r, err := exec.Command("sh", "-c", "cd /mnt/UserData/industrialcool-pcp-wifi-plus/pcp-scripts; nohup ./wp-switcher.sh > /dev/null 2>&1 &").Output()
		r, err := exec.Command("sh", "-c", "cd /mnt/UserData/industrialcool-pcp-wifi-plus/pcp-scripts; ./wp-switcher.sh").Output()
		if err != nil {
			log.Debug("(((((( 1 ))))))")
			pr.ReturnResponse(w, err)
			return
		}

	*/
	pr.Cmd = "./wifi-plus.sh wp_test"
	rc, err := exec.Command("sh", "-c", "cd cgi-bin && ./wifi-plus.sh wp_test").Output()
	if err != nil {
		log.Debug("[[[[[ 0 ]]]]]")
		pr.ReturnResponse(w, err)
		return
	}
	log.Debugf("r is [%s]", string(rc))
	src := `{ "mess": "` + strings.TrimSpace(string(rc)) + `"}`
	var b map[string]interface{}
	err = json.Unmarshal([]byte(src), &b)
	if err != nil {
		log.Debug("[[[[[ 1 ]]]]]")
		log.Fatal(err)
	}
	pr.Data = b
	pr.ReturnResponse(w, nil)
	/*
		r := `{"boopy": "beep"}`
		var b map[string]interface{}
		err := json.Unmarshal([]byte(r), &b)
		if err != nil {
			log.Fatal()
		}
		pr.Data = b
		pr.ReturnResponse(w, nil)
	*/
}

func (a *App) wpSwitcher(w http.ResponseWriter, r *http.Request) {

	log.Debug("wpSwitcher - attempting to switch between wifi and wap")

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
	//args := []string{"status"}

	a.sysPCPConfig(&pr, r.Method, &err, &pc)
	md := textToMap(pc)

	if _, err := os.Stat(os.Getenv("TCEPCPLOC")); errors.Is(err, os.ErrNotExist) {
		err = errors.New("Unable to switch. apmode is not installed")
		pr.ReturnResponse(w, err)
		return
	}
	/*
		_, err = a.ExecCmd("/usr/local/etc/init.d/pcp-apmode", args)
		if err != nil {
			err = errors.New("Unable to switch. apmode is not installed")
			pr.ReturnResponse(w, err)
			return
		}

	*/

	err = a.wapConfig(&pr, http.MethodGet, nil)
	if err != nil {
		pr.ReturnResponse(w, err)
		return
	}

	si.APStatus = 200
	si.APMode = md["APMODE"]
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
		pcpInitdFileExists = true
	}
	if si.APMode == "on" && pcpInitdFileExists {
		wapRunning = true
	}

	if wifiRunning && wapRunning {
		err = errors.New("Both wifi and wap are running")
		pr.ReturnResponse(w, err)
		return
	} else if !wifiRunning && !wapRunning {
		err = errors.New("Both wifi and wap are not running")
		pr.ReturnResponse(w, err)
		return
	}

	if wifiRunning {
		// switch to wap
		pr.Message = "Switch to wap"
	} else if wapRunning {
		// switch to wifi
		pr.Message = "switch to wifi"
	}

	pr.Cmd = ""
	pr.Data = si
	pr.ReturnResponse(w, err)
}

func (a *App) getWPACliStatus(w http.ResponseWriter, _ *http.Request) {

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
	pr := WifiPlusResponse{
		Function:   "return404",
		Action:     "rest",
		StatusCode: 404,
		Message:    "Nowt ere chap",
	}
	pr.ReturnResponse(w, nil)
}
