// This file is part of libraries-repository-engine.
//
// Copyright 2021 ARDUINO SA (http://www.arduino.cc/)
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published
// by the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.
//
// You can be released from the requirements of the above licenses by purchasing
// a commercial license. Buying such a license is mandatory if you want to
// modify or otherwise use the software for commercial activities involving the
// Arduino software without disclosing the source code of your own applications.
// To purchase a commercial license, send an email to license@arduino.cc.

// Package cli defines the command line interface.
package cli

import (
	"os"

	"github.com/arduino/libraries-repository-engine/internal/feedback"
	"github.com/spf13/cobra"
)

// rootCmd defines the base CLI command.
var rootCmd = &cobra.Command{
	Use:  "libraries-repository-engine",
	Long: "The tool for managing the Arduino Library Manager content.",
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
	Aliases: []string{"sync"},
}

func init() {
	rootCmd.PersistentFlags().String("config-file", "config.json", "Configuration file path")
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	// Support the old `libraries-repository-engine [CONFIG_FILE [REGISTRY_FILE]]` interface.
	func() {
		if len(os.Args) > 1 {
			if os.Args[1][0] == '-' {
				// The argument is a flag, so assume the new interface.
				return
			}

			for _, command := range rootCmd.Commands() {
				for _, alias := range append(command.Aliases, command.Name(), "help") { // Hacky to check "help" redundantly, but whatever.
					if os.Args[1] == alias {
						// The argument is a registered subcommand, so assume the new interface.
						return
					}
				}
			}

			// The argument is not a registered subcommand, so assume it is the CONFIG_FILE positional argument of the old interface.
			rootCmd.PersistentFlags().Set("config-file", os.Args[1]) // Transfer the argument to the new interface's flag.
			os.Args = append([]string{os.Args[0]}, os.Args[2:]...)   // Remove the argument.
		}

		feedback.Warning(
			`Using deprecated command line syntax. New syntax:
libraries-repository-engine sync [--config-file=CONFIG_FILE] [REGISTRY_FILE]
`,
		)
		// Set the subcommand to the root's alias.
		os.Args = append([]string{os.Args[0], rootCmd.Aliases[0]}, os.Args[1:]...)
	}()

	return rootCmd.Execute()
}
