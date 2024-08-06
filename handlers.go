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

	var rc []byte
	var err error
	pr := WifiPlusResponse{
		Method: "sysAction",
		Action: sysAction,
	}

	switch sysAction {
	case "status":
		a.sysStatus(w, &pr, rc, err)
	case "picore":
		a.sysPiCoreDetails(w, &pr, rc, err)
	case "reboot":
		a.sysReboot(w, &pr)
		/*
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

		*/
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

		var netArr []WifiNetwork
		for i := 0; i < len(lines); i++ {
			wifiDetails := strings.Split(lines[0], "\t")
			wn := WifiNetwork{SSID: wifiDetails[4],
				BSSID: wifiDetails[0],
				Flags: wifiDetails[3]}
			netArr = append(netArr, wn)
		}
		pr.Data = netArr

	case "status":
		args = []string{"wlan0", "status"}
		pr.Cmd = "/usr/local/etc/init.d/wifi"
		statret, err := a.ExecCmd("/usr/local/etc/init.d/wifi", args)
		if err != nil {
			pr.ReturnResponse(w, err)
		}
		statuses := strings.Split(statret, "\n")
		pr.Message = "init.d/wifi wlan0 status"
		pr.Data = WifiStatus{
			WPASupplicantStatus: statuses[0],
			UDHCPStatus:         statuses[1],
		}
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
			pr.Data = SSID{SSID: sr}
		}

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
