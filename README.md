# Piñata (alpha)
----------------
Play blindfold chess against any UCI compatible engines like [Stockfish](https://stockfishchess.org/). 

## Usage
Piñata looks for `stockfish` in the standard executable search path by default. If not found 
```
Usage:
  pinata [flags]

Flags:
  -d, --depth int       engine search depth (default 10)
  -e, --engine string   path to UCI compatible chess engine executable (default "stockfish")
  -p, --play string     choose black or white (default "white")
  -v, --visual          cheat blindfold
      --no-color        disable colors 
      --version         version for pinata
  -h, --help            help for pinata
  
```

## Playing Blind
By default, the computer engine plays black. You make your first move. Use <TAB> to auto-complete possible moves or commands.
```
$ ./pinata 
[W]< e4
[B]> e5
[W]< Nc3 
[B]> f6
[W]< 
Ke2      Qe2      Qf3      Qg4      Qh5      Rb1      Be2      Bd3      Bc4      Bb5      Ba6      Nge2     Nf3      Nh3      Nb1
Nce2     Na4      Nb5      Nd5      a3       a4       b3       b4       d3       d4       f3       f4       g3       g4       h3
h4       resign   /visual  /quit    /keys
```
## Playing Visual
```
You can cheat the blindfold with `--visual` or `/visual` options and play interactively. Use `/visual` to toggle the board display in practice sessions to verify your memory.
$ ./pinata -v
(W) e4
(B) e6
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
Please follow Piñata [Contributor's Guide](https://github.com/abperiasamy/pinata/blob/master/CONTRIBUTING.md)

## License
Piñata is free software, licensed under [GNU AGPL v3 or later](https://github.com/abperiasamy/pinata/blob/master/LICENSE)
