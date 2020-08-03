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

	if playerColor() == chess.Black {
		err = engineMove(eng, gGame)
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
			gGame.Resign(playerColor())
			isGameOver(gGame)
			goto end

		case strings.HasPrefix(cmd, "/visual"):
			if gVisual {
				gVisual = false
				fmt.Println("You are playing", Bold(Yellow("blind")), "now.")
			} else {
				gVisual = true
				fmt.Println("You are playing", Bold(Yellow("visual")), "now.")
				fmt.Print(gGame.Position().Board().Draw())
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
					fmt.Println("Allowed arguments are", Bold(Yellow("[vi|emacs]")))
					continue
				}
			}

			if l.IsVimMode() {
				fmt.Println(Bold(Yellow("vi")), "key bindings active")
			} else {
				fmt.Println(Bold(Yellow("emacs")), "key bindings active")
			}
		case cmd == "/quit":
			goto end
		default:
			engineCounterMove(eng, gGame, cmd)
			if isGameOver(gGame) {
				goto end
			}
		}
	}
end:
}

func isGameOver(game *chess.Game) bool {
	switch game.Outcome() {
	case chess.NoOutcome:
		return false
	case chess.Draw:
		fmt.Println(Bold(Yellow("Game draw!!")))
	case chess.WhiteWon:
		fmt.Println(Bold(Yellow("White won the game!!")))
	case chess.BlackWon:
		fmt.Println(Bold(Yellow("Black won the game!!")))
	default:
		panic(game.Outcome()) // should never happen
	}
	return true // The end.
}

// Human's turn
func engineMove(engine *uci.Engine, game *chess.Game) error {
	results, err := engine.GoDepth(gEngineDepth, uci.HighestDepthOnly)
	if err != nil {
		fmt.Println(game.Position().Board().Draw())
		fmt.Println(err)
		return err
	}

	moveSAN, err := chess.LongAlgebraicNotation{}.Decode(game.Position(), results.BestMove)
	if err != nil {
		fmt.Println("SANe " + game.Position().Board().Draw())
		fmt.Println(err)
		return err
	}

	err = game.Move(moveSAN)
	if err != nil {
		fmt.Println(game.Position().Board().Draw())
		fmt.Println(err)
		return err
	}

	fmt.Println(enginePrompt() + chess.AlgebraicNotation{}.Encode(game.Position(), moveSAN))
	if gVisual {
		fmt.Print(game.Position().Board().Draw())
	}
	return nil
}

// Engine's turn
func engineCounterMove(engine *uci.Engine, game *chess.Game, moveStr string) error {
	err := game.MoveStr(moveStr)
	if err != nil {
		fmt.Println("Allowed moves:", Bold(Yellow(validMoves())))
		return err
	}
	engine.SetFEN(game.FEN())
	results, err := engine.GoDepth(gEngineDepth, uci.HighestDepthOnly)
	if err != nil {
		fmt.Println(game.Position().Board().Draw())
		fmt.Println(err)
		return err
	}

	moveSAN, err := chess.LongAlgebraicNotation{}.Decode(game.Position(), results.BestMove)
	if err != nil {
		fmt.Println("SANe " + game.Position().Board().Draw())
		fmt.Println(err)
		return err
	}

	err = game.Move(moveSAN)
	if err != nil {
		fmt.Println(game.Position().Board().Draw())
		fmt.Println(err)
		return err
	}

	fmt.Println(enginePrompt() + chess.AlgebraicNotation{}.Encode(game.Position(), moveSAN))
	if gVisual {
		fmt.Print(game.Position().Board().Draw())
	}
	return nil
}
