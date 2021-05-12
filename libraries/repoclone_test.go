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
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCloneRepos(t *testing.T) {
	meta := &Repo{URL: "https://github.com/arduino-libraries/Servo.git"}

	subfolder, err := meta.AsFolder()
	require.NoError(t, err)

	repo, err := CloneOrFetch(meta, filepath.Join("/tmp", subfolder))

	require.NoError(t, err)
	require.NotNil(t, repo)
	require.Equal(t, "/tmp/github.com/arduino-libraries/Servo", repo.FolderPath)

	defer os.RemoveAll(repo.FolderPath)

	_, err = os.Stat(repo.FolderPath)
	require.NoError(t, err)
}
