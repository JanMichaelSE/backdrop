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
    echo '  -s, --slideshow      Will configure and set a custom slideshow of images you select with fzf.'
    echo '                       To select multiple images hit "Tab" on the images you desire to select, then hit "Enter" to'
    echo '                       confirm.'
    echo '  -u, --url            Provide an image url to be set as wallpaper. The image will be downloaded and previewed.'
    echo '                       If confirmed, the image will be downloaded to the directory were all images are found '
    echo '                       (check "IMAGES" section). If image is NOT accepted by user, the image gets deleted and previous '
    echo '                       wallpaper is set.'
    echo '  -v, --version        Print version information.'
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

expose_image_path_and_wallpapers() {
    local PICTURES_PATH="$HOME/Pictures/wallpapers"
    local CONFIGS_PATH="$HOME/.config/backdrop/wallpapers"

    # Chosen in order of priority.
    if [[ -d "$BACKDROP_IMAGE_PATH" ]]; then
        SELECTED_PATH=${BACKDROP_IMAGE_PATH}
        WALLPAPERS=$(find -L "${BACKDROP_IMAGE_PATH}" -maxdepth 1 -type f | awk -F '/' '{print $NF}')
    elif [[ -d "$CONFIGS_PATH" ]]; then
        SELECTED_PATH=${CONFIGS_PATH}
        WALLPAPERS=$(find -L "${CONFIGS_PATH}" -maxdepth 1 -type f | awk -F '/' '{print $NF}')
    elif [[ -d $PICTURES_PATH ]]; then
        SELECTED_PATH=${PICTURES_PATH}
        WALLPAPERS=$(find -L "${PICTURES_PATH}" -maxdepth 1 -type f | awk -F '/' '{print $NF}')
    else
        echo "No valid directories found to list images. Please assure you have one of the following configured:"
        echo "     - $CONFIGS_PATH"
        echo "     - $PICTURES_PATH"
        echo '     - Set a custom path with the "--path" or "-p" flag.'
        echo ""
        exit 1
    fi
}

get_previous_wallpaper() {
    if gsettings list-schemas | grep -iq mate.background; then
        echo "$(gsettings get org.mate.background picture-filename)"
    elif gsettings list-schemas | grep -iq gnome.desktop.background; then
        echo "$(gsettings get org.gnome.desktop.background picture-uri | awk -F "://" '{print $2}' | sed "s/'//g")"
    fi
}

set_wallpaper() {
    if gsettings list-schemas | grep -iq mate.background; then
        gsettings set org.mate.background picture-filename "$1"
    elif gsettings list-schemas | grep -iq gnome.desktop.background; then
        gsettings set org.gnome.desktop.background picture-uri "file://$1"
        gsettings set org.gnome.desktop.background picture-uri-dark "file://$1"
    fi
    check_command_status "Changing background image"
}

append_backdrop_custom_image_path() {
   local CONFIG_FILE="$1"
   local CONFIG_PATH="$HOME/$CONFIG_FILE"
   local IMAGE_PATH="$2"

   echo '' >> "$CONFIG_PATH"
   echo "# Setup Backdrop Custom Image Path." >> "$CONFIG_PATH"
   if [[ $CONFIG_FILE = '.zshrc' || $CONFIG_FILE = '.bashrc' ]]; then
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

setup_fzf_wallpaper() {
    expose_image_path_and_wallpapers

    while true; do
        PREVIOUS_WALLPAPER=$(get_previous_wallpaper)
        SELECTED_WALLPAPER=$(find -L "$SELECTED_PATH" -maxdepth 1 -type f | awk -F '/' '{print $NF}' | fzf --layout=reverse)

        if [[ -f "$SELECTED_PATH/$SELECTED_WALLPAPER" ]]; then
            set_wallpaper "$SELECTED_PATH/$SELECTED_WALLPAPER"
        else
            echo "No image selected, exiting..."
            break
        fi

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
}

setup_slideshow() {
    expose_image_path_and_wallpapers

    local PREVIOUS_WALLPAPER=$(get_previous_wallpaper)
    local SELECTED_WALLPAPERS=$(find -L "$SELECTED_PATH" -maxdepth 1 -type f | awk -F '/' '{print $NF}' | fzf --layout=reverse --multi)
    mapfile -t WALLPAPERS_ARRAY <<< "$SELECTED_WALLPAPERS"

    # Exit if no wallpaper was selected.
    if [[ -z "${WALLPAPERS_ARRAY[0]}" ]]; then
        echo "No image selected, exiting..."
        exit 0
    fi

    # Ask user for duration per slide
    read -p "What should be the duration per slide? (In Seconds): " SLIDE_DURATION

    # Create the slideshow folder if it doesn't exist
    local SLIDESHOW_PATH="$HOME/.local/share/gnome-background-properties"
    local SLIDESHOW_CONFIG_PATH="$HOME/.local/share/backgrounds/backdrop_settings"
    if [[ ! -d "$SLIDESHOW_PATH"  || ! -d "$SLIDESHOW_CONFIG_PATH" ]]; then
       mkdir -p "$SLIDESHOW_PATH"
       mkdir -p "$SLIDESHOW_CONFIG_PATH"
    fi

    local SLIDESHOW_FILE="$SLIDESHOW_PATH/backdrop_slideshow.xml"
    local SLIDESHOW_CONFIG_FILE="$SLIDESHOW_CONFIG_PATH/backdrop_settings.xml"

    # Create First Xml file that holds Slideshow images
    echo '<?xml version="1.0" encoding="UTF-8"?>' > "$SLIDESHOW_FILE"
    echo '<!DOCTYPE wallpapers SYSTEM "gnome-wp-list.dtd">' >> "$SLIDESHOW_FILE"
    echo '<wallpapers>' >> "$SLIDESHOW_FILE"
    echo '  <wallpaper>' >> "$SLIDESHOW_FILE"
    echo "    <name>Backdrop Slideshow</name>" >> "$SLIDESHOW_FILE"
    echo "    <filename>$SLIDESHOW_CONFIG_FILE</filename>" >> "$SLIDESHOW_FILE"
    echo '    <options>zoom</options>' >> "$SLIDESHOW_FILE"
    echo '    <pcolor>#2c001e</pcolor>' >> "$SLIDESHOW_FILE"
    echo '    <scolor>#2c001e</scolor>' >> "$SLIDESHOW_FILE"
    echo '    <shade_type>solid</shade_type>' >> "$SLIDESHOW_FILE"
    echo '  </wallpaper>' >> "$SLIDESHOW_FILE"
    echo '</wallpapers>' >> "$SLIDESHOW_FILE"

    # Create Second Xml file that holds Slideshow configuration
    echo '<background>' > "$SLIDESHOW_CONFIG_FILE"
    echo '  <starttime>' >> "$SLIDESHOW_CONFIG_FILE"
    echo '    <year>2012</year>' >> "$SLIDESHOW_CONFIG_FILE"
    echo '    <month>01</month>' >> "$SLIDESHOW_CONFIG_FILE"
    echo '    <day>01</day>' >> "$SLIDESHOW_CONFIG_FILE"
    echo '    <hour>00</hour>' >> "$SLIDESHOW_CONFIG_FILE"
    echo '    <minute>00</minute>' >> "$SLIDESHOW_CONFIG_FILE"
    echo '    <second>00</second>' >> "$SLIDESHOW_CONFIG_FILE"
    echo '  </starttime>' >> "$SLIDESHOW_CONFIG_FILE"

    local TOTAL_LENGTH=$((${#WALLPAPERS_ARRAY[@]} - 1))
    for index in "${!WALLPAPERS_ARRAY[@]}"; do
      echo '  <static>' >> "$SLIDESHOW_CONFIG_FILE"
      echo "    <duration>${SLIDE_DURATION}.0</duration>" >> "$SLIDESHOW_CONFIG_FILE"
      echo "    <file>$SELECTED_PATH/${WALLPAPERS_ARRAY[$index]}</file>" >> "$SLIDESHOW_CONFIG_FILE"
      echo '  </static>' >> "$SLIDESHOW_CONFIG_FILE"
      echo '  <transition>' >> "$SLIDESHOW_CONFIG_FILE"
      echo '    <duration>0.5</duration>' >> "$SLIDESHOW_CONFIG_FILE"
      echo "    <from>$SELECTED_PATH/${WALLPAPERS_ARRAY[$index]}</from>" >> "$SLIDESHOW_CONFIG_FILE"
      echo "    <to>$SELECTED_PATH/${WALLPAPERS_ARRAY[$((index + 1))]}</to>" >> "$SLIDESHOW_CONFIG_FILE"
      echo '  </transition>' >> "$SLIDESHOW_CONFIG_FILE"

      # Break early from for loop
      if [[ $((index + 1)) -eq $TOTAL_LENGTH ]]; then
          echo '  <static>' >> "$SLIDESHOW_CONFIG_FILE"
          echo "    <duration>${SLIDE_DURATION}.0</duration>" >> "$SLIDESHOW_CONFIG_FILE"
          echo "    <file>$SELECTED_PATH/${WALLPAPERS_ARRAY[$((index + 1))]}</file>" >> "$SLIDESHOW_CONFIG_FILE"
          echo '  </static>' >> "$SLIDESHOW_CONFIG_FILE"
          break 2
      fi
    done

    # Transition again back to the first image.
    echo '  <transition>' >> "$SLIDESHOW_CONFIG_FILE"
    echo '    <duration>0.5</duration>' >> "$SLIDESHOW_CONFIG_FILE"
    echo "    <from>$SELECTED_PATH/${WALLPAPERS_ARRAY[$TOTAL_LENGTH]}</from>" >> "$SLIDESHOW_CONFIG_FILE"
    echo "    <to>$SELECTED_PATH/${WALLPAPERS_ARRAY[0]}</to>" >> "$SLIDESHOW_CONFIG_FILE"
    echo '  </transition>' >> "$SLIDESHOW_CONFIG_FILE"
    echo '</background>' >> "$SLIDESHOW_CONFIG_FILE"

    set_wallpaper "$SLIDESHOW_CONFIG_FILE"
}

setup_url_image() {
    expose_image_path_and_wallpapers

    while true; do
        PREVIOUS_WALLPAPER=$(get_previous_wallpaper)

        read -rp 'Provide Image Url: ' IMAGE_URL
        if [[ $IMAGE_URL = '' ]]; then
            echo "No value provided. Exiting..."
            return
        fi

        local URL_IMAGES_PATH="$PWD/url_images"
        if [[ ! -d "$URL_IMAGES_PATH" ]]; then
            mkdir -p "$URL_IMAGES_PATH" 
        fi

        wget -nv -P "$URL_IMAGES_PATH" "$IMAGE_URL"
        check_command_status "Getting image from url"

        IMAGE_FROM_URL=$(ls "$URL_IMAGES_PATH")
        set_wallpaper "$URL_IMAGES_PATH/$IMAGE_FROM_URL"

        while true; do
            read -p "Want to save this change? [y/N]: " CHOICE
            case "$CHOICE" in
                [yY])
                    mv $URL_IMAGES_PATH/* "$SELECTED_PATH"
                    set_wallpaper "$SELECTED_PATH/$IMAGE_FROM_URL"
                    rm -rf "$URL_IMAGES_PATH"
                    echo "Successfully changed background image."
                    return
                    ;;
                [nN]|"")
                    set_wallpaper "$PREVIOUS_WALLPAPER"
                    rm -rf "$URL_IMAGES_PATH"
                    break
                    ;;
                *)
                    echo "Invalid input..."
                    ;;
            esac
        done
    done
}

# This is how to read long and short options, the "--" is to know when we are done.
OPTIONS=$(getopt -o p:fsuvh -l path:,fuzzy,slideshow,url,version,help,uninstall -- "$@")
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
            setup_fzf_wallpaper
            check_command_status "Setting Wallpaper with fzf"
            exit 0
            ;;
        -s|--slideshow)
            setup_slideshow
            check_command_status "Setting Slideshow configuration"
            echo "Successfully configure slideshow!"
            exit 0
            ;;
        -u|--url)
            setup_url_image
            check_command_status "Setting Url Image"
            echo "Successfully Setup Image from Url!"
            exit 0
            ;;
        --uninstall)
            echo "Uninstalling Backdrop..."
            "$HOME/.backdrop/scripts/uninstall.sh"
            exit 0
            ;;
        -v|--version)
            VERSION="v0.0.1"
            echo "Backdrop $VERSION"
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



# Convert wallpapers to an array so we can reference them by index
expose_image_path_and_wallpapers
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
        PREVIOUS_WALLPAPER=$(get_previous_wallpaper)
        SELECTED_WALLPAPER=${WALLPAPERS_ARRAY[$REPLY-1]}

        set_wallpaper "$SELECTED_PATH/$SELECTED_WALLPAPER"

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


exit 0

# Future Tasks:
# * Add an update flag to update the software without having to uninstall and re-install.
# * Provide flags to revert the last image selected "--revert or -r" (Optional, still thinking it's use)
# * Need to see how to handle subfolders
# * Need to see how to handle multiple valid directories
#       For now just gonna give priority to ".config/backdrop/wallpapers" if exists.
# * Make prompt experience more pretty (Low priority but it's bound to happen)
# * Add support for the following platforms:
#   - Add support for Mac (For Omar)
# * Super future: see how midjourney or DALL-E could be a cool integration with this tool.
#   - The user could be given a prompt to generate a wallpaper.
#   - Midjourney could provide 4 images that are presented by using the URL path provided?
#   - If the user likes one and accepts it then its downloaded to the machine and saved in his folder.
#   - If the user does not like the image then he could do another try to get 4 more images based on his last prompt.
#   - The user can also quit and provide a new prompt if he desires to do so.
#   - Must show a count of available midjourney prompts to the user so he knows at all times how
#       many times he can use this tool daily.

