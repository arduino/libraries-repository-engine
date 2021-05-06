package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"arduino.cc/repository/libraries"
	"arduino.cc/repository/libraries/db"
	"arduino.cc/repository/libraries/hash"
	cc "github.com/arduino/golang-concurrent-workers"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

type Config struct {
	BaseDownloadURL string
	LibrariesFolder string
	LogsFolder      string
	LibrariesDB     string
	LibrariesIndex  string
	GitClonesFolder string
	DoNotRunClamav  bool
}

func logError(err error) bool {
	if err != nil {
		log.Println(err)
		return true
	}
	return false
}

var config *Config

func main() {
	var configFile string
	if len(os.Args) > 1 {
		configFile = os.Args[1]
	} else {
		configFile = "./config.json"
	}

	config = readConf(configFile)

	setup(config)

	var reposFile string
	if len(os.Args) > 2 {
		reposFile = os.Args[2]
	} else {
		reposFile = "./repos.txt"
	}

	syncLibraries(reposFile)
}

func syncLibraries(reposFile string) {
	if _, err := os.Stat(reposFile); os.IsNotExist(err) {
		logError(err)
		os.Exit(1)
	}

	log.Println("Synchronizing libraries...")
	repos, err := libraries.ListRepos(reposFile)
	if logError(err) {
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
			logger := log.New(buffer, fmt.Sprintf("JOB %03d - ", job.id), log.LstdFlags|log.LUTC)
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
		logError(err)
	}

	libraryIndex, err := libraryDb.OutputLibraryIndex()
	if logError(err) {
		os.Exit(1)
	}

	serializeLibraryIndex(libraryIndex, config.LibrariesIndex)

	log.Println("...DONE")
}

func serializeLibraryIndex(libraryIndex interface{}, libraryIndexFile string) {
	file, err := os.Create(libraryIndexFile)
	if logError(err) {
		os.Exit(1)
	}
	defer file.Close()

	b, err := json.MarshalIndent(libraryIndex, "", "  ")
	if logError(err) {
		os.Exit(1)
	}

	_, err = file.Write(b)
	if logError(err) {
		os.Exit(1)
	}
}

func readConf(configFile string) *Config {
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		logError(err)
		os.Exit(1)
	}

	file, err := os.Open(configFile)
	if logError(err) {
		os.Exit(1)
	}
	decoder := json.NewDecoder(file)
	config := Config{}
	err = decoder.Decode(&config)
	if logError(err) {
		os.Exit(1)
	}
	return &config
}

func setup(config *Config) {
	err := os.MkdirAll(config.GitClonesFolder, os.FileMode(0777))
	if logError(err) {
		os.Exit(1)
	}
	err = os.MkdirAll(config.LibrariesFolder, os.FileMode(0777))
	if logError(err) {
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
	tags, err := repo.Repository.Tags()
	if err != nil {
		logger.Printf("Error retrieving git-tags: %s", err)
		return
	}

	for {
		tag, err := tags.Next()
		if err != nil {
			// Reached end of tags
			break
		}

		// Sync the library release for each git-tag
		err = syncLibraryTaggedRelease(logger, repo, tag, repoMetadata, libraryDb)
		if err != nil {
			logger.Printf("Error syncing library: %s", err)
		}
	}
}

func syncLibraryTaggedRelease(logger *log.Logger, repo *libraries.Repository, tag *plumbing.Reference, repoMeta *libraries.Repo, libraryDb *db.DB) error {
	// Checkout desired tag
	logger.Printf("Checking out tag: %s", tag.Name())

	repoTree, err := repo.Repository.Worktree()
	if err != nil {
		return err
	}

	if err = repoTree.Checkout(&git.CheckoutOptions{Hash: tag.Hash()}); err != nil {
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

	// If the release is already checked in, skip
	if libraryDb.HasLibrary(library.Name) && libraryDb.HasReleaseByNameVersion(library.Name, library.Version) {
		logger.Printf("Release %s:%s already loaded, skipping", library.Name, library.Version)
		return nil
	}

	// Check if the library has undesried files
	if err := libraries.FailIfHasUndesiredFiles(repo.FolderPath); err != nil {
		return err
	}

	if !config.DoNotRunClamav {
		if out, err := libraries.RunAntiVirus(repo.FolderPath); err != nil {
			logger.Printf("clamav output:\n%s", out)
			return err
		}
	}

	zipName := libraries.ZipFolderName(library)
	lib := filepath.Base(filepath.Clean(filepath.Join(repo.FolderPath, "..")))
	host := filepath.Base(filepath.Clean(filepath.Join(repo.FolderPath, "..", "..")))
	logger.Printf("path=%s/%s, archivename=%s", host, lib, zipName)
	zipFilePath, err := libraries.ZipRepo(repo.FolderPath, filepath.Join(config.LibrariesFolder, host, lib), zipName)
	if err != nil {
		return fmt.Errorf("Error while zipping library: %s", err)
	}

	size, checksum, err := getSizeAndCalculateChecksum(zipFilePath)
	if err != nil {
		return fmt.Errorf("Error while calculating checksums: %s", err)
	}
	release := db.FromLibraryToRelease(library)
	release.URL = config.BaseDownloadURL + host + "/" + lib + "/" + zipName + ".zip"
	release.ArchiveFileName = zipName + ".zip"
	release.Size = size
	release.Checksum = checksum

	if err := libraries.UpdateLibrary(release, repo.URL, libraryDb); err != nil {
		return fmt.Errorf("Error while updating library DB: %s", err)
	}

	return nil
}

func getSizeAndCalculateChecksum(filePath string) (int64, string, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return -1, "", err
	}

	size := info.Size()

	checksum, err := hash.Checksum(filePath)
	if err != nil {
		return -1, "", err
	}

	return size, checksum, nil
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
