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
	"github.com/notnil/chess"
)

// Global defaults. Avoid global variables as much as possible.
var (
	gCfgFile      string
	gEngineBinary string
	gPlayerColor  string
	gVisual       bool

	gVersion     = "1.0-alpha"
	gGame        = chess.NewGame(chess.UseNotation(chess.AlgebraicNotation{}))
	gBlackPrompt = "(B) "
	gWhitePrompt = "(W) "
)