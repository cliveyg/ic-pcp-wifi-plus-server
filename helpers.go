package main

import (
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"os/exec"
	"strings"
)

func (a *App) ExecCmd(command string, args []string) (string, error) {

	stdout, err := exec.Command(command, args...).Output()
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Debug("Something went wrong")
		return "", err
	}
	return strings.TrimSpace(string(stdout)), nil

}

func textToMap(sg string) map[string]string {

	output := map[string]string{}
	for _, pair := range strings.Split(strings.TrimSpace(sg), "\n") {
		kv := strings.Split(pair, "=")
		rs := strings.ReplaceAll(kv[1], "\"", "")
		output[kv[0]] = rs
	}
	return output
}

func encryptPass(wd *WifiDetails, err *error) {
	var hashed []byte
	hashed, *err = bcrypt.GenerateFromPassword([]byte(wd.Password), 8)
	log.Debugf("Hash is %s", hashed)
	wd.Password = string(hashed)
}

func passMatch(wd *WifiDetails, hp string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hp), []byte(wd.Password))
	if err == nil {
		return true
	}
	return false
}
