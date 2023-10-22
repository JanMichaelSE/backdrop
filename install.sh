#!/bin/bash

if [[ -f "$HOME/.zshrc" ]]; then
   echo "Zshrc Config Exists" 
fi

if [[ -f "$HOME/.bashrc" ]]; then
   echo "Bash Config Exists" 
fi

if [[ -f "$HOME/.config/config.fish" ]]; then
   echo "Fish Config Exists" 
fi

exit 0
