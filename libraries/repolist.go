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

package libraries

import (
	"github.com/arduino/libraries-repository-engine/internal/libraries"
)

// LoadRepoListFromFile returns an unfiltered list of library registry entries loaded from the given data file.
func LoadRepoListFromFile(filename string) ([]*Repo, error) {
	return libraries.LoadRepoListFromFile(filename)
}

// Repo is the type for the library repository data.
type Repo = libraries.Repo

// ListRepos returns a filtered list of library registry entries loaded from the given data file.
func ListRepos(reposFilename string) ([]*Repo, error) {
	return libraries.ListRepos(reposFilename)
}
