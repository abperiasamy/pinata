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
	"io"
	"log"
	"strings"

	"github.com/abperiasamy/chess"
	"github.com/chzyer/readline"
)

func filterInput(r rune) (rune, bool) {
	switch r {
	// block CtrlZ feature
	case readline.CharCtrlZ:
		return r, false
	}
	return r, true
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func shell() {
	game := NewGame()

	eng, err := NewEngine(gEngineBinary)
	if err != nil {
		log.Fatal(err)
	}
	defer eng.Close()
	eng.SendOption("Threads", "8")

	completer := readline.NewPrefixCompleter(
		readline.PcItemDynamic(validMovesConstructor(game)),
		readline.PcItem("resign"),
		readline.PcItem("/visual"),
		readline.PcItem("/quit"),

		readline.PcItem("/keys",
			readline.PcItem("vi"),
			readline.PcItem("emacs"),
		))

	l, err := readline.NewEx(&readline.Config{
		// Prompt: "\033[31m»\033[0m ",
		// HistoryFile:     "/tmp/readline.tmp",
		AutoComplete:    completer,
		InterruptPrompt: "^C",
		EOFPrompt:       "/quit",

		HistorySearchFold:   true,
		FuncFilterInputRune: filterInput,
	})
	if err != nil {
		panic(err)
	}
	defer l.Close()

	if playerColor() == chess.Black {
		err = engineMove(eng, game)
		if err != nil {
			fmt.Errorf("Engine Failure: ", err)
		}
	}

	l.SetPrompt(playerPrompt())

	for {
		cmd, err := l.Readline()
		if err == readline.ErrInterrupt {
			if len(cmd) == 0 {
				break
			} else {
				continue
			}
		} else if err == io.EOF {
			break
		}

		cmd = strings.TrimSpace(cmd)
		switch {
		case cmd == "": // no input, do nothing.

		case cmd == "resign":
			game.Resign(playerColor())
			isGameOver(game)
			goto end

		case strings.HasPrefix(cmd, "/visual"):
			if gVisual {
				gVisual = false
				fmt.Println("You are playing", gConsole.Bold(gConsole.Yellow("blind")), "now.")
			} else {
				gVisual = true
				fmt.Println("You are playing", gConsole.Bold(gConsole.Yellow("visual")), "now.")
				fmt.Print(game.Position().Board().Draw())
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
			goto end
		default:
			engineCounterMove(eng, game, cmd)
			if isGameOver(game) {
				goto end
			}
		}
	}
end:
}
