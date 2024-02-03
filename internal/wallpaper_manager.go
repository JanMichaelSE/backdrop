package internal

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
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

	configPath := fmt.Sprintf("%v/.config/backdrop/wallpapers", homePath)
	if stat, err := os.Stat(configPath); err == nil && stat.IsDir() {
		return configPath, nil
	}

	picturesPath := fmt.Sprintf("%v/Pictures/wallpapers", homePath)
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

	return "", ErrNoCompatibleDesktopEnvironment
}

func setWallpaper(wallpaper string) error {
	schemas, err := listSchemas()
	if err != nil {
		return err
	}

	wallpaper = fmt.Sprintf("file://%s", wallpaper)

	if strings.Contains(schemas.String(), "gnome.desktop.background") {
		cmdSetPicture := exec.Command("gsettings", "set", "org.gnome.desktop.background", "picture-uri", wallpaper)
		if err := cmdSetPicture.Run(); err != nil {
			return fmt.Errorf("%w : %v", ErrCouldNotSetBackground, err)
		}

		cmdSetPictureDark := exec.Command("gsettings", "set", "org.gnome.desktop.background", "picture-uri-dark", wallpaper)
		if err := cmdSetPictureDark.Run(); err != nil {
			return fmt.Errorf("%w : %v", ErrCouldNotSetBackground, err)
		}

		return nil
	}

	if strings.Contains(schemas.String(), "mate.desktop.background") {
		cmdSetPicture := exec.Command("gsettings", "set", "org.mate.desktop.background", "picture-uri", wallpaper)
		err := cmdSetPicture.Run()
		if err != nil {
			return fmt.Errorf("%w : %v", ErrCouldNotSetBackground, err)
		}

		return nil
	}

	return ErrNoCompatibleDesktopEnvironment
}

// NOTE: Might move to a separate file later.
func downloadImage(imageUrl, wallpapersPath string) (string, error) {
	resp, err := http.Get(imageUrl)
	if err != nil {
		return "", fmt.Errorf("Could not fetch image, got error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Got bad response code '%v' from request. Cannot proceed.", resp.StatusCode)
	}

	fileName := strings.Split(imageUrl, "/")
	sanitizedFileName := sanitizeFilename(fileName[len(fileName)-1])
	filePath := filepath.Join(wallpapersPath, sanitizedFileName)
	file, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("Could not create file for image url, got error: %v", err)
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return "", fmt.Errorf("Could not copy contents from image url to file, got error: %v", err)
	}

	return filePath, nil
}

func sanitizeFilename(filename string) string {
	// URL Decode
	decodedFilename, err := url.QueryUnescape(filename)
	if err != nil {
		decodedFilename = filename // Use the original filename if decoding fails
	}

	// Define a list of invalid characters and their replacements
	replacements := map[string]string{
		"/":  "_",
		"\\": "_",
		"<":  "_",
		">":  "_",
		":":  "_",
		"\"": "_",
		"|":  "_",
		"?":  "_",
		"*":  "_",
	}

	// Replace invalid characters
	for oldChar, newChar := range replacements {
		decodedFilename = strings.ReplaceAll(decodedFilename, oldChar, newChar)
	}

	// Normalize spaces (optional)
	decodedFilename = strings.Join(strings.Fields(decodedFilename), " ")

	return decodedFilename
}
