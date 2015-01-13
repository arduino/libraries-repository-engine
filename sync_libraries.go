package main

import (
	"encoding/json"
	"os"
	"arduino.cc/repository/libraries"
	"arduino.cc/repository/libraries/db"
	"github.com/robfig/cron"
	"log"
	"time"
	"arduino.cc/repository/libraries/hash"
)

type Config struct {
	BaseDownloadUrl    string
	LibrariesFolder    string
	LibrariesDB        string
	LibrariesIndex     string
	GitClonesFolder    string
	CronTabEntry       string
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

	for _, repo := range repos {
		log.Println("... " + repo)
		handleRepo(repo, libraryDb, config)
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

	encoder := json.NewEncoder(file)
	err = encoder.Encode(libraryIndex)
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

func handleRepo(repoURL string, libraryDb *db.DB, config *Config) {
	repo, err := libraries.CloneOrFetch(repoURL, config.GitClonesFolder)
	if logError(err) {
		return
	}

	err = libraries.CheckoutLastTag(repo)
	if logError(err) {
		return
	}

	library, err := libraries.GenerateLibraryFromRepo(repo.Workdir())
	if logError(err) {
		return
	}

	err = libraries.ZipRepo(repo.Workdir(), library, config.LibrariesFolder)
	if logError(err) {
		return
	}

	release := db.FromLibraryToRelease(library, config.BaseDownloadUrl)

	err = setSizeAndChecksum(release, config.LibrariesFolder)
	if logError(err) {
		return
	}


	err = libraries.UpdateLibrary(release, libraryDb)
	if logError(err) {
		return
	}

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
