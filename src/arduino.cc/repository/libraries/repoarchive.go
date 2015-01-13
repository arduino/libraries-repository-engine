package libraries

import (
	"arduino.cc/repository/libraries/metadata"
	"arduino.cc/repository/libraries/zip"
	"os"
	"path"
)

func ZipRepo(repoFolder string, library *metadata.LibraryMetadata, librariesBaseFolder string) error {
	libraryName := library.Name + "-" + library.Version
	absoluteFileName := path.Join(librariesBaseFolder, libraryName+".zip")
	zipFile, err := os.Create(absoluteFileName)
	if err != nil {
		return err
	}

	defer zipFile.Close()

	return zip.ZipDirectory(repoFolder, libraryName, zipFile)
}
