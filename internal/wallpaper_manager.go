package internal

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/viper"
)

func configureWallpaperPath(path string) error {
	homePath, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	viper.Set("WallpapersPath", path)
	configPath := filepath.Join(homePath, ".backdrop.yaml")
	if err := viper.WriteConfigAs(configPath); err != nil {
		return fmt.Errorf("%w : %v", ErrCouldNotConfigureImagePath, err)
	}
	return nil
}

func getUserWallpapersPath() (string, error) {
	homePath, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	wallpapersPath, ok := viper.Get("WallpapersPath").(string)
	if ok {
		if stat, err := os.Stat(wallpapersPath); err == nil && stat.IsDir() {
			return wallpapersPath, nil
		}
	}

	configPath := filepath.Join(homePath, ".config", "backdrop", "wallpapers")
	if stat, err := os.Stat(configPath); err == nil && stat.IsDir() {
		return configPath, nil
	}

	picturesPath := filepath.Join(homePath, "Pictures", "wallpapers")
	if stat, err := os.Stat(picturesPath); err == nil && stat.IsDir() {
		return picturesPath, nil
	}

	return "", ErrNoValidImagesPath
}

func getWallpapers(path string) []string {
	fileEntries, _ := os.ReadDir(path)

	files := make([]string, 0, len(fileEntries))
	for _, file := range fileEntries {
		files = append(files, file.Name())
	}

	return files
}

func listSchemas() (*bytes.Buffer, error) {
	cmd := exec.Command("gsettings", "list-schemas")
	var outListSchemas bytes.Buffer
	cmd.Stdout = &outListSchemas
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("%w : %v", ErrCouldNotListSchemas, err)
	}

	return &outListSchemas, nil
}

func getPreviousWallpaper() (string, error) {
	switch runtime.GOOS {
	case "linux":
		schemas, err := listSchemas()
		if err != nil {
			return "", err
		}

		if strings.Contains(schemas.String(), "gnome.desktop.background") {
			cmdGetPicture := exec.Command("gsettings", "get", "org.gnome.desktop.background", "picture-uri")
			var outGetPicture bytes.Buffer
			cmdGetPicture.Stdout = &outGetPicture
			err := cmdGetPicture.Run()
			if err != nil {
				return "", err
			}

			uri := strings.ReplaceAll(strings.Trim(outGetPicture.String(), "\n"), "'", "")
			if strings.Contains(uri, "://") {
				parts := strings.SplitN(uri, "://", 2)
				if len(parts) == 2 {
					return parts[1], nil
				}
			}

			return uri, nil
		}

		if strings.Contains(schemas.String(), "mate.desktop.background") {
			cmdGetPicture := exec.Command("gsettings", "get", "org.mate.desktop.background", "picture-uri")
			var outGetPicture bytes.Buffer
			cmdGetPicture.Stdout = &outGetPicture
			err := cmdGetPicture.Run()
			if err != nil {
				return "", err
			}

			uri := strings.ReplaceAll(strings.Trim(outGetPicture.String(), "\n"), "'", "")
			if strings.Contains(uri, "://") {
				parts := strings.SplitN(uri, "://", 2)
				if len(parts) == 2 {
					return parts[1], nil
				}
			}

			return uri, nil
		}

	case "windows":
		cmdGetPicture := exec.Command("powershell", "-Command", "Get-ItemProperty -Path 'HKCU:\\Control Panel\\Desktop\\' -Name Wallpaper")
		var outGetPicture bytes.Buffer
		cmdGetPicture.Stdout = &outGetPicture
		err := cmdGetPicture.Run()
		if err != nil {
			return "", err
		}
		wallpaperPath := strings.TrimSpace(outGetPicture.String())
		return wallpaperPath, nil
	}
	return "", ErrNoCompatibleDesktopEnvironment

}

func setWallpaper(wallpaper string) error {
	switch runtime.GOOS {
	case "linux":
		return setWallpaperLinux(wallpaper)
	case "windows":
		return setWallpaperWindows(wallpaper)
	}
	return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)

}

func setWallpaperWindows(wallpaper string) error {
	psCommand := fmt.Sprintf(`RUNDLL32.EXE user32.dll, UpdatePerUserSystemParameters ,1, True; $path = "%s"; [SystemParametersInfo]::SystemParametersInfo(20, 0, $path, 3)`, wallpaper)
	cmdSetPicture := exec.Command("powershell", "-Command", psCommand)
	if err := cmdSetPicture.Run(); err != nil {
		return fmt.Errorf("%w : %v", ErrCouldNotSetBackground, err)
	}
	return nil
}

func setWallpaperLinux(wallpaper string) error {
	schemas, err := listSchemas()
	if err != nil {
		return err
	}

	wallpaperURI := fmt.Sprintf("file://%s", wallpaper)

	if strings.Contains(schemas.String(), "gnome.desktop.background") {
		cmdSetPicture := exec.Command("gsettings", "set", "org.gnome.desktop.background", "picture-uri", wallpaperURI)
		if err := cmdSetPicture.Run(); err != nil {
			return fmt.Errorf("%w : %v", ErrCouldNotSetBackground, err)
		}

		cmdSetPictureDark := exec.Command("gsettings", "set", "org.gnome.desktop.background", "picture-uri-dark", wallpaperURI)
		if err := cmdSetPictureDark.Run(); err != nil {
			return fmt.Errorf("%w : %v", ErrCouldNotSetBackground, err)
		}

		return nil
	}

	if strings.Contains(schemas.String(), "mate.desktop.background") {
		cmdSetPicture := exec.Command("gsettings", "set", "org.mate.desktop.background", "picture-uri", wallpaperURI)
		err := cmdSetPicture.Run()
		if err != nil {
			return fmt.Errorf("%w : %v", ErrCouldNotSetBackground, err)
		}

		return nil
	}

	return ErrNoCompatibleDesktopEnvironment
}
