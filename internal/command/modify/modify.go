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

// Package modify implements the `modify` CLI subcommand used by the maintainer for modifications to the library registration data.
package modify

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
	"github.com/spf13/pflag"
)

var config *configuration.Config
var libraryName string
var libraryData *db.Library
var releasesData []*db.Release

// Run executes the command.
func Run(command *cobra.Command, cliArguments []string) {
	var err error
	config = configuration.ReadConf(command.Flags())

	libraryName = cliArguments[0]

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

	// Load all the library's data from the DB.
	librariesDb := db.Init(librariesDBPath.String())
	if !librariesDb.HasLibrary(libraryName) {
		feedback.Errorf("Library of name %s not found", libraryName)
		os.Exit(1)
	}
	libraryData, err = librariesDb.FindLibrary(libraryName)
	if err != nil {
		panic(err)
	}
	releasesData = librariesDb.FindReleasesOfLibrary(libraryData)

	restore, err := modifications(command.Flags())
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

func modifications(flags *pflag.FlagSet) (bool, error) {
	didModify := false // Require at least one modification operation was specified by user.

	newRepositoryURL, err := flags.GetString("repo-url")
	if err != nil {
		return false, err
	}
	newTypes, err := flags.GetString("types")
	if err != nil {
		return false, err
	}

	if newRepositoryURL != "" {
		if err := modifyRepositoryURL(newRepositoryURL); err != nil {
			return true, err
		}

		didModify = true
	}

	if newTypes != "" {
		if err := modifyTypes(newTypes); err != nil {
			return false, err
		}

		didModify = true
	}

	if !didModify {
		return false, fmt.Errorf("No modification flags provided so nothing happened. See 'libraries-repository-engine modify --help'")
	}

	return false, nil
}

func modifyRepositoryURL(newRepositoryURL string) error {
	if !libraries.RepoURLValid(newRepositoryURL) {
		return fmt.Errorf("Library URL %s does not have a valid format", newRepositoryURL)
	}

	if libraryData.Repository == newRepositoryURL {
		return fmt.Errorf("Library %s already has URL %s", libraryName, newRepositoryURL)
	}

	oldRepositoryURL := libraryData.Repository

	fmt.Printf("Changing URL of library %s from %s to %s\n", libraryName, oldRepositoryURL, newRepositoryURL)

	// Remove the library Git clone folder. It will be cloned from the new URL on the next sync.
	if err := libraries.BackupAndDeleteGitClone(config, &libraries.Repo{URL: libraryData.Repository}); err != nil {
		return err
	}

	// Update the library repository URL in the database.
	libraryData.Repository = newRepositoryURL

	// Update library releases.
	oldRepositoryObject := libraries.Repository{URL: oldRepositoryURL}
	newRepositoryObject := libraries.Repository{URL: newRepositoryURL}
	libraryMetadata := metadata.LibraryMetadata{Name: libraryData.Name}
	for _, releaseData := range releasesData {
		libraryMetadata.Version = releaseData.Version.String()
		oldArchiveObject, err := archive.New(&oldRepositoryObject, &libraryMetadata, config)
		if err != nil {
			return err
		}
		newArchiveObject, err := archive.New(&newRepositoryObject, &libraryMetadata, config)
		if err != nil {
			return err
		}

		// Move the release archive to the correct path for the new URL (some path components are based on the library repo URL).
		oldArchiveObjectPath := paths.New(oldArchiveObject.Path)
		newArchiveObjectPath := paths.New(newArchiveObject.Path)
		if err := newArchiveObjectPath.Parent().MkdirAll(); err != nil {
			return fmt.Errorf("While creating new library release archives path: %w", err)
		}
		if err := backup.Backup(oldArchiveObjectPath); err != nil {
			return fmt.Errorf("While backing up library release archive: %w", err)
		}
		if err := oldArchiveObjectPath.Rename(newArchiveObjectPath); err != nil {
			return fmt.Errorf("While moving library release archive: %w", err)
		}

		// Update the release download URL in the database.
		releaseData.URL = newArchiveObject.URL
	}

	return nil
}

func modifyTypes(rawTypes string) error {
	newTypes := strings.Split(rawTypes, ",")
	for i := range newTypes {
		newTypes[i] = strings.TrimSpace(newTypes[i])
	}

	sameTypes := func(oldTypes []string) bool {
		if len(oldTypes) != len(newTypes) {
			return false
		}

		for _, oldType := range oldTypes {
			found := false
			for _, newType := range newTypes {
				if oldType == newType {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		}

		return true
	}

	typesChanged := false

	for _, releaseData := range releasesData {
		if !typesChanged {
			// Compare old and new types for this release
			typesChanged = !sameTypes(releaseData.Types)
		}

		releaseData.Types = newTypes
	}

	if !typesChanged {
		return fmt.Errorf("Library %s already has types %s", libraryName, rawTypes)
	}

	return nil
}
