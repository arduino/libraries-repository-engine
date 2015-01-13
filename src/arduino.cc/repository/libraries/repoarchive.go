package libraries

import (
	"arduino.cc/repository/libraries/metadata"
	"arduino.cc/repository/libraries/zip"
	"os"
	"path"
)

func ZipRepo(repoFolder string, library *metadata.LibraryMetadata, librariesBaseFolder string) error {
	absoluteFileName := path.Join(librariesBaseFolder, library.Name+"-"+library.Version+".zip")
	zipFile, err := os.Create(absoluteFileName)
	if err != nil {
		return err
	}

	defer zipFile.Close()

	return zip.ZipDirectory(repoFolder, zipFile)
}
