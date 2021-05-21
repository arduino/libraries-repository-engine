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

package gitutils

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

// ResolveTag returns the commit hash associated with a tag.
func ResolveTag(tag *plumbing.Reference, repository *git.Repository) (*plumbing.Hash, error) {
	// Annotated tags have their own hash, different from the commit hash, so the tag must be resolved to get the has for
	// the associated commit.
	// Tags may point to any Git object. Although not common, this can include tree and blob objects in addition to commits.
	// Resolving non-commit objects results in an error.
	return repository.ResolveRevision(plumbing.Revision(tag.Hash().String()))
}
