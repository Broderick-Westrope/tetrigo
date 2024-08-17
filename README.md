# Tetrigo

*teh·tree·go*

![app demo](./docs/readme-demo.gif)

A Golang implementation of Tetris, attempting to follow the official [2009 Tetris Design Guideline](./docs/2009-Tetris-Design-Guideline.pdf).

This project is a work in progress. It consists of three main components:
1. `pkg/tetris/`: The core Tetris logic, including things like Tetrminimos, the Matrix, and scoring. This can be used to create game modes with your own ruleset and requirements.
2. `pkg/tetris/modes/`: The functionality for different Tetris game modes. This can be used to easily create a Tetris game with your own UI but without needing to know the ruleset.
3. `cmd/tetrigo/`: A TUI (Text User Interface) allowing you to play it out of the box. It also serves as a demonstration on how to use the packages and how to create a TUI using [Bubble Tea](https://github.com/charmbracelet/bubbletea).

Please feel free to open issues with suggestions, bugs, etc.

## Installation

Tetrigo can be installed by downloading the binary or by building from source. See the instructions below for your preferred method.

### Binary

You can download the binary corresponding to your operating system from the [releases page on GitHub].

Once downloaded you can run the binary from the command line:

```bash
# Linux or macOS
./tetrigo

# Windows
tetrigo.exe
```

Optionally, you can move the binary to a directory in your `$PATH` to run it from anywhere ([example](https://gist.github.com/nex3/c395b2f8fd4b02068be37c961301caa7)).

### Build From Source

Ensure that you have a supported version of Go properly installed and setup. You can find the minimum required version of Go in the [go.mod](./go.mod) file.

You can then install the latest release globally by running:

```bash
go install github.com/Broderick-Westrope/tetrigo/cmd/tetrigo@latest
```

Or you can install into another directory:

```bash
env GOBIN=/bin go install github.com/Broderick-Westrope/tetrigo/cmd/tetrigo@latest
```

## Usage

### Controls

The default game controls are as follows:

- **Move Left**: `A`
- **Move Right**: `D`
- **Toggle Soft Drop On/Off**: `S`
- **Hard Drop**: `W`
- **Rotate Clockwise**: `E`
- **Rotate Counter-Clockwise**: `Q`
- **Hold Tetrimino / Submit menu option**: `Space` or `Enter`
- **Pause Game / Exit**: `Escape`
- **Force Quit game**: `Ctrl+C`
- **Show Controls Help**: `?`

The game controls can be changed in the configuration file.

The menu, leaderboard, etc can be navigated using the arrow keys (moving), escape (exit), and enter (submit). These controls are not configurable.

## Configuration

### CLI

Starting Tetrigo with no subcommand or flags will start the game in the menu where you can manually configure simple settings like the player name and game mode:

```bash
./tetrigo
```

You're also able to start the game directly in a game mode (eg. Marathon), skipping the menu:

```bash
# Start the game in Marathon mode with a level of 5 and the player name "Brodie"
./tetrigo marathon --level=5 --name=Brodie 
```

To see more options for starting the game you can run:

```bash
./tetrigo --help
```

### TOML

More complex configuration can be done using a TOML file. If no config file is found sensible defaults will be used.

By default, Tetrigo will look for the file `config.toml` in the working directory. You can specify a different file using the `--config` flag.

```bash
./tetrigo --config=/path/to/config.toml
```

An example configuration file is provided in [`example.config.toml`](./example.config.toml).

## Data

The game data is stored in a SQLite database. By default, the database is stored in the working directory as `tetrigo.db`. You can specify a different file using the `--db` flag.

```bash
./tetrigo --db=/path/to/data.db
```

## TODO

- Check for lockdown 0.5s after landing on a surface
  - Also on soft drop, but not on hard drop.
  - This resets after each movement & rotation, for a total of 15 movements/rotations.
  - See "Extended Placement Lock Down" in the design guidelines.
- Score points from T-Spins
- More game modes
  - Sprint
  - Ultra 
- Multiplayer