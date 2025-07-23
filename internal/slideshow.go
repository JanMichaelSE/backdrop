package internal

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/janmichaelse/backdrop/internal/os_Specifics"
)

var (
	inputDuration io.Reader = os.Stdin
)

func handleSlideshow(out io.Writer, wallpapersPath string, wallpapers []string, imageSelection FuzzySelection) error {
	for hasConfirmed := false; !hasConfirmed; {
		if err := processSlideshow(wallpapersPath, wallpapers, imageSelection); err != nil {
			return err
		}

		previousWallpaper, err := getPreviousWallpaper()
		if err != nil {
			return err
		}

		hasConfirmed, err = handleSelectionConfirmation(previousWallpaper, out, &SelectionOptions{
			Prompt:         "Save slideshow configuration? [y/N]: ",
			SuccessMessage: "Slideshow has been set successfully.",
			Cleanup:        func() {},
		})

		if err != nil {
			return err
		}
	}
	return nil
}

func processSlideshow(wallpapersPath string, wallpapers []string, imageSelection FuzzySelection) error {
	selectedWallpaper, err := imageSelection(wallpapers)
	if err != nil {
		return err
	}

	duration, err := getDurationFromUser(inputDuration)
	if err != nil {
		return err
	}

	configuredWallpaper, err := configureSlideShow(selectedWallpaper, wallpapersPath, duration)
	if err != nil {
		return err
	}

	if runtime.GOOS != "windows" {
		if err := setWallpaper(configuredWallpaper); err != nil {
			return err
		}
	}
	return nil
}

func getDurationFromUser(r io.Reader) (int, error) {
	fmt.Print("What should be the duration per slide? (In minutes): ")
	input, err := bufio.NewReader(r).ReadString('\n')
	if err != nil {
		return 0, fmt.Errorf("failed to read duration input: %w", err)
	}

	durationStr := strings.TrimSpace(input)
	duration, err := strconv.Atoi(durationStr)
	if err != nil || duration <= 0 {
		return 0, fmt.Errorf("invalid duration: please enter a positive integer")
	}

	return duration * 60_000, nil
}

func configureSlideShow(imageText, wallpapersPath string, duration int) (string, error) {
	images := strings.Split(imageText, ";")
	switch runtime.GOOS {
	case "linux":
		return os_Specifics.ConfigureSlideShowLinux(images, wallpapersPath, duration)
	case "windows":
		return os_Specifics.ConfigureSlideShowWindows(images, wallpapersPath, duration)
	default:
		return "", ErrNoCompatibleOS
	}
}
