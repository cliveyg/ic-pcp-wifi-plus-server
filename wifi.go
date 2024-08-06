package main

import (
	log "github.com/sirupsen/logrus"
	"net/http"
	"os/exec"
	"strings"
)

func (a *App) wifiSSID(w http.ResponseWriter, pr *WifiPlusResponse, err error) {
	args := []string{"-r"}
	var sr string
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
}

func (a *App) wifiStatus(w http.ResponseWriter, pr *WifiPlusResponse, err error) {

	args := []string{"wlan0", "status"}
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

}

func (a *App) wifiScan(w http.ResponseWriter, pr *WifiPlusResponse, err error) {

	var rc []byte
	pr.StatusCode = 200
	pr.Message = "Searching for networks..."
	pr.Cmd = "wpa_cli scan wlan0; wpa_cli scan_results"
	rc, err = exec.Command("sh", "-c", "wpa_cli scan wlan0; wpa_cli scan_results").Output()
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
