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

type SelectionOptions struct {
	Prompt         string
	SuccessMessage string
	Cleanup        func()
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

func handleSelectionConfirmation(previousWallpaper string, out io.Writer, opts *SelectionOptions) (bool, error) {
	prompt := opts.Prompt
	if prompt == "" {
		prompt = "Want to save this change? [y/N]: "
	}

	successMessage := opts.SuccessMessage
	if successMessage == "" {
		successMessage = "Successfully changed background image!"
	}
	for {
		userInput, err := userConfirmationWithPrompt(inputConfirmation, prompt)
		if err != nil {
			return false, err
		}

		switch userInput {
		case "y":
			fmt.Fprintln(out, successMessage)
			return true, nil
		case "n", "":
			if err := setWallpaper(previousWallpaper); err != nil {
				return false, err
			}

			opts.Cleanup()
			return false, nil
		default:
			fmt.Fprintln(out, "Invalid input...")
		}
	}
}

func userConfirmationWithPrompt(r io.Reader, prompt string) (string, error) {
	reader := bufio.NewReader(r)
	fmt.Print(prompt)

	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	input = strings.ToLower(strings.TrimSpace(input))
	return input, nil
}
