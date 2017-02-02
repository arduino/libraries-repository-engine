package main

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"fmt"

	"arduino.cc/repository/libraries"
	"arduino.cc/repository/libraries/db"
	"arduino.cc/repository/libraries/hash"
	"github.com/arduino/arduino-modules/git"
	"github.com/arduino/golang-concurrent-workers"
)

// TODO(cm): Merge this struct with config/config.go
type Config struct {
	BaseDownloadUrl string
	LibrariesFolder string
	LibrariesDB     string
	LibrariesIndex  string
	GitClonesFolder string
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
			logger := log.New(os.Stdout, fmt.Sprintf("JOB %03d - ", job.id), log.LstdFlags)
			logger.Printf("Scraping %s", job.repoMetadata.Url)

			// Clone repository
			repo, err := libraries.CloneOrFetch(job.repoMetadata.Url, config.GitClonesFolder)
			if err != nil {
				logger.Printf("Error fetching repository: %s", err)
				continue
			}

			// Retrieve the list of git-tags
			tags, err := repo.ListTags()
			if err != nil {
				logger.Printf("Error retrieving git-tags: %s", err)
				continue
			}

			for _, tag := range tags {
				// Sync the library release for each git-tag
				err = syncLibraryTaggedRelease(logger, repo, tag, job.repoMetadata, libraryDb)
				if err != nil {
					logger.Printf("Error syncing library: %s", err)
				}
			}
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

// TODO(cm): Merge this struct with config/config.go
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

func syncLibraryTaggedRelease(logger *log.Logger, repo *git.Repository, tag string, repoMeta *libraries.Repo, libraryDb *db.DB) error {
	logger.Printf("Checking out tag: %s", tag)
	if out, err := repo.CheckoutTagWithOutput(tag); err != nil {
		logger.Printf("git output: %s", out)
		return fmt.Errorf("Error checking out repo: %s", err)
	}

	library, err := libraries.GenerateLibraryFromRepo(repo)
	if err != nil {
		return fmt.Errorf("Error generating library from repo: %s", err)
	}
	library.Types = repoMeta.Types
	library.Name = repoMeta.LibraryName

	if libraryDb.HasLibrary(library.Name) && libraryDb.HasReleaseByNameVersion(library.Name, library.Version) {
		logger.Printf("Release %s:%s already loaded, skipping", library.Name, library.Version)
		return nil
	}

	if err := libraries.FailIfHasUndesiredFiles(repo.FolderPath); err != nil {
		return err
	}

	if out, err := libraries.RunAntiVirus(repo.FolderPath); err != nil {
		logger.Printf("clamav output:\n%s", out)
		return err
	}

	zipFolderName := libraries.ZipFolderName(library)
	libFolder := filepath.Base(filepath.Clean(filepath.Join(repo.FolderPath, "..")))
	zipFilePath, err := libraries.ZipRepo(repo.FolderPath, filepath.Join(config.LibrariesFolder, libFolder), zipFolderName)
	if err != nil {
		return fmt.Errorf("Error while zipping library: %s", err)
	}

	size, checksum, err := getSizeAndCalculateChecksum(zipFilePath)
	if err != nil {
		return fmt.Errorf("Error while calculating checksums: %s", err)
	}
	release := db.FromLibraryToRelease(library)
	release.URL = config.BaseDownloadUrl + libFolder + "/" + zipFolderName + ".zip"
	release.ArchiveFileName = zipFolderName + ".zip"
	release.Size = size
	release.Checksum = checksum

	/*
		if strings.Index(repo.Url, "https://github.com") != -1 {
			url, size, checksum, err := libraries.GithubDownloadRelease(repo.Url, tag)
			if logError(err) {
				return err
			}
			release.URL = url
			release.Size = size
			release.Checksum = checksum
			release.ArchiveFileName = zipFolderName + "-github.zip"
		}
	*/
	if err := libraries.UpdateLibrary(release, libraryDb); err != nil {
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
