package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func configureSlideShow(imageText, wallpapersPath string, duration int) (string, error) {
	images := strings.Split(imageText, ";")

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
		return "", "", err
	}

	slideShowPath := filepath.Join(homePath, ".local/share/gnome-background-properties")
	_, err = os.Stat(slideShowPath)
	if err != nil {
		os.MkdirAll(slideShowPath, 0777)
	}

	slideShowConfigPath := filepath.Join(homePath, ".local/share/backgrounds/backdrop_settings")
	_, err = os.Stat(slideShowConfigPath)
	if err != nil {
		os.MkdirAll(slideShowConfigPath, 0777)
	}

	slideShowFile := filepath.Join(slideShowPath, "backdrop_slideshow.xml")
	slideShowConfigFile := filepath.Join(slideShowConfigPath, "backdrop_settings.xml")

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
		return err
	}
	defer file.Close()

	if _, err := file.WriteString(content.String()); err != nil {
		return err
	}

	if err := os.Chmod(outFile, 0777); err != nil {
		return err
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
		return "", err
	}
	defer file.Close()

	if _, err := file.WriteString(content.String()); err != nil {
		return "", err
	}

	if err := os.Chmod(configFile, 0777); err != nil {
		return "", err
	}

	return file.Name(), nil
}
