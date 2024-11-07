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

const (
	gnomeSchema = "org.gnome.desktop.background"
	mateSchema  = "org.mate.desktop.background"
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

func getWallpapers(path string) ([]string, error) {
	fileEntries, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory %s: %w", path, err)
	}

	files := make([]string, 0, len(fileEntries))
	for _, file := range fileEntries {
		files = append(files, file.Name())
	}

	return files, nil
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
		return getPreviousWallpaperLinux()
	case "windows":
		return getPreviousWallpaperWindows()
	}
	return "", fmt.Errorf("unsupported operating system: %s", runtime.GOOS)

}

func getPreviousWallpaperLinux() (string, error) {
	if !commandExist("gsettings") {
		return "", fmt.Errorf("gssettings not available")
	}

	schemas, err := listSchemas()
	if err != nil {
		return "", err
	}

	if strings.Contains(schemas.String(), gnomeSchema) {
		return getGsettingsWallpaper(gnomeSchema)
	}

	if strings.Contains(schemas.String(), mateSchema) {
		return getGsettingsWallpaper(mateSchema)
	}

	return "", ErrNoCompatibleDesktopEnvironment

}

func getPreviousWallpaperWindows() (string, error) {
	if !commandExist("powershell") {
		return "", fmt.Errorf("powershell not available on Windows")
	}
	cmdGetPicture := exec.Command("powershell", "-Command", "(Get-ItemProperty -Path 'HKCU:\\Control Panel\\Desktop' -Name Wallpaper).Wallpaper")
	var outGetPicture bytes.Buffer
	cmdGetPicture.Stdout = &outGetPicture
	err := cmdGetPicture.Run()
	if err != nil {
		return "", err
	}
	wallpaperPath := strings.TrimSpace(outGetPicture.String())
	return wallpaperPath, nil
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
	psCommand := fmt.Sprintf(`Add-Type -TypeDefinition @'
using System;
using System.Runtime.InteropServices;
public class Wallpaper {
	[DllImport("user32.dll", CharSet = CharSet.Auto)]
	public static extern int SystemParametersInfo(int uAction, int uParam, string lpvParam, int fuWinIni);
	public static void SetWallpaper(string path) {
		 SystemParametersInfo(20, 0, path, 0x01 | 0x02);
	}
}
'@; [Wallpaper]::SetWallpaper("%s")`, wallpaper)

	cmd := exec.Command("powershell", "-command", psCommand)
	if err := cmd.Run(); err != nil {
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
