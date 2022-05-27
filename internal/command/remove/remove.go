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

// Package remove implements the `remove` CLI subcommand used by the maintainer for removals of libraries or releases.
package remove

import (
	"fmt"
	"os"
	"strings"

	"github.com/arduino/go-paths-helper"
	"github.com/arduino/libraries-repository-engine/internal/backup"
	"github.com/arduino/libraries-repository-engine/internal/configuration"
	"github.com/arduino/libraries-repository-engine/internal/feedback"
	"github.com/arduino/libraries-repository-engine/internal/libraries"
	"github.com/arduino/libraries-repository-engine/internal/libraries/archive"
	"github.com/arduino/libraries-repository-engine/internal/libraries/db"
	"github.com/arduino/libraries-repository-engine/internal/libraries/metadata"

	"github.com/spf13/cobra"
)

var config *configuration.Config
var librariesDb *db.DB
var libraryData *db.Library

// Run executes the command.
func Run(command *cobra.Command, cliArguments []string) {
	config = configuration.ReadConf(command.Flags())

	if len(cliArguments) == 0 {
		feedback.Error("LIBRARY_NAME argument is required")
		os.Exit(1)
	}

	librariesDBPath := paths.New(config.LibrariesDB)
	exist, err := librariesDBPath.ExistCheck()
	if err != nil {
		feedback.Errorf("While checking existence of database file: %s", err)
		os.Exit(1)
	}
	if !exist {
		feedback.Errorf("Database file not found at %s. Check the LibrariesDB configuration value.", librariesDBPath)
		os.Exit(1)
	}

	if err := backup.Backup(librariesDBPath); err != nil {
		feedback.Errorf("While backing up database: %s", err)
		os.Exit(1)
	}

	librariesDb = db.Init(config.LibrariesDB)

	restore, err := removals(cliArguments)
	if err != nil {
		feedback.Error(err)
		if restore {
			if err := backup.Restore(); err != nil {
				feedback.Errorf("While restoring the content from backup: %s", err)
			}
			fmt.Println("Original files were restored.")
		} else {
			if err := backup.Clean(); err != nil {
				feedback.Errorf("While cleaning up the backup content: %s", err)
			}
		}
		os.Exit(1)
	}

	if err := librariesDb.Commit(); err != nil {
		feedback.Errorf("While saving changes to database: %s", err)
		if err := backup.Restore(); err != nil {
			feedback.Errorf("While restoring the content from backup: %s", err)
		}
		fmt.Println("Original files were restored.")
		os.Exit(1)
	}

	if err := backup.Clean(); err != nil {
		feedback.Errorf("While cleaning up the backup files: %s", err)
		os.Exit(1)
	}

	fmt.Println("Success!")
}

func removals(libraryReferences []string) (bool, error) {
	for _, libraryReference := range libraryReferences {
		referenceComponents := strings.SplitN(libraryReference, "@", 2)
		libraryName := referenceComponents[0]
		var libraryVersion string
		if len(referenceComponents) > 1 {
			if referenceComponents[1] == "" {
				return false, fmt.Errorf("Missing version for library name %s. For full removal, omit the '@'", libraryName)
			}
			libraryVersion = referenceComponents[1]
		}

		if !librariesDb.HasLibrary(libraryName) {
			return false, fmt.Errorf("Library name %s not found", libraryName)
		}

		var err error
		libraryData, err = librariesDb.FindLibrary(libraryName)
		if err != nil {
			return true, err
		}

		if libraryVersion == "" {
			// Remove the library entirely.
			if err := removeLibrary(libraryName); err != nil {
				return true, err
			}
		} else {
			// Remove only a specific release of the library.
			if err := removeRelease(libraryName, libraryVersion); err != nil {
				return true, err
			}
		}
	}

	return false, nil
}

func removeLibrary(libraryName string) error {
	fmt.Printf("Removing %s\n", libraryName)

	// Remove the library's release archive files.
	releasesData := librariesDb.FindReleasesOfLibrary(libraryData)
	for _, releaseData := range releasesData {
		if err := removeReleaseArchive(releaseData.Version.String()); err != nil {
			return err
		}
	}

	// Remove the library and its releases from database.
	if err := librariesDb.RemoveLibrary(libraryName); err != nil {
		return err
	}

	// Remove the library Git clone folder.
	if err := libraries.BackupAndDeleteGitClone(config, &libraries.Repo{URL: libraryData.Repository}); err != nil {
		return err
	}

	return nil
}

func removeRelease(libraryName string, version string) error {
	fmt.Printf("Removing %s@%s\n", libraryName, version)

	if !librariesDb.HasReleaseByNameVersion(libraryName, version) {
		return fmt.Errorf("Library release %s@%s not found", libraryName, version)
	}

	// Remove the release archive file.
	if err := removeReleaseArchive(version); err != nil {
		return err
	}

	// Remove the release from the database.
	if err := librariesDb.RemoveReleaseByNameVersion(libraryName, version); err != nil {
		return err
	}

	return nil
}

func removeReleaseArchive(version string) error {
	repositoryObject := libraries.Repository{URL: libraryData.Repository}
	libraryMetadata := metadata.LibraryMetadata{
		Name:    libraryData.Name,
		Version: version,
	}
	archiveObject, err := archive.New(&repositoryObject, &libraryMetadata, config)
	if err != nil {
		panic(err)
	}

	archivePath := paths.New(archiveObject.Path)
	if err := backup.Backup(archivePath); err != nil {
		return fmt.Errorf("While backing up library release archive: %w", err)
	}
	if err := archivePath.RemoveAll(); err != nil {
		return fmt.Errorf("While removing library release archive: %s", err)
	}

	return nil
}
