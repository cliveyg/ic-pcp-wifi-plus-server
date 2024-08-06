package main

import (
	"encoding/json"
	"net/http"
	"os/exec"
)

func (a *App) piCoreDetails(w http.ResponseWriter, pr *WifiPlusResponse, rc []byte, err error) {

	pr.Cmd = "wifi-plus.sh wp_picore_details"
	rc, err = exec.Command("sh", "-c", "cd cgi-bin && sudo ./wifi-plus.sh wp_picore_details").Output()
	if err != nil {
		pr.ReturnResponse(w, err)
	}

	pr.StatusCode = 200
	pr.Message = "piCore system details"
	picoreData := PiCoreSystemData{}

	err = json.Unmarshal(rc, &picoreData)
	if err != nil {
		pr.ReturnResponse(w, err)
	}
	pr.Data = picoreData

}
