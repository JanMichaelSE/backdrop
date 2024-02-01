package internal

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/ktr0731/go-fuzzyfinder"
)

type ImageSelection func(s []string) (string, error)
type GetImageSelector func(c *Config) ImageSelection

func getImageSelector(c *Config) ImageSelection {
	switch {
	case c.isSlideShow:
		return multiFuzzySelection
	default:
		return fuzzySelection
	}
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

func userDuration(r io.Reader) (int, error) {
	reader := bufio.NewReader(r)
	fmt.Print("What should be the duration per slide? (In Seconds): ")

	input, err := reader.ReadString('\n')
	if err != nil {
		return 0, err
	}

	input = strings.ReplaceAll(input, "\n", "")
	duration, err := strconv.Atoi(input)
	if err != nil {
		return 0, err
	}

	return duration, nil
}
