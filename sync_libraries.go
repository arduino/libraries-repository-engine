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
	log.Println("...DONE")
}

func readConf(configFile string) *Config {
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		logError(err)
		os.Exit(1)
	}

	file, _ := os.Open(configFile)
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
