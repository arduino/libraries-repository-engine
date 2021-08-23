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

package cli

import (
	"github.com/arduino/libraries-repository-engine/internal/command/remove"
	"github.com/spf13/cobra"
)

// removeCmd defines the `remove` CLI subcommand.
var removeCmd = &cobra.Command{
	Short:                 "Remove libraries or releases",
	Long:                  "Remove libraries or library releases from Library Manager",
	DisableFlagsInUseLine: true,
	Use: `remove [FLAG]... LIBRARY_NAME[@RELEASE]...

Remove library name LIBRARY_NAME Library Manager content entirely.
-or-
Remove release RELEASE of library name LIBRARY_NAME from the Library Manager content.`,
	Run: remove.Run,
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
