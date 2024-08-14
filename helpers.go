package main

import (
	"bufio"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"io"
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

func passMatch(wd *WifiDetails, err *error, sa *[]string) (bool, bool) {

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
			*err = bcrypt.CompareHashAndPassword([]byte(hashedp), []byte(wd.Password))
			log.Debugf("Orig hashed pass from file is [%s]", hashedp)
			if *err == nil {
				// passwords match
				*sa = append(*sa, scanner.Text())
				return true, networkFound
			}
			// pass no match but network found so encrypt new pass
			// reformat line and append to sa
			networkFound = true
			hashedp := encryptPass(wd, err)
			if *err != nil {
				return false, networkFound
			}
			editedLine := knownWifi[0] + "+" + knownWifi[1] + "+" + hashedp
			*sa = append(*sa, editedLine)
		} else {
			*sa = append(*sa, scanner.Text())
		}
	}

	return false, networkFound
}

func savedToNewNetConf(wd *WifiDetails, err *error) bool {

	var sa []string
	passMatch(wd, err, &sa)

	f, ferr := os.OpenFile(os.Getenv("KNOWNWIFIFILE")+".temp", os.O_CREATE|os.O_WRONLY, 0644)
	if ferr != nil {
		*err = ferr
		return false
	}
	for _, line := range sa {
		if _, ferr = f.Write([]byte(line)); err != nil {
			*err = ferr
			return false
		}
	}
	if ferr = f.Close(); err != nil {
		*err = ferr
		return false
	}

	return true
}

func fileSwitch(err *error) bool {

	*err = os.Remove(os.Getenv("KNOWNWIFIFILE"))
	if *err != nil {
		return false
	}
	// open new version of file
	src, fer1 := os.Open(os.Getenv("KNOWNWIFIFILE") + ".new")
	if fer1 != nil {
		*err = fer1
		return false
	}
	defer src.Close()

	// create old filename
	dst, fer2 := os.Create(os.Getenv("KNOWNWIFIFILE"))
	if fer2 != nil {
		*err = fer2
		return false
	}
	defer dst.Close()

	// copy source to destination
	bytesCopied, fer3 := io.Copy(dst, src)
	if fer3 != nil {
		*err = fer3
		return false
	}
	log.Debugf("Copied %d bytes from new version to old filename", bytesCopied)

	// delete .new file
	*err = os.Remove(os.Getenv("KNOWNWIFIFILE"))
	if *err != nil {
		return false
	}
	return true
}
