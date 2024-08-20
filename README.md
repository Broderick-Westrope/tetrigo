# Tetrigo

*teh·tree·go*

![App Demo](./docs/readme-demo.gif)

**Tetrigo** is a Golang implementation of Tetris, following the official [2009 Tetris Design Guideline](./docs/2009-Tetris-Design-Guideline.pdf). The project is modular, offering components tailored to different needs—whether you want to play Tetris, create a custom Tetris UI, or develop new game modes.

## Contents

- [Introduction](#introduction)
- [Installation](#installation)
  - [Binary Installation](#binary-installation)
  - [Build from Source](#build-from-source)
- [Usage](#usage)
  - [Controls](#controls)
  - [CLI Usage](#cli-usage)
- [Configuration](#configuration)
  - [TOML Configuration](#toml-configuration)
- [Data Management](#data-management)
- [Development](#development)
  - [Project Structure](#project-structure)
  - [Building](#building)
  - [Testing](#testing)
- [TODO](#todo)
- [Contributing](#contributing)
- [License](#license)

## Introduction

Tetrigo is designed to be flexible and extendable. Depending on your goals, you can choose the appropriate part of the project to work with:

- **Play Tetris**: Use the TUI (Text User Interface) found in `cmd/tetrigo/`. See the [installation](#installation) section for setup instructions.
- **Create a Custom Tetris UI**: Use the game modes in `pkg/tetris/modes/` to integrate Tetris functionality into your own interface.
- **Develop New Game Modes**: Use the core logic in `pkg/tetris/` to design your own Tetris game modes with custom rules and mechanics.

## Installation

You can install Tetrigo by downloading the pre-built binary or by building the project from source.

### Binary Installation

1. Download the appropriate binary for your operating system from the [releases page on GitHub](https://github.com/Broderick-Westrope/tetrigo/releases).
2. Run the binary from the command line:

   ```bash
   # Linux or macOS
   tetrigo

   # Windows
   tetrigo.exe

Optionally, you can move the binary to a directory in your `$PATH` to run it from anywhere ([example](https://gist.github.com/nex3/c395b2f8fd4b02068be37c961301caa7)).

### Build From Source

Ensure that you have a supported version of Go properly installed and setup. You can find the minimum required version of Go in the [go.mod](./go.mod) file.

Using the Makefile:

You can then install the latest release globally by running:

```bash
go install github.com/Broderick-Westrope/tetrigo/cmd/tetrigo@latest
```

Or you can install into another directory:

```bash
env GOBIN=/bin go install github.com/Broderick-Westrope/tetrigo/cmd/tetrigo@latest
```

## Usage

For general information on how to play Tetris see [this beginners guide](https://tetris.com/article/33/tetris-tips-for-beginners).

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

### CLI Usage

Starting Tetrigo with no subcommand or flags will start the game in the menu where you can manually configure simple settings like the player name and game mode:

```bash
tetrigo
```

You're also able to start the game directly in a game mode (eg. Marathon), skipping the menu:

```bash
# Start the game in Marathon mode with a level of 5 and the player name "Brodie"
tetrigo marathon --level=5 --name=Brodie 
```

To see more options for starting the game you can run:

```bash
tetrigo --help
```

## Configuration

### TOML Configuration

More complex configuration can be done using a TOML file. If no config file is found sensible defaults will be used.

By default, Tetrigo will look for the file `config.toml` in the working directory. You can specify a different file using the `--config` flag.

```bash
tetrigo --config=/path/to/config.toml
```

An example configuration file is provided in [`example.config.toml`](./example.config.toml).

## Data Management

The game data is stored in a SQLite database. By default, the database is stored in the working directory as `tetrigo.db`. You can specify a different file using the `--db` flag.

```bash
tetrigo --db=/path/to/data.db
```

## Development

### Project Structure

This project consists of three main components:
1. `cmd/tetrigo/`: A TUI (Text User Interface) allowing you to play it out of the box. It also serves as a demonstration on how to use the packages and how to create a TUI using [Bubble Tea](https://github.com/charmbracelet/bubbletea).
2. `pkg/tetris/modes/`: The functionality for different Tetris game modes. This can be used to easily create a Tetris game with your own UI but without needing to know the ruleset.
3. `pkg/tetris/`: The core Tetris logic, including things like Tetrminimos, the Matrix, and scoring. This can be used to create game modes with your own ruleset and requirements.

[Task](https://taskfile.dev/) is the build tool used in this project. The Task config lives in [Taskfile.yaml](./Taskfile.yaml). Once the Task CLI is installed, you can see all available tasks by running:
```bash
task -l
```

You can run the TUI using the `run` task:
```bash
task run
```

### Building

You can build the project using the `build` task:
```bash
task build
```

This will create a binary in the `bin/` directory which can be run using the instructions in the [Installation](#installation) section.

### Testing

Tests can be run using the `test` task:
```bash
task test
```

You can also use the `cover` task to generate and open a coverage report:
```bash
task cover
```

The ordered priorities for testing are:
1. `pkg/tetris/`
2. `pkg/tetris/modes/`
3. `cmd/tetrigo/`

### TODO

- Add more tests
  - Finish Super Rotation System testing in [tetrimino_test.go](./pkg/tetris/tetrimino_test.go).
  - Add tests for scoring endOnMaxLevel.
  - Revisit scoring tests.
- Add the remaining Lock Down options.
- Check for Lock Down 0.5s after landing on a surface
  - Also on Soft Drop, but not on Hard Drop.
  - This resets after each movement & rotation, for a total of 15 movements/rotations.
  - See "Extended Placement Lock Down" in the design guidelines.
- Score points from T-Spins
- More game modes
  - Sprint
  - Ultra 
- Multiplayer

## Contributing
Contributions are welcome! If you have suggestions, bug reports, or feature requests, please open an issue on GitHub.

## License
This project is licensed under the MIT License. See the LICENSE file for details.
