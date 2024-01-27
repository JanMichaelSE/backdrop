package internal

import "errors"

// TODO: Might later move to a errors.go
var (
	ErrNoCompatibleDesktopEnvironment = errors.New("No compatible desktop environment found for setting wallpaper.")
	ErrUserCanceledSelection          = errors.New("User canceled selection, exiting program.")
	ErrCouldNotSetBackground          = errors.New("Error setting background wallpaper.")
	ErrCouldNotListSchemas            = errors.New("Error listing schemas.")
	ErrNoValidImagesPath              = errors.New(`User does not have valid images path configured.
    IMAGES
      Images must be stored in ONE of the following paths:
         - $HOME/.config/backdrop/wallpapers (This one has priority)
         - $HOME/Pictures/wallpapers
      Note: If "BACKDROP_IMAGE_PATH" shell variable is set, it will have priority and be used to list images.
            This is set by using the "--path" or "-p" flag mentioned above.
    `)
)
