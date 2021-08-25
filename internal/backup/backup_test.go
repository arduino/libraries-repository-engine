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

package backup

import (
	"testing"

	"github.com/arduino/go-paths-helper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testDataPath string

func TestAll(t *testing.T) {
	var err error
	originalsFolder, err := paths.MkTempDir("", "backup-test-testall")
	require.NoError(t, err)

	// Generate test content.
	originalContent := []byte("foo")
	modifyFile, err := paths.WriteToTempFile(originalContent, originalsFolder, "")
	require.NoError(t, err)
	modifyFolder, err := originalsFolder.MkTempDir("")
	require.NoError(t, err)
	modifyFolderFile, err := paths.WriteToTempFile(originalContent, modifyFolder, "")
	require.NoError(t, err)
	deleteFile, err := paths.WriteToTempFile(originalContent, originalsFolder, "")
	require.NoError(t, err)
	deleteFolder, err := originalsFolder.MkTempDir("")
	require.NoError(t, err)
	deleteFolderFile, err := paths.WriteToTempFile(originalContent, deleteFolder, "")
	require.NoError(t, err)

	// Backup test content.
	err = Backup(modifyFile)
	require.NoError(t, err)
	err = Backup(modifyFolder)
	require.NoError(t, err)
	err = Backup(deleteFile)
	require.NoError(t, err)
	err = Backup(deleteFolder)
	require.NoError(t, err)

	// Change the originals.
	err = modifyFile.WriteFile([]byte("bar"))
	require.NoError(t, err)
	err = modifyFolderFile.WriteFile([]byte("bar"))
	require.NoError(t, err)
	err = deleteFile.Remove()
	require.NoError(t, err)
	err = deleteFolder.RemoveAll()
	require.NoError(t, err)

	err = Restore()
	require.NoError(t, err)

	// Verify changes to originals were reverted.
	content, err := modifyFile.ReadFile()
	require.NoError(t, err)
	assert.Equal(t, originalContent, content)

	content, err = modifyFolderFile.ReadFile()
	require.NoError(t, err)
	assert.Equal(t, originalContent, content)

	assert.True(t, deleteFile.Exist())
	assert.True(t, deleteFolderFile.Exist())

	// Clean the backups.
	err = Clean()
	require.NoError(t, err)
}
