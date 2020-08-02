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
	"log"

	"github.com/freeeve/uci"
	"github.com/notnil/chess"
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func NewEngine(enginePath string) (*uci.Engine, error) {
	eng, err := uci.NewEngine(enginePath)
	if err != nil {
		log.Fatal(err)
	}

	return eng, err
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func playerColor() chess.Color {
	if gPlayerColor == "white" {
		return chess.White
	}
	return chess.Black
}

// Readline completion of all the valid moves left.
func validMovesConstructor(string) (moves []string) {
	for _, move := range gGame.Position().ValidMoves() {
		moveSAN := chess.Encoder.Encode(chess.AlgebraicNotation{}, gGame.Position(), move)
		moves = append(moves, moveSAN)
	}
	return moves
}
