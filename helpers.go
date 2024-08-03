package main

import (
	log "github.com/sirupsen/logrus"
	"os/exec"
)

func (a *App) ExecCmd(command string, args []string) (string, error) {

	log.Info("Starting ExecCmd")

	stdout, err := exec.Command(command, args...).Output()
	if err != nil {
		return "", err
	}
	return string(stdout), nil

}
