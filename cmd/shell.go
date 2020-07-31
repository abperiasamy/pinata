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
	"github.com/notnil/chess"
)

var completer = readline.NewPrefixCompleter(
	readline.PcItemDynamic(ValidMovesConstructor),
	readline.PcItem("/keys",
		readline.PcItem("vi"),
		readline.PcItem("emacs"),
	),

	readline.PcItem("/quit"),
	readline.PcItem("/new"),
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
	eng, err := NewEngine("/usr/games/stockfish")
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
		Prompt: "\033[31m»\033[0m ",
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
			goto exit
		default:
			err := game.MoveStr(line)
			if err != nil {
				fmt.Println(game.Position().Board().Draw())
				fmt.Println(err)
				continue
			}
			eng.SetFEN(game.FEN())
			results, err := eng.GoDepth(10, uci.HighestDepthOnly)
			if err != nil {
				fmt.Println(game.Position().Board().Draw())
				fmt.Println(err)
				continue
			}

			moveSAN, err := chess.LongAlgebraicNotation{}.Decode(game.Position(), results.BestMove)
			if err != nil {
				fmt.Println(game.Position().Board().Draw())
				fmt.Println(err)
				continue
			}

			fmt.Println("SAN: " + chess.AlgebraicNotation{}.Encode(game.Position(), moveSAN))

			err = game.Move(moveSAN)
			if err != nil {
				fmt.Println(game.Position().Board().Draw())
				fmt.Println(err)
				continue
			}
			fmt.Println(game.Position().Board().Draw())
		}
	}
exit:
}
