# Universal Chess Interface

Small wrapper for interacting with uci compliant chess engines.
[Godoc](https://godoc.org/github.com/dolegi/uci)


# Usage
Import library
```go
import "github.com/dolegi/uci"
```
Create a new engine
```go
eng, _ := uci.NewEngine("path-to-engine")
```
Set options
```go
eng.SetOption("Ponder", false)
eng.SetOption("Threads", "2")
```
Check if ready
```go
eng.IsReady()
```
Create a new game, use full algorithmic positioning and assign white
```go
eng.NewGame(uci.NewGameOpts{uci.ALG, uci.W})
```
Set the moves after start position
```go
eng.Position("e2e4")
```
Find best move
```go
resp := eng.Go(uci.GoOpts{MoveTime: 100})
fmt.Println(resp.Bestmove)
```

See example folder for a working example. Run with `go run example/main.go <path-to-uci-engine>`

# Missing features
- [ ] go infinite + stop
- [ ] ponderhit
- [ ] debug mode
- [ ] register

# References
- [UCI reference documention](https://www.shredderchess.com/chess-info/features/uci-universal-chess-interface.html)
- [Stockfish](https://github.com/official-stockfish/Stockfish)
