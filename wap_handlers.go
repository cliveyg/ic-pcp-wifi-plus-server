package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/exec"
)

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
		if _, err := os.Stat("/usr/local/etc/init.d/pcp-apmode"); errors.Is(err, os.ErrNotExist) {
			pr.StatusCode = 404
			pr.Message = "WAP mode is not installed"
		} else {
			if r.Method == http.MethodGet {
				a.wapStopStart(w, &pr, wa)
			} else {
				pr.StatusCode = 405
				pr.Message = "Incorrect method for action"
			}
		}
	default:
		// do nowt
		pr.StatusCode = 400
		pr.Message = "Action does not exist"
	}

	log.WithFields(log.Fields{"Full response is ": pr}).Debug()
	pr.ReturnResponse(w, err)

}

func (a *App) wapStopStart(w http.ResponseWriter, pr *WifiPlusResponse, ac string) {
	log.Debugf("In wapStopStart and our action is [%s]", ac)
	pr.Cmd = "nowt yet"
	pr.StatusCode = 200
	pr.Function = "wapStopStart"
	pr.Message = fmt.Sprintf("Action is [%s]", ac)
	pr.ReturnResponse(w, nil)
}

func (a *App) wapInfo(w http.ResponseWriter, r *http.Request) {
	log.Debugf("In wapInfo and our action is [%s]", r.Method)
	pr := WifiPlusResponse{
		Function: "wapInfo",
		Action:   r.Method,
	}
	pr.Cmd = "nowt yet"
	pr.StatusCode = 200
	pr.Message = fmt.Sprintf("Action is [%s]", r.Method)
	pr.ReturnResponse(w, nil)
}

func (a *App) wapAddRemove(w http.ResponseWriter, r *http.Request) {

	// http 'post' is to add the tcz files, 'delete' is to remove and
	//'get' is fetch the current details if installed
	log.Debug(r.Method)
	pr := WifiPlusResponse{
		Function: "wapAddRemove",
		Action:   r.Method,
	}
	var err error
	//var rc []byte

	if r.Method == http.MethodPost {

		log.Debug("We get here")
		pr.Cmd = "wifi-plus.sh wp_wap_add"
		_, err = exec.Command("sh", "-c", "cd cgi-bin && ./wifi-plus.sh wp_wap_add").Output()
		if err != nil {
			log.Info("==============================")
			pr.ReturnResponse(w, err)
		}

		pr.StatusCode = 200
		pr.Message = "Installing wap extensions"
		r := `{"noop": "doop"}`
		var b map[string]interface{}
		err = json.Unmarshal([]byte(r), &b)
		if err != nil {
			log.Info("+++++++++++++++++++++++++++++++++")
			log.Fatal(err)
		}
		pr.Data = b

	} else if r.Method == http.MethodDelete {
		log.Debug("We should be removing the ap mode stuff")
		pr.Cmd = "wifi-plus.sh wp_wap_remove"
		pr.Message = "Deleting wap extensions"
		pr.StatusCode = http.StatusGone
	} else {
		pr.StatusCode = 405
		pr.Function = "wapAddRemove"
		pr.Cmd = "meep"
		err = fmt.Errorf("HTTP method [%s] not valid for this resource", r.Method)
	}
	pr.ReturnResponse(w, err)
}
