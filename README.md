# mkrg [![CI Status](https://github.com/itchyny/mkrg/workflows/CI/badge.svg)](https://github.com/itchyny/mkrg/actions)
[Mackerel](https://mackerel.io) graph viewer in terminal.

## Screenshots
On iTerm2, export `MKRG_VIEWER=iTerm2` to get the graphical viewer (exporting the environment variable is required when you're using ssh).
![mkrg](https://user-images.githubusercontent.com/375258/47090208-65696e80-d25d-11e8-936a-3fe80879ebe7.png)

Sixel viewer is implemented but does not fit to the terminal width (patches welcolme).

The command has simple graph viewer using Braille.
![mkrg](https://user-images.githubusercontent.com/375258/47095115-8c2ca280-d267-11e8-99de-85dfb7401798.png)

## Installation
### Homebrew
```sh
brew install itchyny/tap/mkrg
```

### Build from source
```sh
go get github.com/itchyny/mkrg/cmd/mkrg
```

## Bug Tracker
Report bug at [Issues・itchyny/mkrg - GitHub](https://github.com/itchyny/mkrg/issues).

## Author
itchyny (https://github.com/itchyny)

## License
This software is released under the MIT License, see LICENSE.
