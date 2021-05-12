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

package metadata

import (
	"strings"

	"github.com/blang/semver"
)

// VersionToSemverCompliant transforms a truncated version to a semver compliant version,
// for example "1.0" is converted to "1.0.0".
func VersionToSemverCompliant(version string) (string, error) {
	versionParts := len(strings.Split(version, "."))
	for versionParts < 3 {
		versionParts++
		version += ".0"
	}

	newVersion, err := semver.Parse(version)
	if err != nil {
		return "", err
	}
	return newVersion.String(), nil
}
