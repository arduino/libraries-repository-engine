// This file is part of libraries-repository-engine.
//
// Copyright 2025 ARDUINO SA (http://www.arduino.cc/)
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

package checkregistry

import (
	"fmt"
	"os"
	"reflect"

	"github.com/arduino/libraries-repository-engine/internal/libraries"
)

// CheckRegistry runs the check-registry action
func CheckRegistry(reposFile string) {
	info, err := os.Stat(reposFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: While loading registry data file: %s\n", err)
		os.Exit(1)
	}

	if info.IsDir() {
		fmt.Fprintf(os.Stderr, "error: Registry data file argument %s is a folder, not a file\n", reposFile)
		os.Exit(1)
	}

	rawRepos, err := libraries.LoadRepoListFromFile(reposFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: While loading registry data file: %s\n", err)
		os.Exit(1)
	}

	filteredRepos, err := libraries.ListRepos(reposFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: While filtering registry data file: %s\n", err)
		os.Exit(1)
	}

	if !reflect.DeepEqual(rawRepos, filteredRepos) {
		fmt.Fprintln(os.Stderr, "error: Registry data file contains duplicate URLs")
		os.Exit(1)
	}

	validTypes := map[string]bool{
		"Arduino":     true,
		"Contributed": true,
		"Partner":     true,
		"Recommended": true,
		"Retired":     true,
	}

	nameMap := make(map[string]bool)
	for _, entry := range rawRepos {
		// Check entry types
		if len(entry.Types) == 0 {
			fmt.Fprintf(os.Stderr, "error: Type not specified for library '%s'\n", entry.LibraryName)
			os.Exit(1)
		}
		for _, entryType := range entry.Types {
			if _, valid := validTypes[entryType]; !valid {
				fmt.Fprintf(os.Stderr, "error: Invalid type '%s' used by library '%s'\n", entryType, entry.LibraryName)
				os.Exit(1)
			}
		}

		// Check library name of the entry
		if _, found := nameMap[entry.LibraryName]; found {
			fmt.Fprintf(os.Stderr, "error: Registry data file contains duplicates of name '%s'\n", entry.LibraryName)
			os.Exit(1)
		}
		nameMap[entry.LibraryName] = true
	}
}
