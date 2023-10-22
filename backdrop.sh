#!/usr/bin/env bash

usage() {
    echo "Usage: ${0} [-rh] [-p PATH]"
    echo 'DESCRIPTION'
    echo '  backdrop is a command-line utility for managing wallpapers on your desktop.'
    echo '  It allows you to set a new wallpaper, revert to a previous wallpaper, and specify the directory'
    echo '  where your wallpaper images are stored.'
    echo ''
    echo 'OPTIONS'
    echo "  -p, --path <PATH>    Set a custom path to find wallpaper images. If not provided, a default"
    echo '                       path will be used.'
    echo '  -f, --fuzzy          Performs a fuzzy finding (Requires fzf).'
    echo "  -h, --help           Displays help information on how to use the ${0} command, listing all"
    echo '  --uninstall          Will uninstall Backdrop by removing all PATHs and Backdrop files.'
    echo '                       available options.'
    echo ''
    echo 'IMAGES'
    echo '  Images must be stored in ONE of the following paths:'
    echo "     - $HOME/.config/backdrop/wallpapers (This one has priority)"
    echo "     - $HOME/Pictures/wallpapers"
    echo '  Note: If "BACKDROP_IMAGE_PATH" shell variable is set, it will have priority and be used to list images.'
    echo '        This is set by using the "--path" or "-p" flag mentioned above.'
    exit 1
}

check_command_status() {
    # Check if last command was successful
    if [[ "${?}" -ne 0 ]]; then
      echo "${1} was not success full." >&2
      exit 1
    fi
}

select_image_path() {
    SELECTED_PATH=${1}
    WALLPAPERS=$(find -L "${1}" -maxdepth 1 -type f | awk -F '/' '{print $NF}')
}

set_wallpaper() {
    gsettings set org.gnome.desktop.background picture-uri "$1"
    gsettings set org.gnome.desktop.background picture-uri-dark "$1"
    check_command_status "Changing background image"
}

append_backdrop_custom_image_path() {
   local CONFIG_FILE="$1"
   local CONFIG_PATH="$HOME/$CONFIG_FILE"
   local IMAGE_PATH="$2"

   echo "Config Path: $CONFIG_PATH"
   echo "Image Path: $IMAGE_PATH"
   echo '' >> "$CONFIG_PATH"
   echo "# Setup Backdrop Custom Image Path." >> "$CONFIG_PATH"
   if [[ $CONFIG_FILE = '.zshrc' || $CONFIG_FILE = '.bashrc' ]]; then
       echo "Setting values for: $CONFIG_FILE"
      echo "export BACKDROP_IMAGE_PATH=$IMAGE_PATH" >> "$CONFIG_PATH"
   else
      echo "set -gx BACKDROP_IMAGE_PATH $IMAGE_PATH" >> "$CONFIG_PATH"
   fi
   echo '' >> "$CONFIG_PATH"
   echo "<<< Completed Setup for $CONFIG_FILE >>>"
   echo '------------------'
   echo 'IMPORTANT:'
   echo "  - Remember to SOURCE your $CONFIG_PATH for changes to take affect. If not re-open your terminal emulator."
   echo '-----------------'
   echo ''
}

set_custom_path() {
    local IMAGE_PATH="$1"
    if [[ -f "$HOME/.zshrc" || -L "$HOME/.zshrc" ]]; then
       if ! grep -q "BACKDROP_IMAGE_PATH" "$HOME/.zshrc"; then
          echo ''
          echo "Backdrop custom image path not in zshrc, adding BACKDROP_IMAGE_PATH."
          append_backdrop_custom_image_path ".zshrc" "$IMAGE_PATH"
       fi
    fi

    if [[ -f "$HOME/.bashrc" || -L "$HOME/.bashrc" ]]; then
       if ! grep -q "BACKDROP_IMAGE_PATH" "$HOME/.bashrc"; then
          echo ''
          echo "Backdrop custom image path not in bashrc, adding BACKDROP_IMAGE_PATH."
          append_backdrop_custom_image_path ".bashrc" "$IMAGE_PATH"
       fi
    fi

    if [[ -f "$HOME/.config/config.fish" || -L "$HOME/.config/config.fish" ]]; then
       if ! grep -q "BACKDROP_IMAGE_PATH" "$HOME/.config/config.fish"; then
          echo ''
          echo "Backdrop custom image path not in config.fish, adding BACKDROP_IMAGE_PATH."
          append_backdrop_custom_image_path ".config/config.fish" "$IMAGE_PATH"
       fi
    fi
}

# This is how to read long and short options, the "--" is to know when we are done.
OPTIONS=$(getopt -o p:fh -l path:,fuzzy,help,uninstall -- "$@")
check_command_status "Getting command options"

# Reorder the arguments to ensure they are correct
eval set -- "$OPTIONS"

# * Provide a flag to just give a filename if the user knows it and automatically set that new bg image "--image or -i" (Flag name could change)
while true; do
    case "$1" in
        -p|--path)
            shift # Move to next argument to get the path value passed.
            CUSTOM_IMAGE_PATH="$1"
            if [[ ! -d "$HOME/$CUSTOM_IMAGE_PATH" && ! -d "$CUSTOM_IMAGE_PATH" ]]; then
               echo "The directory those not exist. Need to provide a path under HOME or absolute path." >&2
               exit 1
            fi
            set_custom_path "$CUSTOM_IMAGE_PATH"
            check_command_status "Setting Custom Image Path"
            echo "New custom path to find images has been set!"
            exit 0
            ;;
        -f|--fuzzy)
            IS_FUZZY_FINDING=true
            ;;
        --uninstall)
            echo "Uninstalling Backdrop..."
            "$HOME/.backdrop/scripts/uninstall.sh"
            exit 0
            ;;
        -h|--help) usage;;
        --)
            shift # End of opitons
            break
            ;;
        *) usage;;
    esac
    shift # Move to next option
done

# FUTURE TODO: Need to see how to handle subfolders
# FUTURE TODO: Need to see how to handle multiple valid directories
#   For now just gonna give priority to ".config/backdrop/wallpapers" if exists.

# Check if image directories exist.
PICTURES_PATH="$HOME/Pictures/wallpapers"
CONFIGS_PATH="$HOME/.config/backdrop/wallpapers"

# Chosen in order of priority.
if [[ -d "$BACKDROP_IMAGE_PATH" ]]; then
    select_image_path $BACKDROP_IMAGE_PATH
elif [[ -d "$CONFIGS_PATH" ]]; then
    select_image_path $CONFIGS_PATH
elif [[ -d $PICTURES_PATH ]]; then
    select_image_path $PICTURES_PATH
else
    echo "No valid directories found to list images. Please assure you have one of the following configured:"
    echo "     - $CONFIGS_PATH"
    echo "     - $PICTURES_PATH"
    echo '     - Set a custom path with the "--path" or "-p" flag.'
    echo ""
    exit 1
fi

# Get user selection
# TODO: NEED TO REFACTOR THIS
if [[ $IS_FUZZY_FINDING = 'true' ]]; then
    while true; do
        PREVIOUS_WALLPAPER=$(gsettings get org.gnome.desktop.background picture-uri)
        SELECTED_WALLPAPER=$(find -L "$SELECTED_PATH" -maxdepth 1 -type f | awk -F '/' '{print $NF}' | fzf --layout=reverse)

        set_wallpaper "file://$SELECTED_PATH/$SELECTED_WALLPAPER"

        while true; do
            read -p "Want to save this change? [y/N]: " CHOICE
            case "$CHOICE" in
                [yY])
                    echo "Successfully changed background image."
                    break 2
                    ;;
                [nN]|"")
                    set_wallpaper "$PREVIOUS_WALLPAPER"
                    break
                    ;;
                *)
                    echo "Invalid input..."
                    ;;
            esac
        done     
    done
else
    # Conver wallpapers to an array so we can reference them by index
    IFS=$'\n' WALLPAPERS_ARRAY=($WALLPAPERS)

    # Display options to user
    i=1
    for wallpaper in "${WALLPAPERS_ARRAY[@]}"; do
        echo "$i) $wallpaper"
        ((i++))
    done
    while true; do
        read -p "Please select image background image (or press 'q' to quit): " REPLY
        if [[ $REPLY == 'q' ]]; then
           echo "Exiting backdrop..." 
           exit 0
        elif (( REPLY > 0 && REPLY <= ${#WALLPAPERS_ARRAY[@]} )); then
            PREVIOUS_WALLPAPER=$(gsettings get org.gnome.desktop.background picture-uri)
            SELECTED_WALLPAPER=${WALLPAPERS_ARRAY[$REPLY-1]}

            set_wallpaper "file://$SELECTED_PATH/$SELECTED_WALLPAPER"

            while true; do
                read -p "Want to save this change? [y/N]: " CHOICE
                case "$CHOICE" in
                    [yY])
                        echo "Successfully changed background image."
                        break 2
                        ;;
                    [nN]|"")
                        set_wallpaper "$PREVIOUS_WALLPAPER"
                        break
                        ;;
                    *)
                        echo "Invalid input..."
                        ;;
                esac
            done     
        else
            echo "Invalid selection, please try again."
        fi
    done
fi


exit 0

# Future Tasks:
# * Improve Readme.md for usage on this tool and installation.
# * Provide flags to revert the last image selected "--revert or -r" (Optional, still thinking it's use)
# * See how a slide show implementation could fit here.
# * See why it doesn't work with zsh, but it does work in Bash? (Might need to migrate to Golang by then.)
#   - Then the user hits enter and it previews the image.
#   - If confirmed the background will stay changed.
#   - If denied, the background will revert to the one the user had.
# * Make prompt experience more pretty (Low priority but it's bound to happen)
# * Add support for CentOS (Because thats what I use at work)
# * Add support for Mac (For Omar)
# * Super future: see how midjourney could be a cool integration with this tool.

