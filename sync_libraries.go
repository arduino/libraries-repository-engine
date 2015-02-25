package main

import (
	"arduino.cc/repository/libraries"
	"arduino.cc/repository/libraries/db"
	"arduino.cc/repository/libraries/hash"
	"encoding/json"
	"github.com/robfig/cron"
	"log"
	"os"
	"time"
)

// TODO(cm): Merge this struct with config/config.go
type Config struct {
	BaseDownloadUrl string
	LibrariesFolder string
	LibrariesDB     string
	LibrariesIndex  string
	GitClonesFolder string
	CronTabEntry    string
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

	if config.CronTabEntry == "" {
		return
	}

	crontab := cron.New()
	crontab.AddFunc(config.CronTabEntry, func() { syncLibraries(reposFile) })
	crontab.Start()

	for {
		time.Sleep(time.Hour)
	}
}

var running bool = false

func syncLibraries(reposFile string) {
	if _, err := os.Stat(reposFile); os.IsNotExist(err) {
		logError(err)
		os.Exit(1)
	}

	if running {
		log.Println("...still synchronizing...")
		return
	}

	running = true
	defer func() { running = false }()

	log.Println("Synchronizing libraries...")
	repos, err := libraries.ListRepos(reposFile)
	if logError(err) {
		os.Exit(1)
	}

	libraryDb := db.Init(config.LibrariesDB)

	var errorWithARepo bool
	for _, repo := range repos {
		log.Println("... " + repo)
		err := handleRepo(repo, libraryDb, config)
		errorWithARepo = errorWithARepo || err != nil
	}

	libraryIndex, err := libraryDb.OutputLibraryIndex()
	if logError(err) {
		os.Exit(1)
	}

	serializeLibraryIndex(libraryIndex, config.LibrariesIndex)

	log.Println("...DONE")

	if errorWithARepo {
		os.Exit(1)
	}
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

func handleRepo(repoURL string, libraryDb *db.DB, config *Config) error {
	repoFolder, err := libraries.CloneOrFetch(repoURL, config.GitClonesFolder)
	if logError(err) {
		return err
	}

	err = libraries.CheckoutLastTag(repoFolder)
	if logError(err) {
		return err
	}

	library, err := libraries.GenerateLibraryFromRepo(repoFolder)
	if logError(err) {
		return err
	}

	zipFolderName := libraries.ZipFolderName(library)

	err = libraries.ZipRepo(repoFolder, config.LibrariesFolder, zipFolderName)
	if logError(err) {
		return err
	}

	release := db.FromLibraryToRelease(library, config.BaseDownloadUrl, zipFolderName+".zip")

	err = setSizeAndChecksum(release, config.LibrariesFolder)
	if logError(err) {
		return err
	}

	err = libraries.UpdateLibrary(release, libraryDb)
	if logError(err) {
		return err
	}

	return nil
}

func setSizeAndChecksum(release *db.Release, librariesFolder string) error {
	info, err := os.Stat(librariesFolder + release.ArchiveFileName)
	if err != nil {
		return err
	}

	release.Size = info.Size()
	checksum, err := hash.Checksum(librariesFolder + release.ArchiveFileName)
	if err != nil {
		return err
	}
	release.Checksum = checksum

	return nil
}
