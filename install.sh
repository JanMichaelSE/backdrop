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
      # TODO: Need to test this tomorrow when I get to work.
      sudo yum update -y
      sudo yum install fzf -y
   else
     echo "Unsupported Distribution/Operating System." >&2
     echo "Could not install fzf for this system. Make sure this is installed before using 'backdrop -f' feature." >&2
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

if [[ -f "$HOME/.config/config.fish" || -L "$HOME/.config/config.fish" ]]; then
   if ! grep -q ".backdrop/bin" "$HOME/.config/config.fish"; then
      echo ''
      echo "Backdrop not in config.fish, adding PATH."
      append_backdrop_to_path ".config/config.fish"
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
fi

exit 0
