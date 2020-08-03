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

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "pinata",
	Short:   "Piñata - play blindfold chess against UCI compatible engines.",
	Version: "1.0aplha",
	Long:    ``,

	// Transfer control to readline shell.
	Run: func(cmd *cobra.Command, args []string) {
		Shell()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// Load config file and register flags.
func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVarP(&gCfgFile, "config", "c", "pinata.toml", "config file")
	rootCmd.PersistentFlags().StringVarP(&gEngineBinary, "engine", "e", "stockfish", "path to UCI compatible chess engine executable")
	rootCmd.PersistentFlags().StringVarP(&gPlayerColor, "play", "p", "white", "choose black or white")
	rootCmd.PersistentFlags().BoolVarP(&gVisual, "visual", "v", false, "cheat blindfold")
	rootCmd.PersistentFlags().IntVarP(&gEngineDepth, "depth", "d", 10, "engine search depth")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if gCfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(gCfgFile)
	} else {
		/*
			// Find home directory.
			home, err := homedir.Dir()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		*/

		// Search config in home directory with name ".pinata" (without extension).
		// viper.AddConfigPath(home) // search first in the home dir.
		viper.AddConfigPath(".") // search in the current dir.

		// Default to TOML format, though viper supports a variety of standard formats.
		viper.SetConfigName(gCfgFile)
	}

	viper.AutomaticEnv() // read in environment variables that match
}
