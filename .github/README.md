# Backdrop: Your Command-Line Wallpaper Manager

[![Issues Badge](https://img.shields.io/github/issues/JanMichaelSE/backdrop)](https://github.com/JanMichaelSE/backdrop/issues)
[![Pull Requests Badge](https://img.shields.io/github/issues-pr/JanMichaelSE/backdrop)](https://github.com/JanMichaelSE/backdrop/pulls)
[![Contributors Badge](https://img.shields.io/github/contributors/JanMichaelSE/backdrop)](https://github.com/JanMichaelSE/backdrop/graphs/contributors)
[![contributions welcome](https://img.shields.io/badge/contributions-welcome-brightgreen.svg?style=flat)](https://github.com/dwyl/esta/issues)

[![made-with-Go](https://img.shields.io/badge/Made%20with-Go-1f425f.svg)](https://go.dev/)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/JanMichaelSE/backdrop)

Welcome to **Backdrop**, a command-line utility designed to manage your desktop wallpapers with ease through the terminal. This tool allows you to set a new wallpaper, use fuzzy finding to select wallpaper, and specify the directory where your wallpaper images are stored.

## :star: Features

- Set a new wallpaper.
- Use `fzf` for fuzzy finding and selecting wallpapers.
- Specify a custom directory for wallpaper images.
- Create custom slideshows with a desired duration per slide.
- Set wallpapers using a URL to an image.

### Windows Features

- Full slideshow support on Windows via `.theme` file (multi-monitor aware).
- Auto-cleans cached wallpapers to prevent black screens.

---

## :wrench: Installation

### Install via Go

<details>
<summary>Steps to Install Using Go (v1.21 or Later)</summary>

```bash
go install github.com/janmichaelse/backdrop@latest
```

</details>

### Provided Binary

<details>
<summary>Download and Install the Binary</summary>

1. Download and unzip the `tar.gz` file for your operating system from the [Releases Page](https://github.com/JanMichaelSE/backdrop/releases).
2. Extract the archive:
   ```bash
   tar -xzf backdrop-gnome-desktop-v{current_version}.linux-amd64.tar.gz
   ```
3. Move the binary to the following path:
   ```bash
   mkdir -p $HOME/.backdrop/bin
   cp backdrop $HOME/.backdrop/bin
   ```
4. Add the binary to your `PATH`:
   ```bash
   echo 'export PATH=$HOME/.backdrop/bin:$PATH' >> $HOME/.bashrc
   source $HOME/.bashrc
   ```
   </details>

### Windows Installation

<details>
<summary>Steps to Add Binary to Windows PATH</summary>

1. Download the binary from the [Releases Page](https://github.com/JanMichaelSE/backdrop/releases).
2. Navigate to **Start Menu** and search for "Environment Variables".
3. Open "**Edit the System Environment Variables**".
4. Go to the "**Advanced**" tab and click "**Environment Variables**".
5. Edit the **Path** variable under System Variables and add the folder where the `backdrop.exe` binary is located. (Do not include the binary name.)

</details>

### Install From Source Code

<details>
<summary>Steps to Build and Install From Source</summary>

1. Clone the repository:
   ```bash
   git clone https://github.com/JanMichaelSE/backdrop.git
   cd backdrop
   ```
2. Build the project:
   ```bash
   go build
   ```
3. Move the binary to a directory in your `PATH`:
   ```bash
   mv backdrop $HOME/.backdrop/bin
   ```
4. Add the binary to your `PATH`:
   ```bash
   echo 'export PATH=$HOME/.backdrop/bin:$PATH' >> $HOME/.bashrc
   source $HOME/.bashrc
   ```
   </details>

---

## :wastebasket: Uninstall

### Uninstall Using Go

<details>
<summary>Steps to Uninstall</summary>

```bash
rm $(go env GOPATH)/bin/backdrop
go clean -modcache  # Optional
```

</details>

### Uninstall on Linux

<details>
<summary>Steps to Remove Installed Binary</summary>

1. Locate the binary:
   ```bash
   which backdrop
   ```
2. Delete the binary:
   ```bash
   rm /path/to/backdrop
   ```
3. (Optional) Remove configuration files, if applicable:
   ```bash
   rm -rf $HOME/.backdrop
   ```
4. Edit ~/.bashrc and remove the line that adds Backdrop to PATH
   ```bash
   nano ~/.bashrc
   ```
   Remove the following line (or equivalent):
   ```bash
   export PATH=$HOME/.backdrop/bin:$PATH
   ```

</details>
  
### Uninstall on Windows

<details>
<summary>Steps to Remove Installed Binary</summary>

1. Open a Command Prompt and verify the location of the binary:
   ```powershell
   where backdrop
   ```
2. Remove the folder path from your PATH environment variable: - Navigate to Start Menu and search for "Environment Variables". - Open Edit the System Environment Variables. - In the Advanced tab, click Environment Variables. - Under the System Variables, select Path and click Edit. - Remove the folder containing backdrop.exe from the list.
</details>

---

## &#x2705; What's Supported

### Operating Systems

- **Fully Supported**:
  - Ubuntu/GNOME-based Distros
  - Windows (Slideshow included via `.theme` file)
- **Limited Support**:
  - CentOS/MATE (Slideshows are not supported)
- **Not Supported**:
  - WSL (Windows Subsystem for Linux)
  - macOS (Make a PR because we won't do it)

---

## :computer: Usage

Backdrop provides several options for managing your wallpapers:

- `-p, --path <PATH>`:
  - Set a custom path for wallpaper images. Defaults to `$HOME/.backdrop/images` if not provided.
- `-h, --help`:
  - Displays help information on available commands.
- `-s, --slideshow`:
  - Configure and set a slideshow using images selected with `fzf`.
  - On Windows, selected images are copied to `%APPDATA%\BackdropSlideShow`, and a `.theme` file is generated and applied silently.
  - On Linux, a `backdrop_settings.xml` slideshow config is created in `~/.local/share/backgrounds`.
  - You will be prompted for the slideshow duration in minutes.
- `-u, --url`:
  - Download and set an image from a URL. Unaccepted images are deleted.
- `-v, --version`:
  - Print version information.

### Example: Set a Custom Wallpaper Directory

```bash
backdrop --path /path/to/your/wallpapers
```

---

## ⚠️ Known Issues

- On Windows, applying a slideshow theme may open the Settings window.

---

## :busts_in_silhouette: Authors

Backdrop is maintained by **JanMichaelSE** & **theweak1**.

---

## :handshake: Contributing

Contributions are always welcome! If you'd like to contribute, please follow these steps:

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

For any issues, please open an issue on GitHub.

---

## :email: Support

For commercial support or any other inquiries, please email the maintainer directly.

## :moneybag: Donations

If you find this project helpful and would like to support its development, please email the maintainer for information on how to make a donation.

Thank you for considering Backdrop for your wallpaper management needs! :grinning:
