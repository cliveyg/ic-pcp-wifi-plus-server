package main

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os/exec"
)

func (a *App) wapStopStart(w http.ResponseWriter, pr *WifiPlusResponse, ac string) {
	log.Debugf("In wapStopStart and our action is [%s]", ac)
	pr.Cmd = "nowt yet"
	pr.StatusCode = 200
	pr.Function = "wapStopStart"
	pr.Message = fmt.Sprintf("Action is [%s]", ac)
	pr.ReturnResponse(w, nil)
}

func (a *App) wapAddRemove(w http.ResponseWriter, pr *WifiPlusResponse, m string) {

	if m == http.MethodPost {

		var rc []byte
		pr.Cmd = "wifi-plus.sh wp_wap_add"
		rc, err := exec.Command("sh", "-c", "cd cgi-bin && ./wifi-plus.sh wp_wap_add").Output()
		if err != nil {
			pr.ReturnResponse(w, err)
		}

		pr.StatusCode = 200
		pr.Message = "Installing wap extensions"
		var b map[string]interface{}
		err = json.Unmarshal(rc, &b)
		if err != nil {
			log.Fatal(err)
		}
		pr.Data = b

	} else if m == http.MethodDelete {
		log.Debug("We should be removing the ap mode stuff")
	} else {
		pr.StatusCode = 405
		pr.Function = "wapAddRemove"
		pr.Cmd = "meep"
		err := fmt.Errorf("HTTP method [%s] not valid for this resource", m)
		pr.ReturnResponse(w, err)
	}
}
