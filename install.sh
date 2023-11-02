#!/usr/bin/env bash

append_backdrop_to_path() {
   local CONFIG_FILE="$1"
   local CONFIG_PATH="$HOME/$CONFIG_FILE"

   echo '' >> "$CONFIG_PATH"
   echo "# Appended Backdrop to PATH, done by install script." >> "$CONFIG_PATH"
   if [[ $CONFIG_FILE = '.zshrc' || $CONFIG_FILE = '.bashrc' ]]; then
      echo 'export PATH=$PATH:$HOME/.backdrop/bin' >> "$CONFIG_PATH"
   else
      echo 'set -gx PATH $HOME/.backdrop/bin $PATH' >> "$CONFIG_PATH"
   fi
   echo '' >> "$CONFIG_PATH"
   echo "Completed Setup for $CONFIG_FILE"
   echo '------------------'
   echo 'IMPORTANT:'
   echo "  - Remember to SOURCE your $CONFIG_PATH for changes to take affect. If not re-open your terminal emulator."
   echo '------------------'
   echo ''
}

install_fzf_based_on_os() {
   if grep -qi 'ubuntu' "/etc/os-release"; then
      sudo apt install fzf -y
   elif grep -qi 'centos' "/etc/os-release"; then
      git clone --depth 1 https://github.com/junegunn/fzf.git "$HOME/.fzf"
      "$HOME/.fzf/install"
   else
     echo "Unsupported Distribution/Operating System." >&2
     echo "Could not install fzf for this system. Make sure this is installed before using 'backdrop -f' feature." >&2
   fi
}

update_to_latest_version() {
   local BACKDROP_VERSION=$(./backdrop.sh -v | awk '{print $2}')
   local INSTALLED_VERSION=$(backdrop -v | awk '{print $2}')

   if [ "$BACKDROP_VERSION" == "$INSTALLED_VERSION" ]; then
      echo "You have the latest version of backdrop installed."
   else
      echo "A newer version of backdrop is available. Latest version: $BACKDROP_VERSION, installed version: $INSTALLED_VERSION"

      read -p "Do you want to update to the latest version? (Y/n) " choice

      case "$choice" in
         n | N)
            echo "Update cancelled."
            ;;
         y | Y | "")
            echo "Updating to the latest version..."
            cp -p "./backdrop.sh" "$HOME/.backdrop/bin/backdrop"
            ;;
         *)
            echo "Invalid choice."
            ;;
      esac
   fi

}

# Check if fzf is installed
echo -e "\n<<< Checking if fzf is installed. >>>"
if ! command -v "fzf" &> /dev/null; then
  echo "fzf is not installed. Installing..."
  install_fzf_based_on_os
else
  echo "fzf is already installed."
fi

if [[ -f "$HOME/.zshrc" || -L "$HOME/.zshrc" ]]; then
   if ! grep -q ".backdrop/bin" "$HOME/.zshrc"; then
      echo ''
      echo "Backdrop not in zshrc, adding PATH."
      append_backdrop_to_path ".zshrc"
   fi
fi

if [[ -f "$HOME/.bashrc" || -L "$HOME/.bashrc" ]]; then
   if ! grep -q ".backdrop/bin" "$HOME/.bashrc"; then
      echo ''
      echo "Backdrop not in bashrc, adding PATH."
      append_backdrop_to_path ".bashrc"
   fi
fi

if [[ -f "$HOME/.config/fish/config.fish" || -L "$HOME/.config/fish/config.fish" ]]; then
   if ! grep -q ".backdrop/bin" "$HOME/.config/fish/config.fish"; then
      echo ''
      echo "Backdrop not in config.fish, adding PATH."
      append_backdrop_to_path ".config/fish/config.fish"
   fi
fi

if [[ ! -d "$HOME/.backdrop" ]]; then
   echo "Configuring '.backdrop/' at $HOME dir."
   mkdir -p "$HOME/.backdrop/bin"
   mkdir -p "$HOME/.backdrop/scripts"
   cp -p "./backdrop.sh" "$HOME/.backdrop/bin/backdrop"
   cp -p "./uninstall.sh" "$HOME/.backdrop/scripts/uninstall.sh"
   echo "Successfully configured backdrop!"
else
   echo "Backdrop already configured."
   echo -e "\n<<< Verifying latest version >>>\n"
   update_to_latest_version
fi

exit 0
