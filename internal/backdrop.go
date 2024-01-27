package internal

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/ktr0731/go-fuzzyfinder"
)

type Config struct {
	path        string
	imageUrl    string
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

func BackdropAction(out io.Writer, config *Config, args []string) error {
	wallpapersPath, err := getUserImagesPath()
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

		var selectedWallpaper string
		if config.isFuzzy {
			selectedWallpaper, err = fuzzySelection(wallpapers)
			if err != nil {
				return err
			}
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
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Want to save this change? [y/N]: ")

			input, err := reader.ReadString('\n')
			if err != nil {
				return err
			}

			input = strings.ToLower(strings.TrimSpace(input))
			switch input {
			case "y":
				fmt.Println("Successfully changed background image!")
				break outter
			case "n":
				setWallpaper(previousWallpaper)
				if err != nil {
					return err
				}
				break inner
			default:
				fmt.Println("Invalid input...")
			}
		}

	}

	// TODO: Still not sure what to return here.
	// To a certain point I feel it's gonna stay 'nil'
	return nil
}

// TODO: Have fuzzy finding at the top
//   - For the time being keep at bottom,
//     need to see if either I can contribute to the package to support this.
//   - Or investigate alternative solutions and see if BubbleTea LipGloss Go packages support this.
func fuzzySelection(fileNames []string) (string, error) {
	selectedIndex, err := fuzzyfinder.Find(
		fileNames,
		func(i int) string {
			return fileNames[i]
		},
		fuzzyfinder.WithCursorPosition(1),
	)
	if errors.Is(err, fuzzyfinder.ErrAbort) {
		return "", ErrUserCanceledSelection
	}
	if err != nil {
		return "", err
	}

	return fileNames[selectedIndex], nil
}
