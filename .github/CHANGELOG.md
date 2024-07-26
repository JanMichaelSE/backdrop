# Change Log

## Minor bug Fix

* Added Compression step in CI/CD pipepline to ensure proper `tar.gz` file be created. This was added with the commit ([8fcc5a4](https://github.com/JanMichaelSE/backdrop/commit/8fcc5a4))
* Fixed issue where default value `N` was concidered a invalid input !["Invalid user input"](./images/user-input-decision.png)

## Code improvement

* Changed hard-coded filepath to wallpapers folder, used `filepath.Join()` method instead.
