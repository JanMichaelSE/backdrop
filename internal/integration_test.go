package internal

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/pkg/errors"
)

func TestIntegration(t *testing.T) {
	initialWallpaper, err := getPreviousWallpaper()
	if err != nil {
		t.Fatalf("Error getting system initial wallpaper for eventual cleanup after tests: %v", err)
	}
	defer setWallpaper(initialWallpaper)

	testFile, cleanup := setupTempImage(t)
	defer cleanup()

	testCases := []struct {
		name               string
		input              string
		config             *Config
		expOutput          string
		expError           error
		imageSelectionStub ImageSelection
	}{
		{
			name:               "SuccessSetWallPaper",
			config:             NewConfig("", "", false, false),
			input:              "y\n",
			expOutput:          "Successfully changed background image!",
			expError:           nil,
			imageSelectionStub: func(s []string) (string, error) { return testFile, nil },
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var out bytes.Buffer
			var err error

			imageSelection = tc.imageSelectionStub
			input = strings.NewReader(tc.input)
			err = BackdropAction(&out, tc.config, []string{})
			if tc.expError != nil {
				if err == nil {
					t.Error("Expected an error, but got 'nil' instead")
				}

				if !errors.Is(err, tc.expError) {
					t.Errorf("Expected error '%v', but got '%v' instead", tc.expError, err)
				}
			}

			if err != nil {
				t.Errorf("Expected NO error, but got '%v' instead", err)
			}

			if !strings.Contains(out.String(), tc.expOutput) {
				t.Errorf("Expected output '%v', but got '%v' instead", tc.expOutput, out.String())
			}

		})
	}
}

func setupTempImage(t *testing.T) (string, func()) {
	t.Helper()

	wallpapersPath, err := getUserImagesPath()
	if err != nil {
		t.Fatalf("Error getting userpath in setup function: %v", err)
	}

	file, err := os.CreateTemp(wallpapersPath, "backdropTestFile")
	if err != nil {
		t.Fatalf("Error creating temp file for setup function: %v", err)
	}

	return file.Name(), func() {
		os.Remove(file.Name())
	}
}
