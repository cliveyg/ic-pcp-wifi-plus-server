package main

import (
	"bufio"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"os"
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

func encryptPass(wd *WifiDetails, err *error) string {
	var hashed []byte
	hashed, *err = bcrypt.GenerateFromPassword([]byte(wd.Password), 8)
	log.Debugf("Hash is %s", hashed)
	return string(hashed)
}

func passMatch(wd *WifiDetails, err *error) (bool, bool) {

	var hashedp string
	networkFound := false

	file, ferr := os.Open(os.Getenv("KNOWNWIFIFILE"))
	if ferr != nil {
		*err = ferr
		return false, false
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		knownWifi := strings.Split(scanner.Text(), "+")
		if knownWifi[0] == wd.BSSID {
			hashedp = knownWifi[2]
			networkFound = true
			log.Debugf("Orig hashed pass from file is [%s]", hashedp)
		}
	}

	*err = bcrypt.CompareHashAndPassword([]byte(hashedp), []byte(wd.Password))
	if *err == nil {
		return true, networkFound
	}
	return false, networkFound
}

func savedToNetConf(wd *WifiDetails, err *error) bool {
	f, ferr := os.OpenFile(os.Getenv("KNOWNWIFIFILE"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if ferr != nil {
		*err = ferr
		return false
	}
	hashedp := encryptPass(wd, err)
	if *err != nil {
		return false
	}
	line := wd.BSSID + "+" + wd.SSID + "+" + hashedp
	if _, ferr = f.Write([]byte(line)); err != nil {
		*err = ferr
		return false
	}
	if ferr = f.Close(); err != nil {
		*err = ferr
		return false
	}
	return true
}
