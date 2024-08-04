package main

import (
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os/exec"
)

func (a *App) ExecCmd(command string, args []string) (string, error) {

	log.Info("Starting ExecCmd")

	stdout, err := exec.Command(command, args...).Output()
	if err != nil {
		log.Info("Something went wrong %s", err)
		return "", err
	}
	log.Info("Something went right %s", string(stdout))
	return string(stdout), nil

}

func (a *App) FormatResponse(w http.ResponseWriter, rc string, err error) {
	if err != nil {
		log.Error("Error is %s", err)
		mess := `{"error": "` + err.Error() + `", "rc": "` + rc + `"}`
		w.WriteHeader(500)
		if _, err := io.WriteString(w, mess); err != nil {
			log.Fatal(err)
		}
		return
	}

	mess := `{"message": "` + rc + `"}`
	if _, err := io.WriteString(w, mess); err != nil {
		log.Fatal(err)
	}

}
