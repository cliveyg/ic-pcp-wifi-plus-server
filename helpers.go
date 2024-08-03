package main

import (
	log "github.com/sirupsen/logrus"
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
