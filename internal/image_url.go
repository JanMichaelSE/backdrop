package internal

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
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

		hasConfirmed, err = handleSelectionConfirmation(previousWallpaper, "", "", out, imageCleanup)
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

func downloadImage(imageUrl, wallpapersPath string) (string, error) {
	resp, err := http.Get(imageUrl)
	if err != nil {
		return "", fmt.Errorf("Could not fetch image, got error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Got bad response code '%v' from request. Cannot proceed.", resp.StatusCode)
	}

	fileName := strings.Split(imageUrl, "/")
	sanitizedFileName := sanitizeFilename(fileName[len(fileName)-1])
	filePath := filepath.Join(wallpapersPath, sanitizedFileName)
	file, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("Could not create file for image url, got error: %v", err)
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return "", fmt.Errorf("Could not copy contents from image url to file, got error: %v", err)
	}

	return filePath, nil
}

func sanitizeFilename(filename string) string {
	// URL Decode
	decodedFilename, err := url.QueryUnescape(filename)
	if err != nil {
		decodedFilename = filename // Use the original filename if decoding fails
	}

	// Define a list of invalid characters and their replacements
	replacements := map[string]string{
		"/":  "_",
		"\\": "_",
		"<":  "_",
		">":  "_",
		":":  "_",
		"\"": "_",
		"|":  "_",
		"?":  "_",
		"*":  "_",
	}

	// Replace invalid characters
	for oldChar, newChar := range replacements {
		decodedFilename = strings.ReplaceAll(decodedFilename, oldChar, newChar)
	}

	// Normalize spaces (optional)
	decodedFilename = strings.Join(strings.Fields(decodedFilename), " ")

	return decodedFilename
}
