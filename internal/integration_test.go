package internal

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/viper"
)

type testConfig struct {
	input              string
	config             *Config
	expOutput          string
	expError           error
	imageSelectionStub ImageSelection
}

func TestIntegration(t *testing.T) {
	initialWallpaper, err := getPreviousWallpaper()
	if err != nil {
		t.Fatalf("Error getting system initial wallpaper for eventual cleanup after tests: %v", err)
	}
	defer setWallpaper(initialWallpaper)

	testFile, cleanupTempImageFile := setupTempImageFile(t)
	defer cleanupTempImageFile()

	cleanupCustomPath := cleanupCustomPathTest(t)
	defer cleanupCustomPath()

	t.Run("SuccessSetWallpaper", func(t *testing.T) {
		testCase := testConfig{
			config:             NewConfig("", "", false, false),
			input:              "y\n",
			expOutput:          "Successfully changed background image!",
			expError:           nil,
			imageSelectionStub: func(s []string) (string, error) { return testFile, nil },
		}

		var out bytes.Buffer
		imageSelection = testCase.imageSelectionStub
		input = strings.NewReader(testCase.input)

		err := BackdropAction(&out, testCase.config, []string{})
		if err != nil {
			t.Errorf("Expected NO error, but got '%v' instead", err)
		}

		if !strings.Contains(out.String(), testCase.expOutput) {
			t.Errorf("Expected output '%v', but got '%v' instead", testCase.expOutput, out.String())
		}
	})

	t.Run("SuccessConfigurePath", func(t *testing.T) {
		testCase := testConfig{
			config:             NewConfig("../test/testData", "", false, false),
			input:              "y\n",
			expOutput:          "Successfully changed background image!",
			expError:           nil,
			imageSelectionStub: func(s []string) (string, error) { return "testImage.png", nil },
		}

		var out bytes.Buffer
		imageSelection = testCase.imageSelectionStub
		input = strings.NewReader(testCase.input)

		err := BackdropAction(&out, testCase.config, []string{})
		if err != nil {
			t.Errorf("Expected NO error, but got '%v' instead", err)
		}

		imagePath, err := getUserImagesPath()
		if err != nil {
			t.Errorf("Expected NO error, but got '%v' instead", err)
		}

		fmt.Println("IMAGE PATH DURING TEST:", imagePath)
		if !strings.Contains(imagePath, testCase.config.path) {
			t.Errorf("Expected image path '%v', got image path '%v'", testCase.config.path, imagePath)
		}

		if !strings.Contains(out.String(), testCase.expOutput) {
			t.Errorf("Expected output '%v', but got '%v' instead", testCase.expOutput, out.String())
		}
	})
}

func setupTempImageFile(t *testing.T) (string, func()) {
	t.Helper()

	wallpapersPath, err := getUserImagesPath()
	if err != nil {
		t.Fatalf("Error getting userpath in test setup: %v", err)
	}

	file, err := os.CreateTemp(wallpapersPath, "backdropTestFile")
	if err != nil {
		t.Fatalf("Error creating temp file for test setup: %v", err)
	}

	return file.Name(), func() {
		os.Remove(file.Name())
	}
}

func cleanupCustomPathTest(t *testing.T) func() {
	t.Helper()

	originalCustomImagePath, ok := viper.Get("WallpapersPath").(string)

	return func() {
		if ok {
			err := configureImagePath(originalCustomImagePath)
			if err != nil {
				t.Fatalf("Error during CustomPathTest cleanup, couldn't set original configureImagePath: '%v'", err)
			}
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
