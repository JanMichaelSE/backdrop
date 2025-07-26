package os_Specifics

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func GetWindowsConfigFilePath() (string, error) {
	appData := os.Getenv("APPDATA")
	if appData == "" {
		return "", fmt.Errorf("APPDATA is not set")
	}
	dir := filepath.Join(appData, "Backdrop")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	return filepath.Join(dir, "config.yaml"), nil
}

func GetWindowsConfigPath() string {
	appData := os.Getenv("APPDATA") // e.g., C:\Users\YourName\AppData\Roaming
	return filepath.Join(appData, "Backdrop")
}

func ConfigureSlideShowWindows(images []string, wallpapersPath string, duration int) (string, error) {
	slideShowDir := filepath.Join(os.Getenv("APPDATA"), "BackdropSlideShow")

	if err := os.MkdirAll(slideShowDir, 0777); err != nil {
		return "", fmt.Errorf("failed to create slideshow directory: %w", err)
	}

	entries, err := os.ReadDir(slideShowDir)
	if err != nil {
		return "", fmt.Errorf("failt to read content of slideshow directy: %w", err)
	}

	for _, entry := range entries {
		entryPath := filepath.Join(slideShowDir, entry.Name())
		if err := os.RemoveAll(entryPath); err != nil {
			return "", fmt.Errorf("failed to remove existing file %s: %w", entryPath, err)
		}
	}

	for _, img := range images {
		src := filepath.Join(wallpapersPath, img)
		dst := filepath.Join(slideShowDir, filepath.Base(img))

		data, _ := os.ReadFile(src)

		if err := os.WriteFile(dst, data, 0666); err != nil {
			return "", fmt.Errorf("failed to write image to slideshow folder: %w", err)
		}
	}

	if err := setWindowsSlideShow(slideShowDir, duration); err != nil {
		return "", err
	}

	if len(images) == 0 {
		return "", fmt.Errorf("no images provided")
	}

	firstImagePath := filepath.Join(slideShowDir, filepath.Base(images[0]))
	themeFile, err := createWindowsThemeFile(firstImagePath, slideShowDir, duration)
	if err != nil {
		return "", fmt.Errorf("failed to create theme file: %w", err)
	}

	err = applyWindowsTheme(themeFile)
	if err != nil {
		return "", err
	}

	return slideShowDir, nil
}

func setWindowsSlideShow(folder string, duration int) error {
	cmd := exec.Command("powershell", "-Command", fmt.Sprintf(`
	$RegPath = "HKCU:\Control Panel\Personalization\Desktop Slideshow"
	Set-ItemProperty -Path $RegPath -Name Interval -Value %d
	Set-ItemProperty -Path $RegPath -Name Shuffle -Value 1
	Set-ItemProperty -Path $RegPath -Name SlideshowEnabled -Value 1
	
	$ThemePath = "HKCU:\Software\Microsoft\Windows\CurrentVersion\Themes"
	Set-ItemProperty -Path $ThemePath -Name SlideshowDirectory -Value "%s"
	
	$WallpaperPath = "HKCU:\Control Panel\Desktop"
	Set-ItemProperty -Path $WallpaperPath -Name Wallpaper -Value ""
	Set-ItemProperty -Path $WallpaperPath -Name WallpaperStyle -Value 10
	Set-ItemProperty -Path $WallpaperPath -Name TileWallpaper -Value 0
	
	# Critical: Remove corrupted or cached wallpaper to avoid black background
	$transcoded = "$env:APPDATA\Microsoft\Windows\Themes\TranscodedWallpaper"
	if (Test-Path $transcoded) { Remove-Item $transcoded -Force -ErrorAction SilentlyContinue }
	
	RUNDLL32.EXE user32.dll, UpdatePerUserSystemParameters
	`, duration, folder))

	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to set registry slideshow values: %v\n%s", err, output)
	}

	return nil
}

func createWindowsThemeFile(firstImagePath, slideshowDir string, duration int) (string, error) {
	themesDir := filepath.Join(os.Getenv("LOCALAPPDATA"), "Microsoft", "Windows", "Themes")
	if err := os.MkdirAll(themesDir, 0777); err != nil {
		return "", fmt.Errorf("failed to create themes directory: %w", err)
	}

	themeFilePath := filepath.Join(themesDir, "backdrop.theme")

	content := fmt.Sprintf(`[Theme]
	DisplayName=Backdrop Slideshow

	[Control Panel\Desktop]
	wallpaper=%s
	TileWallpaper=0
	WallpaperStyle=10
	PicturePosition=10
	SlideshowEnabled=1
	MultimonBackgrounds=1

	[Slideshow]
	ImagesRootPath=%s
	Interval=%d
	Shuffle=1

	[VisualStyles]
	Path=%%SystemRoot%%\resources\Themes\Aero\Aero.msstyles
	ColorStyle=NormalColor
	Size=NormalSize
	AutoColorization=0
	VisualStyleVersion=10

	[MasterThemeSelector]
	MTSM=RJSPBS

	[Sounds]
	SchemeName=@mmres.dll,-800
	`, firstImagePath, slideshowDir, duration)

	if err := os.WriteFile(themeFilePath, []byte(content), 0666); err != nil {
		return "", fmt.Errorf("failed to write theme file: %w", err)
	}

	return themeFilePath, nil
}

func applyWindowsTheme(themePath string) error {

	if _, err := os.Stat(themePath); err != nil {
		return fmt.Errorf("theme file does not exist: %w", err)
	}

	if themePath == "" {
		return fmt.Errorf("themePath is empty")
	}
	psScript := fmt.Sprintf(`
	# Set theme path
	Set-ItemProperty -Path "HKCU:\\Software\\Microsoft\\Windows\\CurrentVersion\\Themes" -Name "CurrentTheme" -Value "%s"

	# Set slideshow options
	Set-ItemProperty -Path "HKCU:\\Control Panel\\Desktop" -Name "SlideshowEnabled" -Value 1
	Set-ItemProperty -Path "HKCU:\\Control Panel\\Personalization\\Desktop Slideshow" -Name "SlideshowEnabled" -Value 1
	Set-ItemProperty -Path "HKCU:\\Control Panel\\Desktop\\PerMonitorSettings" -Name "SlideshowEnabled" -Value 1 -ErrorAction SilentlyContinue

	# Clear TranscodedWallpaper (sometimes holds stale data)
	$transcoded = "$env:APPDATA\\Microsoft\\Windows\\Themes\\TranscodedWallpaper"
	if (Test-Path $transcoded) { Remove-Item $transcoded -Force -ErrorAction SilentlyContinue }

	# Refresh system settings
	RUNDLL32.EXE user32.dll, UpdatePerUserSystemParameters

	# 🧠 Apply the .theme file again to force slideshow to start (silent, but effective)
	Start-Process -FilePath "%s" -WindowStyle Hidden

	Start-Sleep -Milliseconds 500
	`, themePath, themePath)

	cmd := exec.Command("powershell", "-NoProfile", "-ExecutionPolicy", "Bypass", "-Command", psScript)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to apply theme: %w\nOutput: %s", err, string(output))
	}

	return nil
}
