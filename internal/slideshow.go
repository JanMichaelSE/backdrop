package internal

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

var (
	inputDuration io.Reader = os.Stdin
)

func handleSlideshow(out io.Writer, wallpapersPath string, wallpapers []string, imageSelection FuzzySelection) error {
	for hasConfirmed := false; !hasConfirmed; {
		previousWallpaper, err := getPreviousWallpaper()
		if err != nil {
			return err
		}

		selectedWallpaper, err := imageSelection(wallpapers)
		if err != nil {
			return err
		}

		duration, err := userDuration(inputDuration)
		if err != nil {
			return err
		}

		selectedWallpaper, err = configureSlideShow(selectedWallpaper, wallpapersPath, duration)
		if err != nil {
			return err
		}

		err = setWallpaper(selectedWallpaper)
		if err != nil {
			return err
		}

		hasConfirmed, err = handleSelectionConfirmation(previousWallpaper, out, func() {})
		if err != nil {
			return err
		}
	}

	return nil
}

func userDuration(r io.Reader) (int, error) {
	reader := bufio.NewReader(r)
	fmt.Print("What should be the duration per slide? (In Seconds): ")

	input, err := reader.ReadString('\n')
	if err != nil {
		return 0, fmt.Errorf("failed to read duration input: %w", err)
	}

	input = strings.TrimSpace(input)
	duration, err := strconv.Atoi(input)
	if err != nil || duration <= 0 {
		return 0, fmt.Errorf("invalid duration: please enter a positive integer")
	}

	return duration, nil
}

func configureSlideShow(imageText, wallpapersPath string, duration int) (string, error) {
	images := strings.Split(imageText, ";")

	switch runtime.GOOS {
	case "linux":
		return configureSlideShowLinux(images, wallpapersPath, duration)
	case "windows":
		return "", fmt.Errorf("SlideShow is currently not supported for Windows.")

	default:
		return "", ErrNoCompatibleOS
	}
}

func configureSlideShowLinux(images []string, wallpapersPath string, duration int) (string, error) {

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

	switch runtime.GOOS {
	case "linux":
		return configureLinuxSlideShowPaths(homePath)
	case "windows":
		return "", "", fmt.Errorf("SlideShow is currently not supported for Windows.")

	default:
		return "", "", ErrNoCompatibleOS
	}
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
	var content strings.Builder

	content.WriteString(`<?xml version="1.0" encoding="UTF-8"?>`)
	content.WriteString("\n")
	content.WriteString(`<!DOCTYPE wallpapers SYSTEM "gnome-wp-list.dtd">`)
	content.WriteString("\n")
	content.WriteString("<wallpapers>\n")
	content.WriteString("  <wallpaper>\n")
	content.WriteString("    <name>Backdrop Slideshow</name>\n")
	content.WriteString(fmt.Sprintf("     <filename>%v</filename>\n", configFile))
	content.WriteString("    <options>zoom</options>\n")
	content.WriteString("    <pcolor>#2c001e</pcolor>\n")
	content.WriteString("    <scolor>#2c001e</scolor>\n")
	content.WriteString("    <shade_type>solid</shade_type>\n")
	content.WriteString("  </wallpaper>\n")
	content.WriteString("</wallpapers>\n")

	file, err := os.Create(outFile)
	if err != nil {
		return fmt.Errorf("failed to create slideshow file at %s: %w", outFile, err)
	}
	defer file.Close()

	if _, err := file.WriteString(content.String()); err != nil {
		return fmt.Errorf("failed to write content to slideshow file %s: %w", outFile, err)
	}

	if err := os.Chmod(outFile, 0777); err != nil {
		return fmt.Errorf("failed to set permissions for slideshow file %s: %w", outFile, err)
	}

	return nil
}

func createSlideShowConfigFile(images []string, configFile, wallpapersPath string, duration int) (string, error) {
	var content strings.Builder

	content.WriteString("<background>\n")
	content.WriteString("  <starttime>\n")
	content.WriteString("    <year>2012</year>\n")
	content.WriteString("    <month>01</month>\n")
	content.WriteString("    <day>01</day>\n")
	content.WriteString("    <hour>00</hour>\n")
	content.WriteString("    <minute>00</minute>\n")
	content.WriteString("    <second>00</second>\n")
	content.WriteString("  </starttime>\n")

	for i := 0; i < len(images)-1; i++ {
		currentImage := filepath.Join(wallpapersPath, images[i])
		nextImage := filepath.Join(wallpapersPath, images[i+1])
		content.WriteString("  <static>\n")
		content.WriteString(fmt.Sprintf("   <duration>%d.0</duration>\n", duration))
		content.WriteString(fmt.Sprintf("   <file>%v</file>\n", currentImage))
		content.WriteString("  </static>\n")
		content.WriteString("   <transition>\n")
		content.WriteString("   <duration>0.5</duration>\n")
		content.WriteString(fmt.Sprintf("   <from>%v</from>\n", currentImage))
		content.WriteString(fmt.Sprintf("   <to>%v</to>\n", nextImage))
		content.WriteString("   </transition>\n")
	}

	startImage := filepath.Join(wallpapersPath, images[0])
	endImage := filepath.Join(wallpapersPath, images[len(images)-1])
	content.WriteString("  <static>\n")
	content.WriteString(fmt.Sprintf("   <duration>%d.0</duration>\n", duration))
	content.WriteString(fmt.Sprintf("   <file>%v</file>\n", endImage))
	content.WriteString("  </static>\n")
	content.WriteString("  <transition>\n")
	content.WriteString("   <duration>0.5</duration>\n")
	content.WriteString(fmt.Sprintf("   <from>%v</from>\n", endImage))
	content.WriteString(fmt.Sprintf("   <to>%v</to>\n", startImage))
	content.WriteString("  </transition>\n")
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
