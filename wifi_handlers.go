package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os/exec"
	"strings"
)

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
		a.wifiScan(&pr, &err)
	case "ssid":
		a.wifiSSID(&pr, &err)
	case "status":
		a.wifiStatus(&pr, &err)
	case "stop", "start":
		a.wifiStopStart(&pr, &err)
	default:
		// do nowt
		pr.StatusCode = 400
		pr.Message = "Action does not exist"
	}

	log.WithFields(log.Fields{"Full response is ": pr}).Debug()
	pr.ReturnResponse(w, err)

}

func (a *App) wifiSwitchNetwork(w http.ResponseWriter, r *http.Request) {

	log.Debug("[[[[[[[[[ a ]]]]]]]]")

	pr := WifiPlusResponse{
		Function: "wifiSwitch",
		Message:  "Switch wifi networks",
	}
	wd := WifiDetails{}
	err := json.NewDecoder(r.Body).Decode(&wd)
	if err != nil {
		log.Debug("[[[[[[[[[ y ]]]]]]]]")
		pr.StatusCode = 400
		pr.Message = "Incorrect input"
		pr.Data = Eek{Error: err.Error()}
		pr.ReturnResponse(w, err)
		return
	}
	log.Debug("[[[[[[[[[ z ]]]]]]]]")
	// check if sent wifi details match details on file
	newNet := false
	connOk := false
	pm, nf := passMatch(&wd, &err, nil)
	if pm && nf {
		log.Debug("[[[[[[[[[ c ]]]]]]]]")
		pr.Message = "Network found and passwords match"
		pr.StatusCode = 418
	} else if !pm && nf {
		log.Debug("[[[[[[[[[ d ]]]]]]]]")
		pr.StatusCode = 403
		pr.Message = "Network found but password doesn't match"
	} else {
		newNet = true
	}
	// try to switch to network
	var rc []byte
	pr.StatusCode = 200
	pr.Message = pr.Message + ". Switching networks..."
	pr.Cmd = "nohup ./wp-wifi-switch.sh"
	rc, err = exec.Command("sh", "-c", "cd /mnt/UserData/industrialcool-pcp-wifi-plus/pcp-scripts; nohup ./wp-wifi-switch.sh > /dev/null 2>&1 &").Output()
	if err != nil {
		log.Debug("[[[[[[[[[ 0 ]]]]]]]]")
		pr.ReturnResponse(w, err)
		return
	}
	sr := ShellResponse{}
	err = json.Unmarshal(rc, &sr)
	if err != nil {
		log.Debug("[[[[[[[[[ 1 ]]]]]]]]")
		log.Fatal(err)
	}
	log.Debug("[[[[[[[[[ b ]]]]]]]]")
	log.Debug(sr)
	if sr.Status == 200 {
		connOk = true
	}

	// if a new network or existing network but with new pass and connected ok then save to file
	if (newNet || (nf && !pm)) && connOk && savedToTempNetConf(&wd, &err) {
		log.Debug("[[[[[[[[[ e ]]]]]]]]")
		// saved to temp file so overwrite old file with new version
		if !fileSwitch(&err) {
			if restoreFromBackup() {
				pr.Message = "Connected but unable to create new version of conf - restored old version"
			} else {
				pr.Message = "Connected but unable to create new version of conf or restore from backup"
				err = errors.New(pr.Message)
			}
			pr.ReturnResponse(w, err)
		}

	} else if (newNet || (nf && !pm)) && connOk {
		log.Debug("[[[[[[[[[ f ]]]]]]]]")
		pr.Message = "Connected but unable to save network details to temp file"
		err = errors.New(pr.Message)
		pr.ReturnResponse(w, err)
	} else if !connOk {
		log.Debug("[[[[[[[[[ g ]]]]]]]]")
		pr.Message = fmt.Sprintf("Unable to switch to %s wifi network", wd.SSID)
		pr.ReturnResponse(w, err)
		return
	}
	log.Debug("[[[[[[[[[ h ]]]]]]]]")
	wd.Password = "********"
	pr.Data = wd
	pr.ReturnResponse(w, err)
}

func (a *App) wifiStopStart(pr *WifiPlusResponse, err *error) {
	log.Debug("Woop")
	pr.Function = "wifiStopStart"

}

func (a *App) wifiSSID(pr *WifiPlusResponse, err *error) {
	args := []string{"-r"}
	var sr string
	pr.Cmd = "iwgetid"
	sr, *err = a.ExecCmd("iwgetid", args)
	if *err != nil {
		return
	}
	if sr == "" {
		pr.StatusCode = 404
		pr.Message = "No SSID found"
	} else {
		pr.StatusCode = 200
		pr.Message = "SSID found"
		pr.Data = SSID{SSID: sr}
	}
}

func (a *App) wifiStatus(pr *WifiPlusResponse, err *error) {

	args := []string{"wlan0", "status"}
	var ret string
	pr.Cmd = "/usr/local/etc/init.d/wifi"
	ret, *err = a.ExecCmd("/usr/local/etc/init.d/wifi", args)
	if *err != nil {
		return
	}
	stats := strings.Split(ret, "\n")
	pr.Message = "init.d/wifi wlan0 status"
	pr.Data = WifiStatus{
		WPASupplicantStatus: stats[0],
		UDHCPStatus:         stats[1],
	}
	if strings.Contains(ret, "not running") {
		pr.StatusCode = 404
	} else {
		pr.StatusCode = 200
	}

}

func (a *App) wifiScan(pr *WifiPlusResponse, err *error) {

	var rc []byte
	pr.StatusCode = 200
	pr.Message = "Searching for networks..."
	pr.Cmd = "wpa_cli scan wlan0; wpa_cli scan_results"
	rc, *err = exec.Command("sh", "-c", "wpa_cli scan wlan0; sleep 3; wpa_cli scan_results").Output()
	if *err != nil {
		return
	}
	log.Debug(string(rc))
	lines := strings.Split(strings.TrimSpace(string(rc)), "\n")
	// remove first 4 lines of returned data
	lines = append(lines[:0], lines[4:]...)
	log.WithFields(log.Fields{"no of wifi networks": len(lines)}).Debug()

	var netArr []WifiNetwork
	for i := 0; i < len(lines); i++ {
		wifiDetails := strings.Split(lines[i], "\t")
		wn := WifiNetwork{SSID: wifiDetails[4],
			BSSID: wifiDetails[0],
			Flags: wifiDetails[3]}
		netArr = append(netArr, wn)
	}
	pr.Data = netArr
}

func (a *App) wifiRestart(w http.ResponseWriter, pr *WifiPlusResponse) {

	pr.StatusCode = 202
	pr.Message = "Now we wait..."
	pr.Cmd = "nohup ./wp-wifi-refresh.sh"
	_, err := exec.Command("sh", "-c", "cd /mnt/UserData/industrialcool-pcp-wifi-plus/pcp-scripts; nohup ./wp-wifi-refresh.sh > /dev/null 2>&1 &").Output()
	if err != nil {
		pr.ReturnResponse(w, err)
	}
	pr.ReturnResponse(w, nil)
}
