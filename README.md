# Tetrigo

*teh·tree·go*

![app demo](./docs/assets/readme-demo.gif)

A Golang implementation of Tetris, attempting to follow the official [2009 Tetris Design Guideline](https://github.com/Broderick-Westrope/tetrigo/tree/main/docs/2009-Tetris-Design-Guideline.pdf).

This project is a work in progress. I have separated the included TUI (Terminal User Interface) from the game logic. This allows others to use the core logic to make their own game variants and other user interfaces.

Please feel free to open issues with suggestions, bugs, etc.

## TODO

Tasks marked with "MVP" are required for the "Minimum Viable Product". These are the tasks I want to complete before the next round of refactoring, thorough testing, etc.

- (MVP) Game over conditions
    - Game over screen
- (MVP) High Score system
- (MVP) Ghost piece
- More game modes
    - Sprint
    - Ultra 
- Multiplayer
- ~~Configuration file~~
- Check for lockdown 0.5s after landing on a surface
    - Also on soft drop, but not on hard drop
    - This resets after each movement & rotation, for a total of 15 movements/rotations.
    - See "Extended Placement Lock Down"
- ~~Drop one row immediately if nothing is blocking~~
- ~~Pause ('P' key?)~~
- ~~Fix Tetrimino rotation axis~~
- Implement SRS (Super Rotation System)
- ~~Score points from soft & hard drops~~
- T-Spins