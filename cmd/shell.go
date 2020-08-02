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

	"github.com/chzyer/readline"
	"github.com/freeeve/uci"
	. "github.com/logrusorgru/aurora"
	"github.com/notnil/chess"
)

var completer = readline.NewPrefixCompleter(
	readline.PcItemDynamic(validMovesConstructor),
	readline.PcItem("resign"),
	readline.PcItem("/quit"),
	readline.PcItem("/new"),

	readline.PcItem("/keys",
		readline.PcItem("vi"),
		readline.PcItem("emacs"),
	),
	readline.PcItem("/blind",
		readline.PcItem("on"),
		readline.PcItem("off"),
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

	/*
		eng.SetOption("Ponder", false)
		eng.SetOption("Threads", "8")

			eng.NewGame(uci.NewGameOpts{
			InitialFen: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR",
			// Side:       uci.Black,
		})
	*/

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

	blindMode := true

	l.SetPrompt(gWhitePrompt)
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

		case strings.HasPrefix(line, "/blind"):
			cmd := strings.SplitN(line, " ", 2)
			if len(cmd) > 1 {
				switch cmd[1] {
				case "on":
					blindMode = true
				case "off":
					blindMode = false
				}
			}

			if blindMode {
				fmt.Println("You are playing", Bold(Yellow("BLIND")))
			} else {
				fmt.Println("You are playing", Bold(Yellow("VISUAL")))
			}

		case strings.HasPrefix(line, "/keys "):
			mode := line[6:]
			switch mode {
			case "vi":
				l.SetVimMode(true)
			case "emacs":
				l.SetVimMode(false)
			default:
				println("invalid mode:", mode)
			}
		case line == "/keys":
			if l.IsVimMode() {
				println("current mode: vim")
			} else {
				println("current mode: emacs")
			}
		case line == "/quit":
			goto end
		default:

			err := gGame.MoveStr(line)
			if err != nil {
				fmt.Println("invalid move: " + line) //
				//fmt.Println("Allowed Moves: "+gGame.ValidMoves())
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

			fmt.Println(gBlackPrompt + chess.AlgebraicNotation{}.Encode(gGame.Position(), moveSAN))

			err = gGame.Move(moveSAN)
			if err != nil {
				fmt.Println(gGame.Position().Board().Draw())
				fmt.Println(err)
				continue
			}
			if !blindMode {
				fmt.Println(gGame.Position().Board().Draw())
			}
		}
	}
end:
}
