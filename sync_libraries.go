package main

import (
	"encoding/json"
	"os"
	"fmt"
	"arduino.cc/repository/libraries"
	"arduino.cc/repository/libraries/db"
)

type Config struct {
	BaseDownloadUrl string
	LibrariesFolder string
	LibrariesDB     string
	GitClonesFolder string
}

func logError(err error) bool {
	if err != nil {
		fmt.Println(err)
		return true
	}
	return false
}

func main() {
	config := readConf()

	setup(config)

	repos, err := libraries.ListRepos("./repos.txt")
	if logError(err) {
		os.Exit(1)
	}

	libraryDb := db.Init(config.LibrariesDB)

	for _, repo := range repos {
		handleRepo(repo, libraryDb, config)
	}

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
