package os_Specifics

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func GetLinuxConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".backdrop.yaml"), nil
}

func GetLinuxConfigPath() (string, error) {
	homePath, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	configPath := filepath.Join(homePath, ".config", "backdrop", "wallpapers")
	return configPath, nil
}

func ConfigureSlideShowLinux(images []string, wallpapersPath string, duration int) (string, error) {
	slideShowFile, slideShowConfigFile, err := createSlideShowDirectory()

	if err != nil {
		return "", err
	}

	if err := createSlideShowFile(slideShowFile, slideShowConfigFile); err != nil {
		return "", err
	}

	slideShowConfigFile, err = createSlideShowConfigFile(images, slideShowConfigFile, wallpapersPath, duration)

	if err != nil {
		return "", err
	}

	return slideShowConfigFile, nil
}

func createSlideShowDirectory() (string, string, error) {
	homePath, err := os.UserHomeDir()

	if err != nil {
		return "", "", fmt.Errorf("unable to find user home directory: %w", err)
	}

	return configureLinuxSlideShowPaths(homePath)
}

func configureLinuxSlideShowPaths(homePath string) (string, string, error) {
	paths := map[string]string{
		"slideShowPath":       filepath.Join(homePath, ".local", "share", "gnome-background-properties"),
		"slideShowConfigPath": filepath.Join(homePath, ".local", "share", "backgrounds", "backdrop_settings"),
	}

	for _, path := range paths {
		if err := os.MkdirAll(path, 0777); err != nil {
			return "", "", fmt.Errorf("failed to create directory %s: %w", path, err)
		}
	}

	slideShowFile := filepath.Join(paths["slideShowPath"], "backdrop_slideshow.xml")
	slideShowConfigFile := filepath.Join(paths["slideShowConfigPath"], "backdrop_settings.xml")

	return slideShowFile, slideShowConfigFile, nil
}

func createSlideShowFile(outFile, configFile string) error {
	content := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
	<!DOCTYPE wallpapers SYSTEM "gnome-wp-list.dtd">
	<wallpapers>
	  <wallpaper>
		    <name>Backdrop Slideshow</name>
				    <filename>%v</filename>
						    <options>zoom</options>
								    <pcolor>#2c001e</pcolor>
										    <scolor>#2c001e</scolor>
												    <shade_type>solid</shade_type>
														  </wallpaper>
															</wallpapers>
															`, configFile)

	file, err := os.Create(outFile)

	if err != nil {
		return fmt.Errorf("failed to create slideshow file at %s: %w", outFile, err)
	}
	defer file.Close()

	if _, err := file.WriteString(content); err != nil {
		return fmt.Errorf("failed to write content to slideshow file %s: %w", outFile, err)
	}

	if err := os.Chmod(outFile, 0777); err != nil {
		return fmt.Errorf("failed to set permissions for slideshow file %s: %w", outFile, err)
	}

	return nil
}

func createSlideShowConfigFile(images []string, configFile, wallpapersPath string, duration int) (string, error) {
	var content strings.Builder

	content.WriteString(`<background>
<starttime>
<year>2012</year>
<month>01</month>
<day>01</day>
<hour>00</hour>
<minute>00</minute>
<second>00</second>
</starttime>
`)

	writeStatic := func(path string) {
		fmt.Fprintf(&content, "  <static>\n    <duration>%d.0</duration>\n    <file>%s</file>\n  </static>\n", duration, path)
	}

	writeTransition := func(from, to string) {
		fmt.Fprintf(&content, "  <transition>\n    <duration>0.5</duration>\n    <from>%s</from>\n    <to>%s</to>\n  </transition>\n", from, to)
	}

	for i := 0; i < len(images)-1; i++ {
		current := filepath.Join(wallpapersPath, images[i])
		next := filepath.Join(wallpapersPath, images[i+1])
		writeStatic(current)
		writeTransition(current, next)
	}

	// Wrap around to the start
	start := filepath.Join(wallpapersPath, images[0])
	end := filepath.Join(wallpapersPath, images[len(images)-1])
	writeStatic(end)
	writeTransition(end, start)

	content.WriteString("</background>\n")

	file, err := os.Create(configFile)
	if err != nil {
		return "", fmt.Errorf("failed to create slideshow config file at %s: %w", configFile, err)
	}
	defer file.Close()

	if _, err := file.WriteString(content.String()); err != nil {
		return "", fmt.Errorf("failed to write content to config file %s: %w", configFile, err)
	}

	if err := os.Chmod(configFile, 0777); err != nil {
		return "", fmt.Errorf("failed to set permissions for config file %s: %w", configFile, err)
	}

	return file.Name(), nil
}
