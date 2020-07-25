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

	"github.com/dolegi/uci"
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func InitUCI(enginePath string) (*uci.Engine, error) {
	eng, err := uci.NewEngine(enginePath)
	if err != nil {
		log.Fatal(err)
	}

	return eng, err

	/*
		eng, err := uci.NewEngine("/usr/games/stockfish")
		if err != nil {
			log.Fatal(err)
		}

		// set some engine options
		eng.SetOptions(uci.Options{
			Hash:    128,
			Ponder:  false,
			OwnBook: true,
			MultiPV: 4,
		})

		// set the starting position
		eng.SetFEN("rnb4r/ppp1k1pp/3bp3/1N3p2/1P2n3/P3BN2/2P1PPPP/R3KB1R b KQ - 4 11")

		// set some result filter options
		//resultOpts := uci.HighestDepthOnly | uci.IncludeUpperbounds | uci.IncludeLowerbounds
		// results, _ := eng.GoDepth(10, resultOpts)

		return eng, err
	*/
}
