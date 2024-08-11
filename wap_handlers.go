package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"os/exec"
)

func (a *App) wapAction(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	wa := vars["action"]

	//TODO: Check input string more thoroughly

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
	case "config":
		err = a.wapConfig(w, &pr, r.Method, r.Body)
	default:
		// do nowt
		pr.StatusCode = 400
		pr.Message = "Action does not exist"
	}

	log.WithFields(log.Fields{"Full response is ": pr}).Debug()
	pr.ReturnResponse(w, err)

}

func (a *App) wapConfig(w http.ResponseWriter, pr *WifiPlusResponse, hm string, bd io.ReadCloser) error {
	log.Debugf("In wapConfig and our action is [%s]", hm)

	pr.Function = "wapConfig"
	if hm == http.MethodPut {
		var cf WAPConfig
		pr.Cmd = "wifi-plus.sh wp_wap_edit_config"
		err := json.NewDecoder(bd).Decode(&cf)
		if err != nil {
			pr.StatusCode = 400
			pr.Message = "Incorrect input"
			pr.Data = Eek{Error: err.Error()}
			return err
		}
		cf.ValidateInput(&err)
		if err != nil {
			pr.StatusCode = 400
			pr.Message = "Failed validation"
			pr.Data = Eek{Error: err.Error()}
			return err
		}
		// sending the wap settings as a single string to script
		sCmd := "cd cgi-bin && ./wifi-plus.sh wp_edit_wap_config " + cf.Stringify()
		rc, er2 := exec.Command("sh", "-c", sCmd).Output()
		if er2 != nil {
			return er2
		}
		var b map[string]interface{}
		err = json.Unmarshal(rc, &b)
		if err != nil {
			log.Fatal(err)
		}
		pr.Message = "Successfully updated WAP config"
	} else {
		// only GET
		pr.Cmd = "wifi-plus.sh wp_fetch_wap_config"
		rc, err := exec.Command("sh", "-c", "cd cgi-bin && ./wifi-plus.sh wp_fetch_wap_config").Output()
		if err != nil {
			return err
		}
		log.Debugf("RC is [%s]", string(rc))
		wapCfg := WAPConfig{}
		err = json.Unmarshal(rc, &wapCfg)
		if err != nil {
			log.Debug("WQRWERWERWERWRWERWER")
			log.Fatal(err)
		}
		pr.Message = "WAP config details"
		pr.Data = wapCfg
	}
	pr.StatusCode = 200
	return nil
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

	pr := WifiPlusResponse{
		Function: "wapInfo",
		Action:   r.Method,
	}
	pr.Cmd = "/usr/local/etc/init.d/pcp-apmode status"
	args := []string{"status"}
	var sr string
	sr, err := a.ExecCmd("/usr/local/etc/init.d/pcp-apmode", args)
	if err != nil {
		// create our own error due to returned error missing info
		err = errors.New("Not all wap processes running")
		pr.ReturnResponse(w, err)
		return
	}

	pr.StatusCode = 200
	pr.Message = sr
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
	var rc []byte

	if r.Method == http.MethodPost {

		pr.Cmd = "wifi-plus.sh wp_wap_add"
		if _, err := os.Stat("/usr/local/etc/init.d/pcp-apmode"); err == nil {
			pr.StatusCode = 409
			pr.Message = "WAP is already installed"
			pr.ReturnResponse(w, nil)
			return
		}

		rc, err = exec.Command("sh", "-c", "cd cgi-bin && ./wifi-plus.sh wp_wap_add").Output()
		if err != nil {
			pr.ReturnResponse(w, err)
			return
		}
		log.Debugf("RC from [wifi-plus.sh wp_wap_add] is %s", string(rc))
		pr.StatusCode = 200
		pr.Message = "Installing wap extensions"
		r := `{"noop": "doop"}`
		var b map[string]interface{}
		err = json.Unmarshal([]byte(r), &b)
		if err != nil {
			log.Fatal(err)
		}
		pr.Data = b

	} else if r.Method == http.MethodDelete {
		log.Debug("We should be removing the ap mode stuff")
		pr.Cmd = "wifi-plus.sh wp_wap_remove"
		pr.Message = "Deleting wap extensions"

		rc, err = exec.Command("sh", "-c", "cd cgi-bin && ./wifi-plus.sh wp_wap_remove").Output()
		if err != nil {
			pr.ReturnResponse(w, err)
			return
		}
		log.Debugf("RC from [wifi-plus.sh wp_wap_remove] is %s", string(rc))
		pr.StatusCode = http.StatusGone
		pr.Message = "Removing wap extensions"
		r := `{"doop": "noop"}`
		var b map[string]interface{}
		err = json.Unmarshal([]byte(r), &b)
		if err != nil {
			log.Fatal(err)
		}
		pr.Data = b

	}
	pr.ReturnResponse(w, err)
}
