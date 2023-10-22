#!/usr/bin/env bash

if [[ -d "$HOME/.backdrop" ]]; then
    echo "Deleting '.backdrop' directory..."
    rm -rf "$HOME/.backdrop" 
fi

if [[ -f "$HOME/.zshrc" || -L "$HOME/.zshrc" ]]; then
    if grep -q ".backdrop/bin" "$HOME/.bashrc"; then
        sed -i --follow-symlinks "/backdrop/Id" "$HOME/.bashrc"
        echo "Removed PATH from '.bashrc'"
    fi
fi

if [[ -f "$HOME/.bashrc" || -L "$HOME/.bashrc" ]]; then
    if grep -q ".backdrop/bin" "$HOME/.zshrc"; then
        sed -i --follow-symlinks "/backdrop/Id" "$HOME/.zshrc"
        echo "Removed PATH from '.zshrc'"
    fi
fi

if [[ -f "$HOME/.config/config.fish" || -L "$HOME/.config/config.fish" ]]; then
    if grep -q ".backdrop/bin" "$HOME/.config/config.fish"; then
        sed -i --follow-symlinks "/backdrop/Id" "$HOME/.config/config.fish"
        echo "Removed PATH from '.config/config.fish'"
    fi
fi

echo "Successfully uninstalled Backdrop!"

exit 0
