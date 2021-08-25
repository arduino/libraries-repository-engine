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

// Package backup does backup and restore of files.
package backup

import (
	"github.com/arduino/go-paths-helper"
)

type backup struct {
	originalPath *paths.Path
	backupPath   *paths.Path
}

var backupsFolder *paths.Path
var backups []backup

// Backup saves a backup copy of the given path.
func Backup(originalPath *paths.Path) error {
	if backupsFolder == nil {
		// Create a parent folder to store all backups of this session.
		var err error
		if backupsFolder, err = paths.MkTempDir("", "libraries-repository-engine-backup"); err != nil {
			return err
		}
	}

	// Create a folder for this individual backup item.
	backupFolder, err := backupsFolder.MkTempDir("")
	if err != nil {
		return err
	}

	backupPath := backupFolder.Join(originalPath.Base())

	isDir, err := originalPath.IsDirCheck()
	if err != nil {
		return err
	}
	if isDir {
		if err := originalPath.CopyDirTo(backupPath); err != nil {
			return err
		}
	} else {
		if err := originalPath.CopyTo(backupPath); err != nil {
			return err
		}
	}

	backups = append(backups, backup{originalPath: originalPath, backupPath: backupPath})

	return nil
}

// Restore restores all backed up files.
func Restore() error {
	for _, backup := range backups {
		isDir, err := backup.backupPath.IsDirCheck()
		if err != nil {
			return err
		}
		if isDir {
			if err := backup.originalPath.RemoveAll(); err != nil {
				return err
			}
			if err := backup.backupPath.CopyDirTo(backup.originalPath); err != nil {
				return err
			}
		} else {
			if err := backup.backupPath.CopyTo(backup.originalPath); err != nil {
				return err
			}
		}
	}

	return nil
}

// Clean deletes all the backup files.
func Clean() error {
	if backupsFolder == nil {
		return nil
	}

	return backupsFolder.RemoveAll()
}
