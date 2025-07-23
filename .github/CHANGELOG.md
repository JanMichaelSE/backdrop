# Changelog

## Version 2.1.0

### ✨ Added

- OS-specific separation of logic using `runtime.GOOS` switch.
- Introduced platform-based implementations:
  - `configureWallpaperPath` now routes to `GetLinuxConfigFilePath` or `GetWindowsConfigFilePath`.
  - Windows-specific slideshow functions (e.g., `ConfigureSlideShowWindows`, `setWindowsSlideshow`, etc.)
  - Linux-specific slideshow and wallpaper logic.
- New helper utilities for:
  - Slide show creation and configuration.
  - GSettings wallpaper retrieval (Linux).
  - Theme file creation and application (Windows).
- Refined CLI prompt with customizable confirmation and success messages.

### 🔧 Changed

- Version updated from `2.0.0` to `2.1.0`.
- Modularized slideshow handling into `windows.go` and `linux.go`.
- Refactored `handleSelectionConfirmation` to use a `SelectionOptions` struct for optional parameters, improving flexibility and readability.
- Updated all usages (including in `handleImageUrl`) to use the new struct-based signature.
- Test cases now include Windows-specific slideshow validation (theme and folder checks).

### 🧹 Removed

- Direct logic previously embedded in cross-platform functions; moved to appropriate OS-specific files.

### 📁 Project Structure

- Significant folder and file modularization by platform:
  - `internal/os_Specifics/windows.go`
  - `internal/os_Specifics/linux.go`

### 📦 Dependency/Tooling

- Retained use of `viper` for configuration management, centralizing path logic.

---

## Known Issues

- 🐛 **Windows Slideshow Theme Bug**: When applying the slideshow theme on Windows, the system settings window briefly opens and must be closed manually by the user.
