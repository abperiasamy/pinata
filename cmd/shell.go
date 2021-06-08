/*
Copyright © 2020 Anand Babu Periasamy https://twitter.com/abperiasamy

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/abperiasamy/chess"
	"github.com/chzyer/readline"
)

// Readline input filter
func filterInput(r rune) (rune, bool) {
	switch r {
	/*
		// block CtrlZ feature
		case readline.CharCtrlZ:
			return r, false
	*/
	}
	return r, true
}

// Readline file listing
func completeLoad(path string) func(string) []string {
	return func(line string) []string {
		names := make([]string, 0)
		files, _ := ioutil.ReadDir(path)
		for _, f := range files {
			if strings.HasSuffix(f.Name(), ".pgn") {
				names = append(names, f.Name())
			}
		}
		return names
	}
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func shell() {
	// Initialize a new game and save it in global gGame.
	gGame = chess.NewGame(chess.UseNotation(chess.AlgebraicNotation{}))

	// Load game from PGN.
	if gGamePath != "" {
		filename := gGamePath
		// Append default name to dir if empty.
		fInfo, err := os.Stat(filename)
		if err == nil && fInfo.IsDir() {
			filename = filepath.Clean(filepath.Join(filename) + "/" + gGameFilename)
		}

		// Check if file exist.
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			fmt.Println(gConsole.Bold(gConsole.Red(filename)), "does not exist.")
			os.Exit(1)
		}

		if gGame = loadPGN(filename); gGame == nil { // Failed to load the PGN.
			// fmt.Println("Unable to open " + gConsole.Bold(gConsole.Red(filename)).String() + ".")
			os.Exit(1)
		}

		// Check to see if the game already ended.
		if isGameOver(gGame) {
			drawBoard(gGame)
			os.Exit(0)
		}
	}

	eng, err := newEngine(gEngineBinary)
	if err != nil {
		log.Fatal(err)
	}
	defer eng.Close()
	eng.SendOption("Threads", "8")

	completer := readline.NewPrefixCompleter(
		readline.PcItemDynamic(validMovesConstructor()),
		readline.PcItem("resign"),
		readline.PcItem("/fen"),
		readline.PcItem("/save", readline.PcItem(gGameFilename)),
		readline.PcItem("/load", readline.PcItemDynamic(completeLoad("."))),
		readline.PcItem("/visual"),
		readline.PcItem("/quit"),
		readline.PcItem("/keys",
			readline.PcItem("vi"),
			readline.PcItem("emacs"),
		))

	l, err := readline.NewEx(&readline.Config{
		// Prompt: "\033[31m»\033[0m ",
		// HistoryFile:     "/tmp/readline.tmp",
		AutoComplete:        completer,
		InterruptPrompt:     "/quit",
		EOFPrompt:           "\n",
		HistorySearchFold:   true,
		FuncFilterInputRune: filterInput,
	})
	if err != nil {
		panic(err)
	}
	defer l.Close()

	gameStarted := false

	if humanColor() == chess.Black {
		err = engineMoveFirst(eng, gGame)
		if err != nil {
			fmt.Errorf("Engine Failure: ", err)
			os.Exit(1)
		}
		gameStarted = true
	}

	for {
		l.SetPrompt(humanPrompt())
		cmd, err := l.Readline()
		if err == readline.ErrInterrupt {
			cmd = "/quit"
		}
		cmd = strings.TrimSpace(cmd)
		switch {
		case cmd == "": // no input, do nothing.

		case cmd == "resign":
			gGame.Resign(humanColor())
			isGameOver(gGame) // Game is over, but print the status.

			// Save the game.
			if savePGN(gGame, gGameFilename) == nil { // Success
				fmt.Println("Game saved to", gConsole.Bold(gConsole.Red(gGameFilename)))
			}

			goto end

		case strings.HasPrefix(cmd, "/fen"):
			cmd := strings.SplitN(cmd, " ", 2)
			if len(cmd) > 1 {
				fenStr := cmd[1]
				eng.SetFEN(fenStr)
				fen, err := chess.FEN(fenStr)
				if err != nil {
					fmt.Println("Not a valid FEN.")
					continue
				}
				gGame = chess.NewGame(fen)
				if isGameOver(gGame) { // No more moves to play.
					goto end
				}
			} else { // Just display the current FEN
				fmt.Println(gGame.FEN())
			}

		case strings.HasPrefix(cmd, "/load"):
			cmd := strings.SplitN(cmd, " ", 2)
			filename := gGameFilename
			if len(cmd) == 2 {
				filename = cmd[1]
			}

			// Append default name to dir if empty.
			fInfo, err := os.Stat(filename)
			if err == nil && fInfo.IsDir() {
				filename = filepath.Clean(filepath.Join(filename) + "/" + gGameFilename)
			}

			// Check if file exist.
			if _, err := os.Stat(filename); os.IsNotExist(err) {
				fmt.Println(gConsole.Bold(gConsole.Red(filename)), "does not exist")
				continue
			} else if !strings.HasSuffix(filename, ".pgn") {
				filename += ".pgn"
			}

			g := loadPGN(filename)
			if g != nil { // Success
				gGame = g              // Overwrite the current game.
				if isGameOver(gGame) { // No more moves to play.
					goto end
				}
			}

		case strings.HasPrefix(cmd, "/save"):
			cmd := strings.SplitN(cmd, " ", 2)
			filename := gGameFilename
			if len(cmd) == 2 {
				filename = strings.TrimSpace(cmd[1])
			}

			// Append default name to dir if empty.
			fInfo, err := os.Stat(filename)
			if err == nil && fInfo.IsDir() {
				filename = filepath.Clean(filepath.Join(filename) + "/" + gGameFilename)
			} else if !strings.HasSuffix(filename, ".pgn") {
				if strings.HasSuffix(filename, ".") { // avoid generating "..pgn"
					filename = strings.TrimSuffix(filename, ".")
				}
				filename += ".pgn"
			}

			if savePGN(gGame, filename) == nil { // Success
				fmt.Println("Game saved to", gConsole.Bold(gConsole.Red(filename)))
			}

		case strings.HasPrefix(cmd, "/visual"):
			if gVisual {
				gVisual = false
				fmt.Println("You are playing", gConsole.Bold(gConsole.Yellow("blind")), "now.")
			} else {
				gVisual = true
				fmt.Println("You are playing", gConsole.Bold(gConsole.Yellow("visual")), "now.")
				drawBoard(gGame)
			}
			continue

		case strings.HasPrefix(cmd, "/keys"):
			cmd := strings.SplitN(cmd, " ", 2)
			if len(cmd) > 1 {
				switch cmd[1] {
				case "vi":
					l.SetVimMode(true)
				case "emacs":
					l.SetVimMode(false)
				default:
					fmt.Println("Allowed arguments are", gConsole.Bold(gConsole.Yellow("[vi|emacs]")))
					continue
				}
			}

			if l.IsVimMode() {
				fmt.Println(gConsole.Bold(gConsole.Yellow("vi")), "key bindings active")
			} else {
				fmt.Println(gConsole.Bold(gConsole.Yellow("emacs")), "key bindings active")
			}

		case cmd == "/quit":

			// Save the game.
			if gameStarted && savePGN(gGame, gGameFilename) == nil { // Success
				fmt.Println("Game saved to", gConsole.Bold(gConsole.Red(gGameFilename)))
			}

			goto end

		default:
			// Send the human move to engine and get a counter move
			engineMoveNext(eng, gGame, cmd)
			gameStarted = true
			if isGameOver(gGame) {
				// Save the game.
				if savePGN(gGame, gGameFilename) == nil { // Success
					fmt.Println("Game saved to", gConsole.Bold(gConsole.Red(gGameFilename)))
				}
				goto end
			}
		}
	}
end:
}
