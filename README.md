# Piñata (v1.4)
Piñata is an interactive shell to play blindfold chess against computers. Install any UCI compatible chess engine like [Stockfish](https://stockfishchess.org/download/) in the standard executable search path and Piñata will pick it up.

## Download
| OS         | Arch           | Link                                                                                                   |
| ---------- | --------       | ------                                                                                                 |
| GNU/Linux  | Intel (64-bit) | [Download](https://github.com/abperiasamy/pinata/releases/download/v1.4/pinata_1.4_darwin_x64.tar.gz)   |
| GNU/Linux  | Arm (64-bit)   | [Download](https://github.com/abperiasamy/pinata/releases/download/v1.4/pinata_1.4_linux_a64.tar.gz)   |
| GNU/Linux  | Arm (32-bit)   | [Download](https://github.com/abperiasamy/pinata/releases/download/v1.4/pinata_1.4_linux_a32v7.tar.gz) |
| Windows    | Intel (64-bit) | [Download](https://github.com/abperiasamy/pinata/releases/download/v1.4/pinata_1.4_windows_x64.zip)    |
| Darwin     | Intel (64-bit) | [Download](https://github.com/abperiasamy/pinata/releases/download/v1.4/pinata_1.4_darwin_x64.tar.gz)  |
| FreeBSD    | Intel (64-bit) | [Download](https://github.com/abperiasamy/pinata/releases/download/v1.4/pinata_1.4_freebsd_x64.tar.gz) |

## Usage
```
Usage:
  pinata [flags]

Flags:
  -b, --black           choose the black side
  -d, --depth int       engine search depth (default 10)
  -e, --engine string   path to UCI compatible chess engine executable (default "stockfish")
  -f, --file string     load game from a PGN file
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
/fen     /save    /load    /visual  /quit    /keys    /fen     /visual  /quit    /keys
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
█ 🙇
```
## Contribute to Piñata Project
Please follow Piñata [Contributor's Guide](https://github.com/abperiasamy/pinata/blob/master/code_of_conduct.md)

## Credits
- [Chess library](https://github.com/notnil/chess) by Logan Spears (notnil)
- [UCI library](https://github.com/freeeve/uci) by Eve Freeman (freeeve)

## License
Piñata is free software, licensed under [GNU AGPL v3 or later](https://github.com/abperiasamy/pinata/blob/master/LICENSE)
