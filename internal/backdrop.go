package internal

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type Config struct {
	path        string
	isImageUrl  bool
	isFuzzy     bool
	isSlideShow bool
}

func NewConfig(path string, isImageUrl, isSlideShow bool) *Config {
	return &Config{
		path:        path,
		isImageUrl:  isImageUrl,
		isSlideShow: isSlideShow,
	}
}

var (
	getSelector       GetFuzzySelector = getFuzzySelector
	inputConfirmation io.Reader        = os.Stdin
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

	wallpapers, err := getWallpapers(wallpapersPath)
	if err != nil {
		return err
	}
	switch {
	case config.isSlideShow:
		imageSelection := getSelector(config)
		err := handleSlideshow(out, wallpapersPath, wallpapers, imageSelection)
		if err != nil {
			return err
		}
	case config.isImageUrl:
		err := handleImageUrl(out, wallpapersPath)
		if err != nil {
			return err
		}
	default:
		imageSelection := getSelector(config)
		err := handleFuzzySearch(out, wallpapersPath, wallpapers, imageSelection)
		if err != nil {
			return err
		}
	}

	return nil
}

func handleSelectionConfirmation(previousWallpaper string, out io.Writer, cleanup func()) (bool, error) {
	for {
		userInput, err := userConfirmation(inputConfirmation)
		if err != nil {
			return false, err
		}

		switch userInput {
		case "y":
			fmt.Fprintln(out, "Successfully changed background image!")
			return true, nil
		case "n", "":
			if err := setWallpaper(previousWallpaper); err != nil {
				return false, err
			}

			cleanup()
			return false, nil
		default:
			fmt.Fprintln(out, "Invalid input...")
		}
	}
}

func userConfirmation(r io.Reader) (string, error) {
	reader := bufio.NewReader(r)
	fmt.Print("Want to save this change? [y/N]: ")

	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	input = strings.ToLower(strings.TrimSpace(input))
	return input, nil
}
