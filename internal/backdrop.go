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
	imageSelection ImageSelection
	input          io.Reader = os.Stdin
)

func BackdropAction(out io.Writer, config *Config, args []string) error {
	if config.path != "" {
		err := configureImagePath(config.path)
		if err != nil {
			return err
		}
	}

	wallpapersPath, err := getUserImagesPath()
	if err != nil {
		return err
	}

	wallpapers := getWallpapers(wallpapersPath)

	switch {
	case config.isFuzzy:
		imageSelection = fuzzySelection
	}

outter:
	for {
		previousWallpaper, err := getPreviousWallpaper()
		if err != nil {
			return err
		}

		selectedWallpaper, err := imageSelection(wallpapers)
		if err != nil {
			return err
		}

		fullSelectedPath := filepath.Join(wallpapersPath, selectedWallpaper)
		if stats, err := os.Stat(fullSelectedPath); err == nil && stats.Mode().IsRegular() {
			err := setWallpaper(fullSelectedPath)
			if err != nil {
				return err
			}
		}

	inner:
		for {
			userInput, err := userConfirmation(input)
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

	// TODO: Still not sure what to return here.
	// To a certain point I feel it's gonna stay 'nil'
	return nil
}
