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

// Package sync implements the `sync` CLI subcommand that updates the Library Manager content.
package sync

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	cc "github.com/arduino/golang-concurrent-workers"
	"github.com/arduino/libraries-repository-engine/internal/configuration"
	"github.com/arduino/libraries-repository-engine/internal/feedback"
	"github.com/arduino/libraries-repository-engine/internal/libraries"
	"github.com/arduino/libraries-repository-engine/internal/libraries/archive"
	"github.com/arduino/libraries-repository-engine/internal/libraries/db"
	"github.com/arduino/libraries-repository-engine/internal/libraries/gitutils"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/spf13/cobra"
)

var config *configuration.Config

// Run executes the command.
func Run(command *cobra.Command, cliArguments []string) {
	config = configuration.ReadConf(command.Flags())

	setup(config)

	var reposFile string
	if len(cliArguments) > 0 {
		reposFile = cliArguments[0]
	} else {
		reposFile = "./repos.txt"
	}

	if len(cliArguments) > 1 {
		feedback.LogError(fmt.Errorf("Multiple arguments are not supported"))
		os.Exit(1)
	}

	syncLibraries(reposFile)
}

func syncLibraries(reposFile string) {
	if _, err := os.Stat(reposFile); os.IsNotExist(err) {
		feedback.LogError(err)
		os.Exit(1)
	}

	log.Println("Synchronizing libraries...")
	repos, err := libraries.ListRepos(reposFile)
	if feedback.LogError(err) {
		os.Exit(1)
	}

	type jobContext struct {
		id           int
		repoMetadata *libraries.Repo
	}

	libraryDb := db.Init(config.LibrariesDB)

	jobQueue := make(chan *jobContext)

	pool := cc.New(4)
	worker := func() {
		log.Println("Started worker...")
		for job := range jobQueue {
			buffer := &bytes.Buffer{}
			logger := log.New(buffer, "", log.LstdFlags|log.LUTC)
			syncLibrary(logger, job.repoMetadata, libraryDb)

			// Output log to file
			if err := outputLogFile(logger, job.repoMetadata, buffer); err != nil {
				logger.Printf("Error writing log file: %s", err.Error())
			}

			// Output log to stdout
			fmt.Println(buffer.String())
		}
		log.Println("Completed worker!")
	}
	pool.Run(worker)
	pool.Run(worker)
	pool.Run(worker)
	pool.Run(worker)
	pool.Wait()

	go func() {
		id := 0
		for _, repo := range repos {
			jobQueue <- &jobContext{
				id:           id,
				repoMetadata: repo,
			}
			id++
		}
		close(jobQueue)
	}()

	for err := range pool.Errors {
		feedback.LogError(err)
	}

	libraryIndex, err := libraryDb.OutputLibraryIndex()
	if feedback.LogError(err) {
		os.Exit(1)
	}

	serializeLibraryIndex(libraryIndex, config.LibrariesIndex)

	log.Println("...DONE")
}

func serializeLibraryIndex(libraryIndex interface{}, libraryIndexFile string) {
	file, err := os.Create(libraryIndexFile)
	if feedback.LogError(err) {
		os.Exit(1)
	}
	defer file.Close()

	b, err := json.MarshalIndent(libraryIndex, "", "  ")
	if feedback.LogError(err) {
		os.Exit(1)
	}

	_, err = file.Write(b)
	if feedback.LogError(err) {
		os.Exit(1)
	}
}

func setup(config *configuration.Config) {
	err := os.MkdirAll(config.GitClonesFolder, os.FileMode(0777))
	if feedback.LogError(err) {
		os.Exit(1)
	}
	err = os.MkdirAll(config.LibrariesFolder, os.FileMode(0777))
	if feedback.LogError(err) {
		os.Exit(1)
	}
}

func syncLibrary(logger *log.Logger, repoMetadata *libraries.Repo, libraryDb *db.DB) {
	logger.Printf("Scraping %s", repoMetadata.URL)

	repoFolderName, err := repoMetadata.AsFolder()
	if err != nil {
		logger.Printf("Invalid URL: %s", err.Error())
		return
	}
	repoFolder := filepath.Join(config.GitClonesFolder, repoFolderName)

	// Clone repository
	repo, err := libraries.CloneOrFetch(repoMetadata, repoFolder)
	if err != nil {
		logger.Printf("Error fetching repository: %s", err)
		logger.Printf("Removing clone and trying again")
		os.RemoveAll(repoFolder)
		repo, err = libraries.CloneOrFetch(repoMetadata, repoFolder)
		if err != nil {
			logger.Printf("Error fetching repository: %s", err)
			logger.Printf("Leaving...")
			return
		}
	}

	// Retrieve the list of git-tags
	tags, err := gitutils.SortedCommitTags(repo.Repository)
	if err != nil {
		logger.Printf("Error retrieving git-tags: %s", err)
		return
	}

	for _, tag := range tags {
		// Sync the library release for each git-tag
		err = syncLibraryTaggedRelease(logger, repo, tag, repoMetadata, libraryDb)
		if err != nil {
			logger.Printf("Error syncing library: %s", err)
		}
	}
}

func syncLibraryTaggedRelease(logger *log.Logger, repo *libraries.Repository, tag *plumbing.Reference, repoMeta *libraries.Repo, libraryDb *db.DB) error {
	var releaseLog string // This string will be displayed in the logs for indexed releases.

	// Checkout desired tag
	logger.Printf("Checking out tag: %s", tag.Name().Short())
	if err := gitutils.CheckoutTag(repo.Repository, tag); err != nil {
		return fmt.Errorf("Error checking out repo: %s", err)
	}

	// Create library metadata from library.properties
	library, err := libraries.GenerateLibraryFromRepo(repo)
	if err != nil {
		return fmt.Errorf("Error generating library from repo: %s", err)
	}
	library.Types = repoMeta.Types

	// If the release name is different from the listed name, skip release...
	if library.Name != repoMeta.LibraryName {
		logger.Printf("Release %s:%s has wrong library name, should be %s", library.Name, library.Version, repoMeta.LibraryName)
		return nil
	}

	releaseQuery := db.Release{
		LibraryName: library.Name,
		Version:     db.VersionFromString(library.Version),
	}
	// If the release is already checked in, skip
	if libraryDb.HasLibrary(library.Name) {
		if release, _ := libraryDb.FindRelease(&releaseQuery); release != nil {
			logger.Printf("Release %s:%s already loaded, skipping", library.Name, library.Version)
			if release.Log != "" {
				logger.Print(release.Log)
			}
			return nil
		}
	}

	if !config.DoNotRunClamav {
		if out, err := libraries.RunAntiVirus(repo.FolderPath); err != nil {
			logger.Printf("clamav output:\n%s", out)
			return err
		}
	}

	report, err := libraries.RunArduinoLint(config.ArduinoLintPath, repo.FolderPath, repoMeta)
	reportTemplate := `<a href="https://arduino.github.io/arduino-lint/latest/">Arduino Lint</a> %s:
<details><summary>Click to expand Arduino Lint report</summary>
<hr>
%s
<hr>
</details>`
	if err != nil {
		logger.Printf(reportTemplate, "found errors", report)
		return err
	}
	if report != nil {
		formattedReport := fmt.Sprintf(reportTemplate, "has suggestions for possible improvements", report)
		logger.Print(formattedReport)
		releaseLog += formattedReport
	}

	zipName := archive.ZipFolderName(library)
	lib := filepath.Base(filepath.Clean(filepath.Join(repo.FolderPath, "..")))
	host := filepath.Base(filepath.Clean(filepath.Join(repo.FolderPath, "..", "..")))
	zipFilePath, err := archive.ZipRepo(repo.FolderPath, filepath.Join(config.LibrariesFolder, host, lib), zipName)
	if err != nil {
		return fmt.Errorf("Error while zipping library: %s", err)
	}

	size, checksum, err := archive.GetSizeAndCalculateChecksum(zipFilePath)
	if err != nil {
		return fmt.Errorf("Error while calculating checksums: %s", err)
	}
	release := db.FromLibraryToRelease(library)
	release.URL = config.BaseDownloadURL + host + "/" + lib + "/" + zipName + ".zip"
	release.ArchiveFileName = zipName + ".zip"
	release.Size = size
	release.Checksum = checksum
	release.Log = releaseLog

	if err := libraries.UpdateLibrary(release, repo.URL, libraryDb); err != nil {
		return fmt.Errorf("Error while updating library DB: %s", err)
	}

	return nil
}

func outputLogFile(logger *log.Logger, repoMetadata *libraries.Repo, buffer *bytes.Buffer) error {
	if config.LogsFolder == "" {
		return nil
	}
	repoSubFolder, err := repoMetadata.AsFolder()
	if err != nil {
		return fmt.Errorf("URL Path: %s", err.Error())
	}
	logFolder := filepath.Join(config.LogsFolder, repoSubFolder)
	if _, err = os.Stat(logFolder); os.IsNotExist(err) {
		err = os.MkdirAll(logFolder, os.FileMode(0755))
	}
	if err != nil {
		return fmt.Errorf("mkdir %s: %s", logFolder, err.Error())
	}
	logFile := filepath.Join(logFolder, "index.html")
	output := "<pre>\n" + buffer.String() + "\n</pre>"
	if err := ioutil.WriteFile(logFile, []byte(output), 0644); err != nil {
		return fmt.Errorf("write log to file: %s", err.Error())
	}
	return nil
}
