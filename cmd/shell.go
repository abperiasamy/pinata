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
	"os"
	"strings"

	"github.com/abperiasamy/chess"
	"github.com/chzyer/readline"
	"github.com/freeeve/uci"
	. "github.com/logrusorgru/aurora"
)

var completer = readline.NewPrefixCompleter(
	readline.PcItemDynamic(validMovesConstructor),
	readline.PcItem("resign"),
	readline.PcItem("/new"),
	readline.PcItem("/visual"),
	readline.PcItem("/quit"),

	readline.PcItem("/keys",
		readline.PcItem("vi"),
		readline.PcItem("emacs"),
	),
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
func Shell() {
	eng, err := NewEngine(gEngineBinary)
	if err != nil {
		log.Fatal(err)
	}
	defer eng.Close()

	eng.SendOption("Threads", "8")

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

	if gPlayerColor == "black" {
		results, err := eng.GoDepth(10, uci.HighestDepthOnly)
		if err != nil {
			fmt.Println(gGame.Position().Board().Draw())
			fmt.Println(err)
			os.Exit(1)
		}

		moveSAN, err := chess.LongAlgebraicNotation{}.Decode(gGame.Position(), results.BestMove)
		if err != nil {
			fmt.Println("SANe " + gGame.Position().Board().Draw())
			fmt.Println(err)
			os.Exit(1)
		}

		err = gGame.Move(moveSAN)
		if err != nil {
			fmt.Println(gGame.Position().Board().Draw())
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println(enginePrompt() + chess.AlgebraicNotation{}.Encode(gGame.Position(), moveSAN))
		if gVisual {
			fmt.Print(gGame.Position().Board().Draw())
		}
	}

	l.SetPrompt(playerPrompt())

	for {

		line, err := l.Readline()
		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				break
			} else {
				continue
			}
		} else if err == io.EOF {
			break
		}

		line = strings.TrimSpace(line)
		switch {
		case line == "":
		case line == "resign":
			gGame.Resign(chess.Black)
			// fmt.Println(Outcome().String())
			goto end

		case strings.HasPrefix(line, "/visual"):
			if gVisual {
				gVisual = false
				fmt.Println("You are playing", Bold(Yellow("blind")), "now.")
			} else {
				gVisual = true
				fmt.Print(gGame.Position().Board().Draw())
			}
			continue

		case strings.HasPrefix(line, "/keys"):
			cmd := strings.SplitN(line, " ", 2)
			if len(cmd) > 1 {
				switch cmd[1] {
				case "vi":
					l.SetVimMode(true)
				case "emacs":
					l.SetVimMode(false)
				default:
					fmt.Println("Allowed arguments are", Bold(Yellow("[vi|emacs]")))
					continue
				}
			}

			if l.IsVimMode() {
				fmt.Println(Bold(Yellow("vi")), "key bindings active")
			} else {
				fmt.Println(Bold(Yellow("emacs")), "key bindings active")
			}
		case line == "/quit":
			goto end
		default:
			err := gGame.MoveStr(line)
			if err != nil {
				fmt.Println("Allowed moves:", Bold(Yellow(validMoves())))
				continue
			}
			eng.SetFEN(gGame.FEN())
			results, err := eng.GoDepth(10, uci.HighestDepthOnly)
			if err != nil {
				fmt.Println(gGame.Position().Board().Draw())
				fmt.Println(err)
				continue
			}

			moveSAN, err := chess.LongAlgebraicNotation{}.Decode(gGame.Position(), results.BestMove)
			if err != nil {
				fmt.Println("SANe " + gGame.Position().Board().Draw())
				fmt.Println(err)
				continue
			}

			err = gGame.Move(moveSAN)
			if err != nil {
				fmt.Println(gGame.Position().Board().Draw())
				fmt.Println(err)
				continue
			}

			fmt.Println(playerPrompt() + chess.AlgebraicNotation{}.Encode(gGame.Position(), moveSAN))
			if gVisual {
				fmt.Print(gGame.Position().Board().Draw())
			}
		}
	}
end:
}
