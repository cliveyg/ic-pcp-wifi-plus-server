package main

import (
	log "github.com/sirupsen/logrus"
	"os/exec"
	"strings"
)

func (a *App) ExecCmd(command string, args []string) (string, error) {

	log.Info("Starting ExecCmd")

	stdout, err := exec.Command(command, args...).Output()
	if err != nil {
		log.Info("Something went wrong %s", err)
		return "", err
	}
	log.Info("Something went right %s", string(stdout))
	return strings.TrimSpace(string(stdout)), nil

}

func substr(input string, start int, length int) string {
	asRunes := []rune(input)
	if start >= len(asRunes) {
		return ""
	}
	if start+length > len(asRunes) {
		length = len(asRunes) - start
	}
	return string(asRunes[start : start+length])
}
