# Piñata (alpha)
Play blindfold chess against any UCI compatible engines like [Stockfish](https://stockfishchess.org/). 

## Usage
Piñata looks for `stockfish` in the standard executable search path by default. If not found 
```
Usage:
  pinata [flags]

Flags:
  -b, --black           choose the black side
  -d, --depth int       engine search depth (default 10)
  -e, --engine string   path to UCI compatible chess engine executable (default "stockfish")
  -h, --help            help for pinata
  -l, --light           invert the colors for lighter console background
      --no-color        disable colors
      --version         version for pinata
  -v, --visual          cheat blindfold
```

## Playing Blind
By default, the computer engine plays black. You make your first move. Use <TAB> to auto-complete possible moves or commands.
```
$ ./pinata 
█ 🙇  e4
░ 🤖  e6
█ 🙇  Nf3 
░ 🤖  d5
█ 🙇 <TAB>
Ke2      Qe2      Rg1      Be2      Bd3      Bc4      Bb5+     Ba6      Na3      Nc3      Ng1      Nd4      Nh4      Ne5      Ng5      
a3       a4       b3       b4       c3       c4       d3       d4       g3       g4       h3       h4       exd5     e5       resign   
/visual  /quit    /keys
```
## Playing Visual
You can cheat the blindfold with `--visual` flag and play interactively. Use `/visual` command to toggle the board display during the practice sessions to verify your memory.
```
$ ./pinata --visual
█ 🙇  e4
░ 🤖  e6
┼───┼───┼───┼───┼───┼───┼───┼───┼───┼
│   │ A │ B │ C │ D │ E │ F │ G │ H │
┼───┼───┼───┼───┼───┼───┼───┼───┼───┼
│ 8 │ ♖ │ ♘ │ ♗ │ ♕ │ ♔ │ ♗ │ ♘ │ ♖ │
┼───┼───┼───┼───┼───┼───┼───┼───┼───┼
│ 7 │ ♙ │ ♙ │ ♙ │ ♙ │   │ ♙ │ ♙ │ ♙ │
┼───┼───┼───┼───┼───┼───┼───┼───┼───┼
│ 6 │   │   │   │   │ ♙ │   │   │   │
┼───┼───┼───┼───┼───┼───┼───┼───┼───┼
│ 5 │   │   │   │   │   │   │   │   │
┼───┼───┼───┼───┼───┼───┼───┼───┼───┼
│ 4 │   │   │   │   │ ♟ │   │   │   │
┼───┼───┼───┼───┼───┼───┼───┼───┼───┼
│ 3 │   │   │   │   │   │   │   │   │
┼───┼───┼───┼───┼───┼───┼───┼───┼───┼
│ 2 │ ♟ │ ♟ │ ♟ │ ♟ │   │ ♟ │ ♟ │ ♟ │
┼───┼───┼───┼───┼───┼───┼───┼───┼───┼
│ 1 │ ♜ │ ♞ │ ♝ │ ♛ │ ♚ │ ♝ │ ♞ │ ♜ │
┼───┼───┼───┼───┼───┼───┼───┼───┼───┼
(W)
```
## Contribute to Piñata Project
Please follow Piñata [Contributor's Guide](https://github.com/abperiasamy/pinata/blob/master/code_of_conduct.md)

## License
Piñata is free software, licensed under [GNU AGPL v3 or later](https://github.com/abperiasamy/pinata/blob/master/LICENSE)
