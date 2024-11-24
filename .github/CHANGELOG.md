# Change Log

## Support for Windows OS

- Was tested using windows 11.
- Supported features are the following:
  - Setting wallpaper from the `wallpapers folder path`
  - Setting wallpaper from url provided
  - Preview wallpaper from the `walppapers folder` before completely setting it as background

## Code improvement

- Changed hard-coded filepath to wallpapers folder, used `filepath.Join()` method instead.

- Replaced hard-coded `path trimming` for `filepath.Split` method within the `integration_test.go` file.

- Created OS specific functions with the purpose to abstract funtionality from the main code for cross-platform compability.
