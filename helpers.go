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
		log.WithFields(log.Fields{"err": err}).Debug("Something went wrong")
		return "", err
	}
	log.WithFields(log.Fields{"stdout": string(stdout)}).Debug("Something went right")
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

func textToMap(sg string) map[string]string {

	log.Debug("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@")
	log.Debug("In textToMap")
	output := map[string]string{}
	for _, pair := range strings.Split(strings.TrimSpace(sg), "\n") {
		log.Debugf("PAIR IS %s", pair)
		kv := strings.Split(pair, "=")
		rs := strings.ReplaceAll(kv[1], "\"", "")
		output[kv[0]] = rs
	}
	log.WithFields(log.Fields{"[[[output]]]": output}).Debug()
	return output
}
