package internal

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

func commandExist(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func getGsettingsWallpaper(schema string) (string, error) {
	cmd := exec.Command("gsettings", "get", schema, "picture-uri")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return "", err
	}

	uri := strings.ReplaceAll(strings.Trim(out.String(), "\n"), "'", "")
	if strings.Contains(uri, "://") {
		parts := strings.SplitN(uri, "://", 2)
		if len(parts) == 2 {
			return parts[1], nil
		}
		return "", fmt.Errorf("unexpected URI format: %s", uri)
	}
	return uri, nil
}
