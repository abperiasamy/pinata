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
	"os"
	"os/exec"

	"github.com/abperiasamy/chess"
	"github.com/freeeve/uci"
)

// The shell initializes the engine upon entry.
func newEngine(enginePath string) (*uci.Engine, error) {
	_, err := exec.LookPath(gEngineBinary)
	if err != nil { // Alternatively look under games dir.
		path, err := exec.LookPath("/usr/games/" + gEngineBinary)
		if err != nil {
			fmt.Println("Unable to find " + gConsole.Bold(gConsole.Red(gEngineBinary)).String() + ". Please use `--engine` flag to choose a UCI compatible engine.")
			os.Exit(1)
		}
		gEngineBinary = path
	}

	eng, err := uci.NewEngine(enginePath)
	if err != nil {
		fmt.Println(gConsole.Red(err))
		fmt.Println("Unable to initialize " + gConsole.Bold(gConsole.Red(gEngineBinary)).String() + ". Please use `--engine` flag to choose a UCI compatible engine.")
		os.Exit(1)
	}

	return eng, err
}

// Engine's first move as white
func engineMoveFirst(engine *uci.Engine, game *chess.Game) error {
	results, err := engine.GoDepth(gEngineDepth, uci.HighestDepthOnly)
	if err != nil {
		fmt.Println(err)
		return err
	}

	moveLAN, err := chess.LongAlgebraicNotation{}.Decode(game.Position(), results.BestMove)
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = game.Move(moveLAN)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println(enginePrompt() + chess.Encoder.Encode(chess.AlgebraicNotation{}, game.Position(), moveLAN))
	drawBoard(game)
	return nil
}

// Engine's turn
func engineCounterMove(engine *uci.Engine, game *chess.Game, moveStr string) error {
	err := game.MoveStr(moveStr)
	if err != nil {
		fmt.Println("Allowed moves:", gConsole.Bold(gConsole.Yellow(validMoves(game))))
		return err
	}
	engine.SetFEN(game.FEN())
	results, err := engine.GoDepth(gEngineDepth, uci.HighestDepthOnly)
	if err != nil {
		fmt.Println(err)
		return err
	}

	moveLAN, err := chess.LongAlgebraicNotation{}.Decode(game.Position(), results.BestMove)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// Only the valid moves list has the equivalent SAN move with tag pairs.
	for _, move := range game.Position().ValidMoves() {
		if moveLAN.String() == chess.Encoder.Encode(chess.LongAlgebraicNotation{}, game.Position(), move) {
			fmt.Println(enginePrompt() + chess.Encoder.Encode(chess.AlgebraicNotation{}, game.Position(), move))
		}
	}

	err = game.Move(moveLAN)
	if err != nil {
		fmt.Println(err)
		return err
	}

	drawBoard(game)
	return nil
}
