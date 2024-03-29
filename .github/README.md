# Backdrop: Your Command-Line Wallpaper Manager

[![Issues Badge](https://img.shields.io/github/issues/JanMichaelSE/backdrop)](https://github.com/JanMichaelSE/backdrop/issues)
[![Pull Requests Badge](https://img.shields.io/github/issues-pr/JanMichaelSE/backdrop)](https://github.com/JanMichaelSE/backdrop/pulls)
[![Contributors Badge](https://img.shields.io/github/contributors/JanMichaelSE/backdrop)](https://github.com/JanMichaelSE/backdrop/graphs/contributors)
[![contributions welcome](https://img.shields.io/badge/contributions-welcome-brightgreen.svg?style=flat)](https://github.com/dwyl/esta/issues)

[![made-with-Go](https://img.shields.io/badge/Made%20with-Go-1f425f.svg)](https://go.dev/)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/JanMichaelSE/backdrop)

Welcome to **Backdrop**, a command-line utility designed to manage your desktop wallpapers with ease through the terminal. This tool allows you to set a new wallpaper, use fuzzy finding to select wallpaper, and specify the directory where your wallpaper images are stored.

## :star: Features

- Set a new wallpaper
- Set wallpapers with `fzf`
- Specify a custom path for wallpaper images
- Set custom slideshow of images with your desired duration per slide.
- Set wallpapers by providing a URL to an image.

## :wrench: Installation

### If you have go version 1.21 or greater
You can install backdrop by using the following go command
```bash
go install github.com/JanMichaelSE/backdrop@latest
```

### Provided binary
you can install Backdrop by downloading and unziping the `tar.gz` file for your operating system (if available) that is included in the release
```bash
tar -xzf backdrop-gnome-desktop-v1.0.0.linux-amd64.tar.gz
```

After place it under the following path:
```bash
mkdir -p $HOME/.backdrop/bin
cp backdrop $HOME/.backdrop/bin
```

Finally add the binary into the `PATH` variable
```bash
echo 'export PATH=$HOME/.backdrop/bin:$PATH' >> $HOME/.bashrc
source $HOME/.bashrc
```

### From Source code
You can install it by compiling the source code
```bash
git clone https://github.com/JanMichaelSE/backdrop.git
cd backdrop
go build
```

Optionally move it to a directory in your PATH.
```bash
mv backdrop $HOME/.backdrop/bin
```

Finally add the binary into the `PATH` variable
```bash
echo 'export PATH=$HOME/.backdrop/bin:$PATH' >> $HOME/.bashrc
source $HOME/.bashrc
```

## :wastebasket: Uninstall

```bash
rm $(go env GOPATH)/bin/backdrop
go clean -modcache # <- Optional
```

## &#x2705; What's Supported

#### Operating Systems
- Ubuntu/GNOME based Distros
- CentOS/MATE (Doesn't support Slideshows)
- Coming Soon:
    - MacOS
    - Microsoft (Make a PR because I won't do it)
- Will not be Supported:
    - WSL

## :computer: Usage

Backdrop provides several options for managing your wallpapers:

- `-p, --path <PATH>`: 
    - Set a custom path to find wallpaper images. If not provided, a default path will be used.
- `-h, --help`: 
    - Displays help information on how to use the command, listing all available options.
- `-s, --slideshow`: 
    - Will configure and set a custom slideshow of images you select with fzf. To select multiple images hit "Tab" on the images you desire to select, then hit "Enter" to confirm.
- `-u, --url`: 
    - You will be prompted to provide an image url to be set as wallpaper. The image will be downloaded and previewed. If confirmed, the image will be downloaded to the directory were all images are found (check "IMAGES" section). If image is NOT accepted by user, the image gets deleted and previous wallpaper is set.
- `-v, --version`: 
    - Print version information.

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
