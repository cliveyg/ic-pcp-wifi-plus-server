package main

import (
	"encoding/json"
	"net/http"
	"os/exec"
)

func (a *App) wapInstall(w http.ResponseWriter, pr *WifiPlusResponse, err error) {

	var rc []byte
	pr.Cmd = "wifi-plus.sh wp_wap_install"
	rc, err = exec.Command("sh", "-c", "cd cgi-bin && ./wifi-plus.sh wp_wap_install").Output()
	if err != nil {
		pr.ReturnResponse(w, err)
	}

	pr.StatusCode = 200
	pr.Message = "Installing wap extensions"
	var b map[string]interface{}
	json.Unmarshal(rc, &b)
	pr.Data = b

}
