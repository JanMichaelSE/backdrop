package internal

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

var (
	ErrNoValidImagesPath              = errors.New("User does not have valid images path configured.")
	ErrNoCompatibleDesktopEnvironment = errors.New("No compatible desktop environment found.")
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

	fmt.Println("Wallpaper path:", wallpapersPath)

	wallpapers := getWallpapers(wallpapersPath)
	fmt.Println("Wallpapers:", wallpapers)

outter:
	for {
		previousWallpaper, err := getPreviousWallpaper()
		if err != nil {
			return err
		}
		fmt.Println("Previous wallpaper:", previousWallpaper)

		var selectedWallpaper string
		if config.isFuzzy {
			fmt.Println("Selected search")
			selectedWallpaper, err = useFZF(wallpapers, out)
			if err != nil {
				return err
			}
		}

		fullSelectedPath := filepath.Join(wallpapersPath, selectedWallpaper)
		if stats, err := os.Stat(fullSelectedPath); err == nil && stats.Mode().IsRegular() {
			setWallpaper(fullSelectedPath)
		}

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
				break
			default:
				fmt.Println("Invalid input...")
			}
		}

	}

	// TODO: Still not sure what to return here.
	return nil
}

// TODO: WHERE I LEFT OFF
// - Currently FZF implementation is more complex than I expected,
// but this could be an opportunity to get rid of this dependency and provide
// an alternative way to select an image that is easy to use and maintain.
// - Need to research above options.
func useFZF(fileNames []string, out io.Writer) (string, error) {
	cmd := exec.Command("fzf", "--layout=reverse")
	cmd.Stdin = strings.NewReader(strings.Join(fileNames, "\n"))
	cmd.Stdout = out
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return "", err
	}

	return "", nil
}

func getUserImagesPath() (string, error) {
	homePath, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	wallpapersPath, ok := viper.Get("WallpapersPath").(string)
	if ok {
		if stat, err := os.Stat(wallpapersPath); err == nil && stat.IsDir() {
			return wallpapersPath, nil
		}
	}

	configPath := fmt.Sprintf("%v/.config/backdrop/wallpapers", homePath)
	if stat, err := os.Stat(configPath); err == nil && stat.IsDir() {
		return configPath, nil
	}

	picturesPath := fmt.Sprintf("%v/Pictures/wallpapers", homePath)
	if stat, err := os.Stat(picturesPath); err == nil && stat.IsDir() {
		return picturesPath, nil
	}

	// TODO: Might want to see how to make this become a good message for the user
	return "", ErrNoValidImagesPath
}

func getWallpapers(path string) []string {
	fileEntries, _ := os.ReadDir(path)

	files := make([]string, 0, len(fileEntries))
	for _, file := range fileEntries {
		files = append(files, file.Name())
	}

	return files
}

func getPreviousWallpaper() (string, error) {
	cmdSettings := exec.Command("gsettings", "list-schemas")
	var outSettings bytes.Buffer
	cmdSettings.Stdout = &outSettings
	if err := cmdSettings.Run(); err != nil {
		return "", err
	}

	if strings.Contains(outSettings.String(), "gnome.desktop.background") {
		cmdGetGnome := exec.Command("gsettings", "get", "org.gnome.desktop.background", "picture-uri")
		var outGetGnome bytes.Buffer
		cmdGetGnome.Stdout = &outGetGnome
		err := cmdGetGnome.Run()
		if err != nil {
			return "", err
		}

		uri := strings.Trim(outGetGnome.String(), "\n")
		if strings.Contains(uri, "://") {
			parts := strings.SplitN(uri, "://", 2)
			if len(parts) == 2 {
				return parts[1], nil
			}
		}

		return uri, nil
	}

	if strings.Contains(outSettings.String(), "mate.desktop.background") {
		cmdGetMate := exec.Command("gsettings", "get", "org.mate.desktop.background", "picture-uri")
		var outGetMate bytes.Buffer
		cmdGetMate.Stdout = &outGetMate
		err := cmdGetMate.Run()
		if err != nil {
			return "", err
		}

		uri := strings.Trim(outGetMate.String(), "\n")
		if strings.Contains(uri, "://") {
			parts := strings.SplitN(uri, "://", 2)
			if len(parts) == 2 {
				return parts[1], nil
			}
		}

		return uri, nil
	}

	return "", ErrNoCompatibleDesktopEnvironment
}

func setWallpaper(wallpaper string) error {
	cmdSettings := exec.Command("gsettings", "list-schemas")
	var outSettings bytes.Buffer
	cmdSettings.Stdout = &outSettings
	if err := cmdSettings.Run(); err != nil {
		return err
	}

	if strings.Contains(outSettings.String(), "gnome.desktop.background") {
		cmdGetGnome := exec.Command("gsettings", "set", "org.gnome.desktop.background", "picture-uri", wallpaper)
		err := cmdGetGnome.Run()
		if err != nil {
			// TODO: NEED TO BETTER HANDLE THIS ERROR
			return err
		}
	}

	if strings.Contains(outSettings.String(), "mate.desktop.background") {
		cmdGetGnome := exec.Command("gsettings", "set", "org.mate.desktop.background", "picture-uri", wallpaper)
		err := cmdGetGnome.Run()
		if err != nil {
			// TODO: NEED TO BETTER HANDLE THIS ERROR
			return err
		}
	}

	// TODO: Might want to see how to make this become a good message for the user
	return ErrNoCompatibleDesktopEnvironment
}
