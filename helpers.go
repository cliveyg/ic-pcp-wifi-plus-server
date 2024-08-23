package main

import (
	"bufio"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

func (a *App) enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS, POST, PUT, DELETE")
	(*w).Header().Add("Access-Control-Allow-Credentials", "true")
	(*w).Header().Add("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
}

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
	//log.Debugf("Hashed is %s", hashed)
	return string(hashed)
}

func passMatch(wd *WifiDetails, err *error, sa *[]string) (bool, bool) {

	var hashedp string
	networkFound := false
	pm := false

	file, ferr := os.Open(os.Getenv("KNOWNWIFIFILE"))
	if ferr != nil {
		*err = ferr
		return false, false
	}
	defer func(file *os.File) {
		*err = file.Close()
		if *err != nil {
			log.Fatal(err)
		}
	}(file)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		knownWifi := strings.Split(line, "+")
		if knownWifi[0] == wd.BSSID {
			hashedp = knownWifi[2]
			*err = bcrypt.CompareHashAndPassword([]byte(hashedp), []byte(wd.Password))
			if *err == nil {
				// passwords match
				*sa = append(*sa, line)
				networkFound = true
				pm = true
			}
			// passwords don't match but network found so encrypt new pass
			// reformat line and append to sa. also use wd ssid in case this
			// has changed too
			networkFound = true
			hashedp = encryptPass(wd, err)
			editedLine := knownWifi[0] + "+" + wd.SSID + "+" + hashedp
			*sa = append(*sa, editedLine)
		} else {
			*sa = append(*sa, line)
		}
	}
	return pm, networkFound
}

func savedToTempNetConf(wd *WifiDetails, err *error) bool {

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

	// create or overwrite file with ending of .backup
	dst, fErr1 := os.Create(os.Getenv("KNOWNWIFIFILE") + ".backup")
	if fErr1 != nil {
		*err = fErr1
		return false
	}
	defer func(dst *os.File) {
		err := dst.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(dst)

	// open original file
	src, fErr2 := os.Open(os.Getenv("KNOWNWIFIFILE"))
	if fErr2 != nil {
		*err = fErr2
		return false
	}
	defer func(src *os.File) {
		err := src.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(src)

	// copy source (original file) to destination (.backup)
	bytesCopied, fErr3 := io.Copy(dst, src)
	if fErr3 != nil {
		*err = fErr3
		return false
	}

	// open temp file created by savedToTempNetConf
	src, fErr2 = os.Open(os.Getenv("KNOWNWIFIFILE") + ".temp")
	if fErr2 != nil {
		*err = fErr2
		return false
	}
	defer func(src *os.File) {
		err := src.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(src)

	// truncate orig file
	dst, fErr1 = os.Create(os.Getenv("KNOWNWIFIFILE"))
	if fErr1 != nil {
		*err = fErr1
		return false
	}
	defer func(dst *os.File) {
		err := dst.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(dst)

	// copy source (.temp) to destination (original)
	bytesCopied, fErr3 = io.Copy(dst, src)
	if fErr3 != nil {
		*err = fErr3
		return false
	}
	log.Debugf("Copied %d bytes from .temp version to original file", bytesCopied)

	// delete .temp file
	*err = os.Remove(os.Getenv("KNOWNWIFIFILE") + ".temp")
	if *err != nil {
		return false
	}
	return true
}

func restoreFromBackup() bool {

	// truncate original file
	dst, fErr1 := os.Create(os.Getenv("KNOWNWIFIFILE"))
	if fErr1 != nil {
		log.Fatal(fErr1)
	}
	defer func(dst *os.File) {
		err := dst.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(dst)

	// open backup file
	src, fErr2 := os.Open(os.Getenv("KNOWNWIFIFILE") + ".backup")
	if fErr2 != nil {
		log.Fatal(fErr2)
	}
	defer func(src *os.File) {
		err := src.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(src)

	// copy source (backup) to destination (original)
	bytesCopied, fErr3 := io.Copy(dst, src)
	if fErr3 != nil {
		log.Fatal(fErr3)
	}
	log.Debugf("Restored %d bytes from .backup version to original file", bytesCopied)
	return true
}
