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
# 3. Need to Present those options to the user. (Use a Select)
# 4. Save user selection and change that background.
# 5. Provide a message confirming backdrop has been changed.

# Future Tasks:
# * Provide flags to revert the last image selected "--revert or -r"
# * Provide a flag to just give a filename if the user knows it and automatically set that new bg image "--image or -i" (Flag name could change)
# * Provide a Preview Functionality before saving.
#   - This could re-use the revert functionality above. Function would be good.
# * See how a slide show implementation could fit here.

exit 0
