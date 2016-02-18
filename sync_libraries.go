package main

import (
	"arduino.cc/repository/libraries"
	"arduino.cc/repository/libraries/db"
	"arduino.cc/repository/libraries/hash"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	//"strings"
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

	libraryDb := db.Init(config.LibrariesDB)

	var errorWithARepo bool
	for _, repo := range repos {
		log.Println("... " + repo.Url)
		errors := syncLibraryInRepo(repo, libraryDb, config)
		errorWithARepo = errorWithARepo || (errors != nil && len(errors) > 0)
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

func syncLibraryInRepo(repo *libraries.Repo, libraryDb *db.DB, config *Config) []error {
	repoFolder, err := libraries.CloneOrFetch(repo.Url, config.GitClonesFolder)
	if logError(err) {
		return []error{err}
	}

	tags, err := libraries.ListTags(repoFolder)
	if logError(err) {
		return []error{err}
	}

	var errors []error
	for _, tag := range tags {
		err = syncLibraryTaggedRelease(repoFolder, tag, repo, libraryDb, config)
		if err != nil {
			errors = append(errors, err)
		}
	}

	return errors
}

func syncLibraryTaggedRelease(repoFolder string, tag string, repo *libraries.Repo, libraryDb *db.DB, config *Config) error {
	log.Println("... ... tag " + tag)

	err := libraries.CheckoutTag(repoFolder, tag)
	if logError(err) {
		return err
	}

	library, err := libraries.GenerateLibraryFromRepo(repoFolder, repo)
	if logError(err) {
		return err
	}

	if libraryDb.HasLibrary(library.Name) && libraryDb.HasReleaseByNameVersion(library.Name, library.Version) {
		log.Println("... ... tag " + tag + " already loaded: skipping")
		return nil
	}

	err = libraries.FailIfHasUndesiredFiles(repoFolder)
	if logError(err) {
		return err
	}

	err = libraries.RunAntiVirus(repoFolder)
	if logError(err) {
		return err
	}

	zipFolderName := libraries.ZipFolderName(library)
	libFolder := filepath.Base(filepath.Clean(filepath.Join(repoFolder, "..")))
	zipFilePath, err := libraries.ZipRepo(repoFolder, filepath.Join(config.LibrariesFolder, libFolder), zipFolderName)
	if logError(err) {
		return err
	}

	size, checksum, err := getSizeAndCalculateChecksum(zipFilePath)
	if logError(err) {
		return err
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
	err = libraries.UpdateLibrary(release, libraryDb)
	if logError(err) {
		return err
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
