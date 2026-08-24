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

package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
	semver "go.bug.st/relaxed-semver"
)

func releasesFromVersions(versions ...string) []*Release {
	releases := make([]*Release, len(versions))
	for i, version := range versions {
		releases[i] = &Release{
			Version:  semver.MustParse(version),
			Size:     1,
			Checksum: "SHA-256:0",
		}
	}
	return releases
}

func versionsOf(releases []*Release) []string {
	versions := make([]string, len(releases))
	for i, release := range releases {
		versions[i] = release.Version.String()
	}
	return versions
}

func TestSelectReleasesWithMajorVersionCoverage(t *testing.T) {
	releases := releasesFromVersions(
		"1.0.0", "1.1.0", "1.2.0",
		"2.0.0", "2.1.0",
		"3.0.0",
	)

	// One slot to spare after the first round covers every major: the second most
	// recent release of the next major in round-robin order (2.x) fills it.
	selected := selectReleasesWithMajorVersionCoverage(releases, 4)
	assert.Equal(t, []string{"3.0.0", "2.1.0", "2.0.0", "1.2.0"}, versionsOf(selected))

	// Exactly enough slots for one round: one release per major version.
	selected = selectReleasesWithMajorVersionCoverage(releases, 3)
	assert.Equal(t, []string{"3.0.0", "2.1.0", "1.2.0"}, versionsOf(selected))

	// Fewer slots than major versions: only the first majors in round-robin order fit.
	selected = selectReleasesWithMajorVersionCoverage(releases, 2)
	assert.Equal(t, []string{"3.0.0", "2.1.0"}, versionsOf(selected))

	// More slots than releases: everything is kept.
	selected = selectReleasesWithMajorVersionCoverage(releases, 100)
	assert.Equal(t, []string{"3.0.0", "2.1.0", "2.0.0", "1.2.0", "1.1.0", "1.0.0"}, versionsOf(selected))
}

func TestSelectReleasesWithMajorVersionCoverageRoundRobinFairness(t *testing.T) {
	// Major 3 has more releases than majors 2 and 1: a naive "fill remaining slots with
	// the most recent releases" approach would exhaust the budget on 3.x's third release
	// (3.7.0) before giving 2.x a second pick. Round-robin instead gives every major a
	// turn each round, so 2.8.0 is picked before 3.7.0.
	releases := releasesFromVersions(
		"3.9.0", "3.8.0", "3.7.0",
		"2.9.0", "2.8.0",
		"1.0.0",
	)

	selected := selectReleasesWithMajorVersionCoverage(releases, 5)
	assert.Equal(t, []string{"3.9.0", "3.8.0", "2.9.0", "2.8.0", "1.0.0"}, versionsOf(selected))
}

func TestOutputLibraryIndexMaxVersionsPerLibraryKeepsOneReleasePerMajor(t *testing.T) {
	testDB := DB{
		Libraries: []*Library{
			{Name: "FooLib", Repository: "https://github.com/Bar/FooLib.git"},
		},
		Releases: releasesFromVersions("1.0.0", "1.1.0", "2.0.0", "3.0.0", "3.1.0"),
	}
	for _, release := range testDB.Releases {
		release.LibraryName = "FooLib"
	}

	index, err := testDB.OutputLibraryIndex(2)
	assert.NoError(t, err)

	output := index.(*indexOutput)
	versions := make([]string, len(output.Libraries))
	for i, lib := range output.Libraries {
		versions[i] = lib.Version.String()
	}
	assert.Equal(t, []string{"3.1.0", "2.0.0"}, versions)
}
