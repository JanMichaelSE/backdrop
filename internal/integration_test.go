package internal

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/spf13/viper"
)

type testConfig struct {
	inputConfirmation string
	inputDuration     string
	inputImageUrl     string
	config            *Config
	expOutput         string
	expError          error
	imageSelectorStub GetFuzzySelector
}

// NOTE: If the tests for TestSetImageUrl fail for fetching image. Try changing these variables.
// I wanted a real network request as part of my tests.
var (
	TEST_IMAGE_URL      = "https://static.frontendmasters.com/assets/teachers/theprimeagen/thumb.webp"
	TEST_IMAGE_URL_NAME = "thumb.webp"
)

func TestSetWallpaper(t *testing.T) {
	initialWallpaper, err := getPreviousWallpaper()
	if err != nil {
		t.Fatalf("Error getting system initial wallpaper for eventual cleanup after tests: %v", err)
	}
	defer setWallpaper(initialWallpaper)

	testFile, cleanupTempImageFile := setupTempImageFile(t)
	defer cleanupTempImageFile()

	cleanupCustomPath := cleanupCustomPath(t)
	defer cleanupCustomPath()

	testCase := testConfig{
		config:            NewConfig("", false, false),
		inputConfirmation: "y\n",
		expOutput:         "Successfully changed background image!",
		expError:          nil,
		imageSelectorStub: func(c *Config) FuzzySelection {
			return func(s []string) (string, error) {
				return testFile, nil
			}
		},
	}

	var out bytes.Buffer
	getSelector = testCase.imageSelectorStub
	inputConfirmation = strings.NewReader(testCase.inputConfirmation)
	if err := BackdropAction(&out, testCase.config, []string{}); err != nil {
		t.Errorf("Expected NO error, but got '%v' instead", err)
	}

	if !strings.Contains(out.String(), testCase.expOutput) {
		t.Errorf("Expected output '%v', but got '%v' instead", testCase.expOutput, out.String())
	}
}

func TestConfigurePath(t *testing.T) {
	initialWallpaper, err := getPreviousWallpaper()
	if err != nil {
		t.Fatalf("Error getting system initial wallpaper for eventual cleanup after tests: %v", err)
	}
	defer setWallpaper(initialWallpaper)

	cleanupCustomPath := cleanupCustomPath(t)
	defer cleanupCustomPath()

	testCase := testConfig{
		config:            NewConfig("../test/testData/images", false, false),
		inputConfirmation: "y\n",
		inputDuration:     "10\n",
		expOutput:         "Successfully changed background image!",
		expError:          nil,
		imageSelectorStub: func(c *Config) FuzzySelection {
			return func(s []string) (string, error) {
				return "testImage.jpg", nil
			}
		},
	}

	var out bytes.Buffer
	getSelector = testCase.imageSelectorStub
	inputConfirmation = strings.NewReader(testCase.inputConfirmation)
	inputDuration = strings.NewReader(testCase.inputDuration)
	if err := BackdropAction(&out, testCase.config, []string{}); err != nil {
		t.Errorf("Expected NO error, but got '%v' instead", err)
	}

	wallpapersPath, err := getUserWallpapersPath()
	if err != nil {
		t.Errorf("Expected NO error, but got '%v' instead", err)
	}

	if !strings.Contains(wallpapersPath, testCase.config.path) {
		t.Errorf("Expected image path '%v', got image path '%v'", testCase.config.path, wallpapersPath)
	}

	if !strings.Contains(out.String(), testCase.expOutput) {
		t.Errorf("Expected output '%v', but got '%v' instead", testCase.expOutput, out.String())
	}
}

func TestSetSlideShow(t *testing.T) {

	if runtime.GOOS == "windows" {
		t.Skip("SlideShow is not supported on Windows")
	}
	initialWallpaper, err := getPreviousWallpaper()
	if err != nil {
		t.Fatalf("Error getting system initial wallpaper for eventual cleanup after tests: %v", err)
	}
	defer setWallpaper(initialWallpaper)

	cleanupCustomPath := cleanupCustomPath(t)
	defer cleanupCustomPath()

	testCase := testConfig{
		config:            NewConfig("../test/testData/images", false, true),
		inputConfirmation: "y\n",
		expOutput:         "Successfully changed background image!",
		expError:          nil,
		imageSelectorStub: func(c *Config) FuzzySelection {
			return func(s []string) (string, error) {
				return "testImage.png;testImage2.png;testImage3.png", nil
			}
		},
	}

	var out bytes.Buffer
	getSelector = testCase.imageSelectorStub
	inputConfirmation = strings.NewReader(testCase.inputConfirmation)
	if err := BackdropAction(&out, testCase.config, []string{}); err != nil {
		t.Errorf("Expected NO error, but got '%v' instead", err)
	}

	wallpapersPath, err := getUserWallpapersPath()
	if err != nil {
		t.Errorf("Expected NO error, but got '%v' instead", err)
	}

	if !strings.Contains(wallpapersPath, testCase.config.path) {
		t.Errorf("Expected image path '%v', got image path '%v'", testCase.config.path, wallpapersPath)
	}

	currentWallpaper, err := getPreviousWallpaper()
	if err != nil {
		t.Fatalf("Unexpected error getting wallpaper: %v", err)
	}
	if !strings.Contains(currentWallpaper, "backdrop_settings.xml") {
		t.Errorf("Expected wallpaper '%v' to be set for slideshow, but got '%v' instead", "backdrop_settings.xml", currentWallpaper)
	}

	currentSlideShowSettings := getCurrentBackdropSlideShowSettings(t)
	expectedSlideShowSettings := getExpectedBackdropSlideShowSettings(t)
	if currentSlideShowSettings != expectedSlideShowSettings {
		t.Errorf("Expected slide show settings:\n %v \nGot slide show settings:\n %v", expectedSlideShowSettings, currentSlideShowSettings)
	}

	if !strings.Contains(out.String(), testCase.expOutput) {
		t.Errorf("Expected output '%v', but got '%v' instead", testCase.expOutput, out.String())
	}
}

func TestSetImageUrl(t *testing.T) {
	initialWallpaper, err := getPreviousWallpaper()
	if err != nil {
		t.Fatalf("Error getting system initial wallpaper for eventual cleanup after tests: %v", err)
	}
	defer setWallpaper(initialWallpaper)

	cleanup := cleanupImageUrl(t)
	defer cleanup()

	testCase := testConfig{
		config:            NewConfig("", true, false),
		inputConfirmation: "y\n",
		inputImageUrl:     fmt.Sprintf("%v\n", TEST_IMAGE_URL),
		expOutput:         "Successfully changed background image!",
		expError:          nil,
	}

	var out bytes.Buffer
	inputConfirmation = strings.NewReader(testCase.inputConfirmation)
	inputImageUrl = strings.NewReader(testCase.inputImageUrl)
	if err := BackdropAction(&out, testCase.config, []string{}); err != nil {
		t.Errorf("Expected NO error, but got '%v' instead", err)
	}

	currentWallpaper, err := getPreviousWallpaper()
	if err != nil {
		t.Fatalf("Unexpected error getting wallpaper: %v", err)
	}
	if !strings.Contains(currentWallpaper, TEST_IMAGE_URL_NAME) {
		t.Errorf("Expected wallpaper '%v' to be set for image url, but got '%v' instead", TEST_IMAGE_URL_NAME, currentWallpaper)
	}

	if !strings.Contains(out.String(), testCase.expOutput) {
		t.Errorf("Expected output '%v', but got '%v' instead", testCase.expOutput, out.String())
	}
}

func getCurrentBackdropSlideShowSettings(t *testing.T) string {
	t.Helper()

	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatal("Got error when getting user home dir in test helper slide show settings.")
	}

	slideShowSettings := filepath.Join(homeDir, ".local", "share", "backgrounds", "backdrop_settings", "backdrop_settings.xml")
	fileBytes, err := os.ReadFile(slideShowSettings)
	if err != nil {
		t.Fatalf("Got error when opening file in test helper slide show settings: %v", err)
	}

	return string(fileBytes)
}

func getExpectedBackdropSlideShowSettings(t *testing.T) string {
	t.Helper()

	slideShowSettings := "../test/testData/backdrop_settings.xml"
	fileBytes, err := os.ReadFile(slideShowSettings)
	if err != nil {
		t.Fatalf("Got error when opening file in test helper slide show settings: %v", err)
	}

	return string(fileBytes)
}

func setupTempImageFile(t *testing.T) (string, func()) {
	t.Helper()

	wallpapersPath, err := getUserWallpapersPath()
	if err != nil {
		t.Fatalf("Error getting userpath in test setup: %v", err)
	}

	file, err := os.CreateTemp(wallpapersPath, "backdropTestFile*.jpg")
	if err != nil {
		t.Fatalf("Error creating temp file for test setup: %v", err)
	}
	file.Close()
	_, fileName := filepath.Split(file.Name())

	return fileName, func() {
		if err := os.Remove(file.Name()); err != nil && !os.IsNotExist(err) {
			t.Fatalf("Error cleaning up temp file %s: %v", fileName, err)
		}
	}
}

func cleanupCustomPath(t *testing.T) func() {
	t.Helper()

	originalCustomImagePath, ok := viper.Get("WallpapersPath").(string)

	return func() {
		if ok {
			err := configureWallpaperPath(originalCustomImagePath)
			if err != nil {
				t.Fatalf("Error during CustomPathTest cleanup, couldn't set original configureImagePath: '%v'", err)
			} else {
				homePath, err := os.UserHomeDir()
				if err != nil {
					t.Fatalf("Error during CustomPathTest cleanup, could'nt get user home directory: '%v'", err)
				}

				configPath := filepath.Join(homePath, ".backdrop.yaml")
				os.Remove(configPath)

			}
		}
	}
}

func cleanupImageUrl(t *testing.T) func() {
	t.Helper()

	wallpapersPath, err := getUserWallpapersPath()
	if err != nil {
		t.Fatalf("Error getting userpath in test setup: %v", err)
	}

	filePath := filepath.Join(wallpapersPath, TEST_IMAGE_URL_NAME)
	return func() {
		os.Remove(filePath)
	}
}
