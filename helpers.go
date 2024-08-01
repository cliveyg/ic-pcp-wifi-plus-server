package main

import "os/exec"

func (a *App) ExecCmd(command string, args []string) (string, error) {

	stdout, err := exec.Command(command, args...).Output()
	if err != nil {
		return "", err
	}
	return string(stdout), nil

}
