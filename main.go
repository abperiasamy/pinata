// +build go1.13

/*
 * blindfold (C) 2020 Anand Babu Periasamy
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"fmt"
	uci "github.com/abperiasamy/blindfold/pkg/uci"
	"log"
)

func main() {
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
	resultOpts := uci.HighestDepthOnly | uci.IncludeUpperbounds | uci.IncludeLowerbounds
	results, _ := eng.GoDepth(10, resultOpts)

	// print it (String() goes to pretty JSON for now)
	fmt.Println(results)
}
