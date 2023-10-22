#!/usr/bin/env bash

append_backdrop_to_path() {
   local CONFIG_FILE="$1"
   local CONFIG_PATH="$HOME/$CONFIG_FILE"

   echo '' >> "$CONFIG_PATH"
   echo "# Appended Backdrop to PATH, done by install script." >> "$CONFIG_PATH"
   if [[ $CONFIG_FILE = '.zshrc' || $CONFIG_FILE = '.bashrc' ]]; then
      echo "Adding to bashrc/zshrc"
      echo 'export PATH=$PATH:$HOME/.backdrop/bin' >> "$CONFIG_PATH"
   else
      echo "Adding to fish"
      echo 'set -gx PATH $HOME/.backdrop/bin $PATH' >> "$CONFIG_PATH"
   fi
   echo '' >> "$CONFIG_PATH"

   echo '------------------'
   echo 'IMPORTANT:'
   echo "  - Remember to SOURCE your $CONFIG_PATH for changes to take affect. If not re-open your terminal emulator."
   echo '------------------'
}

if [[ -f "$HOME/.zshrc" ]]; then
   if ! grep -q ".backdrop/bin" "$HOME/.zshrc"; then
      echo "Backdrop not in zshrc, adding PATH."
      append_backdrop_to_path ".zshrc"
   fi
fi

if [[ -f "$HOME/.bashrc" ]]; then
   if ! grep -q ".backdrop/bin" "$HOME/.bashrc"; then
      echo "Backdrop not in bashrc, adding PATH."
      append_backdrop_to_path ".bashrc"
   fi
fi

if [[ -f "$HOME/.config/config.fish" ]]; then
   if ! grep -q ".backdrop/bin" "$HOME/.config/config.fish"; then
      echo "Backdrop not in config.fish, adding PATH."
      append_backdrop_to_path ".config/config.fish"
   fi
fi

if [[ ! -d "$HOME/.backdrop" ]]; then
   echo "Configuring '.backdrop/' at $HOME dir."
   mkdir -p "$HOME/.backdrop/bin"
   cp -p "./backdrop.sh" "$HOME/.backdrop/bin/backdrop"
   echo "Successfully configured backdrop!"
else
   echo "Backdrop already configured."
fi

exit 0
