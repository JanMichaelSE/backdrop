package internal

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/ktr0731/go-fuzzyfinder"
)

type FuzzySelection func(s []string) (string, error)
type GetFuzzySelector func(c *Config) FuzzySelection

func getFuzzySelector(c *Config) FuzzySelection {
	switch {
	case c.isSlideShow:
		return multiFuzzySelection
	default:
		return fuzzySelection
	}
}

func handleFuzzySearch(out io.Writer, wallpapersPath string, wallpapers []string, imageSelection FuzzySelection) error {
	for hasConfirmed := false; !hasConfirmed; {
		previousWallpaper, err := getPreviousWallpaper()
		if err != nil {
			return err
		}

		selectedWallpaper, err := imageSelection(wallpapers)
		if err != nil {
			return err
		}

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

		hasConfirmed, err = handleSelectionConfirmation(previousWallpaper, out, &SelectionOptions{
			Prompt:         "",
			SuccessMessage: "",
			Cleanup:        func() {},
		})

		if err != nil {
			return err
		}
	}

	return nil
}

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

func multiFuzzySelection(fileNames []string) (string, error) {
	selectedIndexes, err := fuzzyfinder.FindMulti(
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

	images := make([]string, 0, len(selectedIndexes))
	for _, index := range selectedIndexes {
		images = append(images, fileNames[index])
	}

	return strings.Join(images, ";"), nil
}
