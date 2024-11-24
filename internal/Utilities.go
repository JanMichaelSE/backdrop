package internal

import (
	"os/exec"
)

func commandExist(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}
