package internal

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/ktr0731/go-fuzzyfinder"
)

type ImageSelection func([]string) (string, error)

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
