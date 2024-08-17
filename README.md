# Tetrigo

*teh·tree·go*

![app demo](./docs/assets/readme-demo.gif)

A Golang implementation of Tetris, attempting to follow the official [2009 Tetris Design Guideline](https://github.com/Broderick-Westrope/tetrigo/tree/main/docs/2009-Tetris-Design-Guideline.pdf).

This project is a work in progress. It consists of three main components:
1. `pkg/tetris/`: The core Tetris logic, including things like Tetrminimos, the Matrix, and scoring. This can be used to create game modes with your own ruleset and requirements.
2. `pkg/tetris/modes/`: The functionality for different Tetris game modes. This can be used to easily create a Tetris game with your own UI but without needing to know the ruleset.
3. `cmd/tui/`: A TUI (Text User Interface) allowing you to play it out of the box. It also serves as a demonstration on how to use the packages and how to create a TUI using [Bubble Tea](https://github.com/charmbracelet/bubbletea).

Please feel free to open issues with suggestions, bugs, etc.

## TODO

- (MVP) Ghost piece
- Multiplayer
- Do more with the configuration file
- Check for lockdown 0.5s after landing on a surface
    - Also on soft drop, but not on hard drop.
    - This resets after each movement & rotation, for a total of 15 movements/rotations.
    - See "Extended Placement Lock Down" in the design guidelines.
- Implement SRS (Super Rotation System)
- Score points from T-Spins
- More game modes
  - Sprint
  - Ultra 
- Multiplayer