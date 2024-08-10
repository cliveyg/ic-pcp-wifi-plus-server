package main

import (
	"encoding/json"
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
		StatusCode: 418,
		Message:    "testing whoami stuff",
	}
	/*
		pr.Cmd = "nohup ./wp-wifi-to-wap.sh"
		//r, err := exec.Command("sh", "-c", "cd /mnt/UserData/industrialcool-pcp-wifi-plus/pcp-scripts; nohup ./wp-wifi-to-wap.sh > /dev/null 2>&1 &").Output()
		r, err := exec.Command("sh", "-c", "cd /mnt/UserData/industrialcool-pcp-wifi-plus/pcp-scripts; ./wp-wifi-to-wap.sh").Output()
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
	src := `{ "mess": "'+string(rc)+'"}`
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
