package main

import (
	"encoding/json"
	"os"
	"arduino.cc/repository/libraries"
	"arduino.cc/repository/libraries/db"
	"github.com/robfig/cron"
	"log"
	"time"
)

type Config struct {
	BaseDownloadUrl string
	LibrariesFolder string
	LibrariesDB     string
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
	config = readConf()

	setup(config)

	syncLibraries()

	if config.CronTabEntry == "" {
		return
	}

	crontab := cron.New()
	crontab.AddFunc(config.CronTabEntry, func() { syncLibraries() })
	crontab.Start()

	for {
		time.Sleep(time.Hour)
	}
}

var running bool = false

func syncLibraries() {
	if running {
		log.Println("...still synchronizing...")
		return
	}

	running = true
	defer func() { running = false }()

	log.Println("Synchronizing libraries...")
	repos, err := libraries.ListRepos("./repos.txt")
	if logError(err) {
		os.Exit(1)
	}

	libraryDb := db.Init(config.LibrariesDB)

	for _, repo := range repos {
		log.Println("... " + repo)
		handleRepo(repo, libraryDb, config)
	}
	log.Println("...DONE")
}

func readConf() *Config {
	if _, err := os.Stat("./config.json"); os.IsNotExist(err) {
		logError(err)
		os.Exit(1)
	}

	file, _ := os.Open("./config.json")
	decoder := json.NewDecoder(file)
	config := Config{}
	err := decoder.Decode(&config)
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

	err = libraries.UpdateLibrary(library, libraryDb, config.BaseDownloadUrl)
	if logError(err) {
		return
	}

	err = libraries.ZipRepo(repo.Workdir(), library, config.LibrariesFolder)
	if logError(err) {
		return
	}

}
