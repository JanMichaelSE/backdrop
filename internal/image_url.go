package internal

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

var (
	inputImageUrl io.Reader = os.Stdin
)

func handleImageUrl(out io.Writer, wallpapersPath string) error {
	for hasConfirmed := false; !hasConfirmed; {
		previousWallpaper, err := getPreviousWallpaper()
		if err != nil {
			return err
		}

		imageUrl, err := userImageUrl(inputImageUrl)
		if err != nil {
			return err
		}

		image, err := downloadImage(imageUrl, wallpapersPath)
		if err != nil {
			return err
		}

		err = setWallpaper(image)
		if err != nil {
			return err
		}

		imageCleanup := func() {
			os.Remove(image)
		}

		hasConfirmed, err = handleSelectionConfirmation(previousWallpaper, out, imageCleanup)
		if err != nil {
			return err
		}
	}

	return nil
}

func userImageUrl(r io.Reader) (string, error) {
	reader := bufio.NewReader(r)
	fmt.Print("Provide Image Url: ")

	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	input = strings.TrimSpace(input)
	return input, nil
}
