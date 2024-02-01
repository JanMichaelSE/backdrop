package internal

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type Config struct {
	path     string
	imageUrl string
	// TODO: Make Fuzzy finding the default
	isFuzzy     bool
	isSlideShow bool
}

func NewConfig(path, imageUrl string, isFuzzy, isSlideShow bool) *Config {
	return &Config{
		path:        path,
		imageUrl:    imageUrl,
		isFuzzy:     isFuzzy,
		isSlideShow: isSlideShow,
	}
}

var (
	getSelector       GetImageSelector = getImageSelector
	inputConfirmation io.Reader        = os.Stdin
	inputDuration     io.Reader        = os.Stdin
)

func BackdropAction(out io.Writer, config *Config, args []string) error {
	if config.path != "" {
		err := configureWallpaperPath(config.path)
		if err != nil {
			return err
		}
	}

	wallpapersPath, err := getUserWallpapersPath()
	if err != nil {
		return err
	}

	wallpapers := getWallpapers(wallpapersPath)

outter:
	for {
		previousWallpaper, err := getPreviousWallpaper()
		if err != nil {
			return err
		}

		imageSelection := getSelector(config)
		selectedWallpaper, err := imageSelection(wallpapers)
		if err != nil {
			return err
		}

		if config.isSlideShow {
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
		}

		if selectedWallpaper[0] != '/' {
			fmt.Println("I did not happen")
			fullSelectedPath := filepath.Join(wallpapersPath, selectedWallpaper)
			stats, err := os.Stat(fullSelectedPath)
			if err == nil && stats.Mode().IsRegular() {
				err := setWallpaper(fullSelectedPath)
				if err != nil {
					return err
				}
			}
			if err != nil {
				return err
			}
		}

	inner:
		for {
			userInput, err := userConfirmation(inputConfirmation)
			if err != nil {
				return err
			}

			switch userInput {
			case "y":
				fmt.Fprintln(out, "Successfully changed background image!")
				break outter
			case "n":
				setWallpaper(previousWallpaper)
				if err != nil {
					return err
				}
				break inner
			default:
				fmt.Fprintln(out, "Invalid input...")
			}
		}

	}

	return nil
}
