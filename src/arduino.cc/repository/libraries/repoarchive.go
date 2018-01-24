package libraries

import (
	"os"
	"path/filepath"
	"regexp"

	"arduino.cc/repository/libraries/metadata"
	"arduino.cc/repository/libraries/zip"
)

func ZipRepo(repoFolder string, baseFolder string, zipFolderName string) (string, error) {
	err := os.MkdirAll(baseFolder, os.FileMode(0755))
	if err != nil {
		return "", err
	}
	absoluteFileName := filepath.Join(baseFolder, zipFolderName+".zip")
	zipFile, err := os.Create(absoluteFileName)
	if err != nil {
		return "", err
	}
	err = zip.ZipDirectory(repoFolder, zipFolderName, zipFile)
	zipFile.Close()

	if err != nil {
		os.Remove(absoluteFileName)
		return "", err
	}

	return absoluteFileName, nil
}

func ZipFolderName(library *metadata.LibraryMetadata) string {
	pattern := regexp.MustCompile("[^a-zA-Z0-9]")
	return pattern.ReplaceAllString(library.Name, "_") + "-" + library.Version
}
