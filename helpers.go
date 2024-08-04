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

func (a *App) FormatResponse(w http.ResponseWriter, cmd string, sc int, message string, data string, err error) {

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err != nil {
		log.Error("Error is %s", err)
		sc = 500
		jsonStr := "{ \"command\": \"" + cmd + "\", " +
			"\"message\": \"error\"," +
			"\"data\": {" + err.Error() + "} }"

		w.WriteHeader(sc)
		if _, err := io.WriteString(w, jsonStr); err != nil {
			log.Fatal(err)
		}
		return
	}

	jsonStr := "{ \"command\": \"" + cmd + "\", " +
		"\"message\": \"" + message + "\"," +
		"\"data\": {" + data + "} }"

	w.WriteHeader(sc)
	if _, err := io.WriteString(w, jsonStr); err != nil {
		log.Fatal(err)
	}

}
