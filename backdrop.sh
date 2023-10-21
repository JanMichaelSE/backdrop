# FEATURES
# 1.) Needs to run just using the command at first. "backdrop"
# 2.) Needs to search in specific directories to list out images that can be selected as background
# 3.) After selection:
#   - MVP: Just set the background and exit.
#   - IDEAL: Show preview and ask for confirmation before saving. If rejected revert to previous background.

# TODO MVP:
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
    echo '  -r, --revert         Reverts to the last wallpaper that was set prior to the most recent change.'
    echo "  -h, --help           Displays help information on how to use the ${0} command, listing all"
    echo '                       available options.'
    exit 1
}

# TODO: Might not be needed. But have it here just incase
check_command_status() {
    # Check if last command was successful
    if [[ "${?}" -ne 0 ]]; then
      echo "${1} was not success full." >&2
      exit 1
    fi
}

# 1. Need to specify PATHS to search images under.
PICTURES_PATH="$HOME/Pictures/wallpapers" # Not sure about wallpapers but whatever for now in Pictures.
CONFIGS_PATH="$HOME/.config/backdrop/wallpapers"
CUSTOM_PATH=""

# This is how to read long and short options, the "--" is to know when we are done.
OPTIONS=$(getopt -o p:rh -l path:,revert,help -- "$@")
check_command_status "Getting command options"

# Reorder the arguments to ensure they are correct
eval set -- "$OPTIONS"

while true; do
    case "$1" in
        -p|--path)
            shift # Move to next argument to get the path value passed.
            CUSTOM_PATH="$1"
            if [[ ! -d "$HOME/$CUSTOM_PATH" && ! -d "$CUSTOM_PATH" ]]; then
               echo "The directory those not exist. Need to provide a path under HOME or absolute path." >&2
               exit 1
            fi
            # TODO: Need to see how to make this persist? Config file maybe?
            exit 0
            ;;
        -r|--revert)
            echo "Reverting image... STILL PENDING"
            # TODO: Need to see how to make this persist? Config file maybe?
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

# 2. Need to grab the file names from those directories with Full Path
# FUTURE TODO: Need to see how to handle subfolders
# FUTURE TODO: Need to see how to handle multiple valid directories
#   For now just gonna give priority to ".config/backdrop/wallpapers" if exists.
# FUTURE TODO: Need to see how to manage CUSTOM_PATH to persist and provide it as an option.

# Check if image directories exist, give priority to ".config".
if [[ -d $CONFIGS_PATH ]]; then
    WALLPAPERS=$(find "$CONFIGS_PATH" -maxdepth 1 -type f | awk -F '/' '{print $NF}')
elif [[ -d $PICTURES_PATH ]]; then
    WALLPAPERS=$(find "$PICTURES_PATH" -maxdepth 1 -type f | awk -F '/' '{print $NF}')
else
    echo "No valid directories found to list images. Please assure you have one of the following configured:"
    echo "     - $CONFIGS_PATH"
    echo "     - $PICTURES_PATH"
    echo ""
    exit 1
fi

# Conver wallpapers to an array so we can reference them by index
IFS=$'\n' WALLPAPERS_ARRAY=($WALLPAPERS)

# Display options to user
i=1
for wallpaper in "${WALLPAPERS_ARRAY[@]}"; do
    echo "$i) $wallpaper"
    ((i++))
done

# Get user selection
while true; do
    read -p "Please select image background image (or press 'q' to quit): " REPLY
    if [[ $REPLY == 'q' ]]; then
       echo "Exiting backdrop..." 
       exit 0
    elif (( REPLY > 0 && REPLY <= ${#WALLPAPERS_ARRAY[@]} )); then
        SELECTED_WALLPAPAER=${WALLPAPERS_ARRAY[$REPLY-1]}

        gsettings set org.gnome.desktop.background picture-uri "file://$PICTURES_PATH/$SELECTED_WALLPAPAER"
        gsettings set org.gnome.desktop.background picture-uri-dark "file://$PICTURES_PATH/$SELECTED_WALLPAPAER"
        check_command_status "Changing background image"

        echo "Successfully changed background image."
        break
    else
        echo "Invalid selection, please try again."
    fi
done

exit 0

# Future Tasks:
# * See why it doesn't work with zsh, but it does work in Bash? (Might need to migrate to Golang by then.)
# * Provide flags to revert the last image selected "--revert or -r"
# * Provide a flag to just give a filename if the user knows it and automatically set that new bg image "--image or -i" (Flag name could change)
# * Provide a Preview Functionality before saving.
#   - This could re-use the revert functionality above. Function would be good.
# * See how a slide show implementation could fit here.
# * Need to see how to integrate "fzf" to fuzzy find the background image in the directory.
#   - Then the user hits enter and it previews the image.
#   - If confirmed the background will stay changed.
#   - If denied, the background will revert to the one the user had.
# * Make install script so tool is ready to be used by just running one script.
# * Make prompt experience more pretty (Low priority but it's bound to happen)
# * Super future: see how midjourney could be a cool integration with this tool.

