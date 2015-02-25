package libraries

import (
	"arduino.cc/repository/libraries/metadata"
	"arduino.cc/repository/libraries/zip"
	"os"
	"path"
	"regexp"
)

func ZipRepo(repoFolder string, librariesBaseFolder string, zipFolderName string) error {
	absoluteFileName := path.Join(librariesBaseFolder, zipFolderName+".zip")
	zipFile, err := os.Create(absoluteFileName)
	if err != nil {
		return err
	}

	defer zipFile.Close()

	return zip.ZipDirectory(repoFolder, zipFolderName, zipFile)
}

func ZipFolderName(library *metadata.LibraryMetadata) string {
	pattern := regexp.MustCompile("[^a-zA-Z0-9]")
	return pattern.ReplaceAllString(library.Name, "_") + "-" + library.Version
}
