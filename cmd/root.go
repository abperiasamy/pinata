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
	"runtime"

	"github.com/abperiasamy/chess"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "pinata",
	Short:   "Piñata - play blindfold chess against UCI compatible engines.",
	Version: gVersion,
	Long:    ``,

	// Transfer control to readline shell.
	Run: func(cmd *cobra.Command, args []string) {
		onStart() // Perform post initialization
		shell()   // Shell controls the game interaction from start to finish.
		onStop()  // Perform cleanup
	},
}

// Perform post initialization routines right before starting the game.
func onStart() {
	initGlobals()

	// Invert colors on a brighter background
	if gLightBg {
		chess.ConsoleDark = false
	} else {
		chess.ConsoleDark = true // also prepare chess package for dark background
	}
	// Also disable chess package color printing.
	chess.ConsoleColor = !gNoColor
	chess.ConsoleUnicode = !gNoColor // also disable unicode printing

}

// Perform post initialization routines right after the game ends.
func onStop() {

}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// Load config file and register flags.
func init() {
	// fmt.Print("\033[?25l") // Hide cursor
	// Initialize config and ENV first. Command-line flags may override these settings.
	// cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVarP(&gCfgFile, "config", "c", "pinata.toml", "config file")
	rootCmd.PersistentFlags().StringVarP(&gEngineBinary, "engine", "e", "stockfish", "path to UCI compatible chess engine executable")
	rootCmd.PersistentFlags().StringVarP(&gGamePath, "file", "f", "", "load game from a PGN file")
	rootCmd.PersistentFlags().BoolVarP(&gHumanIsBlack, "black", "b", false, "choose the black side")
	rootCmd.PersistentFlags().BoolVarP(&gVisual, "visual", "v", false, "cheat blindfold")
	if runtime.GOOS == "windows" { // disable color and unicode support on Windows by default
		rootCmd.PersistentFlags().BoolVar(&gNoColor, "color", true, "disable colors")
	} else {
		rootCmd.PersistentFlags().BoolVar(&gNoColor, "no-color", false, "disable colors")
	}
	rootCmd.PersistentFlags().BoolVarP(&gLightBg, "light", "l", false, "invert the colors for lighter console background")
	rootCmd.PersistentFlags().IntVarP(&gEngineDepth, "depth", "d", 10, "engine search depth")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
}

/*
// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if gCfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(gCfgFile)
	} else {
			// Find home directory.
			home, err := homedir.Dir()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

		// Search config in home directory with name ".pinata" (without extension).
		// viper.AddConfigPath(home) // search first in the home dir.
		viper.AddConfigPath(".") // search in the current dir.

		// Default to TOML format, though viper supports a variety of standard formats.
		viper.SetConfigName(gCfgFile)
	}

	viper.AutomaticEnv() // read in environment variables that match
}
*/
