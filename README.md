# Backdrop: Your Command-Line Wallpaper Manager

[![Issues Badge](https://img.shields.io/github/issues/JanMichaelSE/backdrop)](https://github.com/JanMichaelSE/backdrop/issues)
[![Pull Requests Badge](https://img.shields.io/github/issues-pr/JanMichaelSE/backdrop)](https://github.com/JanMichaelSE/backdrop/pulls)
[![Contributors Badge](https://img.shields.io/github/contributors/JanMichaelSE/backdrop)](https://github.com/JanMichaelSE/backdrop/graphs/contributors)
[![contributions welcome](https://img.shields.io/badge/contributions-welcome-brightgreen.svg?style=flat)](https://github.com/dwyl/esta/issues)

Welcome to **Backdrop**, a command-line utility designed to manage your desktop wallpapers with ease through the terminal. This tool allows you to set a new wallpaper, use fuzzy finding to select wallpaper, and specify the directory where your wallpaper images are stored.

## :star: Features

- Set a new wallpaper
- Specify a custom path for wallpaper images
- Set custom slideshow of images with your desired duration per slide.
- Fuzzy finding of wallpapers (Requires fzf)
- Uninstall Backdrop

## :wrench: Installation

The installation process is straightforward. Simply clone the repository and run the `install.sh` script.

```bash
git clone https://github.com/JanMichaelSE/backdrop.git
cd backdrop
./install.sh
```

## &#x2705; What's Supported

#### Operating Systems
- Ubuntu/GNOME based Distros
- CentOS/MATE (Doesn't support Slideshows)
- Coming Soon:
    - MacOS
#### Shells
- Bash
- Zsh
- Fish


## :package: Dependencies

Backdrop requires `fzf` for fuzzy finding of wallpapers. If `fzf` is not installed, the script will automatically install it for you.

## :computer: Usage

Backdrop provides several options for managing your wallpapers:

- `-p, --path <PATH>`: Set a custom path to find wallpaper images. If not provided, a default path will be used.
- `-f, --fuzzy`: Performs a fuzzy finding (Requires fzf).
- `-h, --help`: Displays help information on how to use the command, listing all available options.
- `-s, --slideshow`: Will configure and set a custom slideshow of images you select with fzf. To select multiple images hit "Tab" on the images you desire to select, then hit "Enter" to confirm.
- `--uninstall`: Will uninstall Backdrop by removing all PATHs and Backdrop files.

For example, to set a custom path for your wallpapers, you can use the `-p` or `--path` flag:

```bash
backdrop --path /path/to/your/wallpapers
```

## :busts_in_silhouette: Authors

Backdrop is currently maintained by JanMichaelSE.

## :handshake: Contributing

Contributions are always welcome! If you'd like to contribute, please follow these steps:

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

For any issues, please open an issue on GitHub.

## :email: Support

For commercial support or any other inquiries, please email the maintainer directly.

## :moneybag: Donations

If you find this project helpful and would like to support its development, please email the maintainer for information on how to make a donation.

Thank you for considering Backdrop for your wallpaper management needs! :grinning:
