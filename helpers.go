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
	defer dst.Close()

	// open original file
	src, fErr2 := os.Open(os.Getenv("KNOWNWIFIFILE"))
	if fErr2 != nil {
		*err = fErr2
		return false
	}
	defer src.Close()

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
	defer src.Close()

	// truncate orig file
	dst, fErr1 = os.Create(os.Getenv("KNOWNWIFIFILE"))
	if fErr1 != nil {
		*err = fErr1
		return false
	}
	defer dst.Close()

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
	defer dst.Close()

	// open backup file
	src, fErr2 := os.Open(os.Getenv("KNOWNWIFIFILE") + ".backup")
	if fErr2 != nil {
		log.Fatal(fErr2)
	}
	defer src.Close()

	// copy source (backup) to destination (original)
	bytesCopied, fErr3 := io.Copy(dst, src)
	if fErr3 != nil {
		log.Fatal(fErr3)
	}
	log.Debugf("Restored %d bytes from .backup version to original file", bytesCopied)
	return true
}
