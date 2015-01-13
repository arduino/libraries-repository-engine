package zip

import (
	"path/filepath"
	"os"
	"archive/zip"
	"io"
	"arduino.cc/repository/libraries/file"
)

func ZipDirectory(rootFolder string, zipFile *os.File) error {
	rootFolder, err := filepath.Abs(rootFolder)
	if err != nil {
		return err
	}

	zipFileWriter := zip.NewWriter(zipFile)
	defer zipFileWriter.Close()

	addEntryToZip := func(path string, info os.FileInfo, err error) error {
		rel, err := filepath.Rel(rootFolder, path)
		if err != nil {
			return err
		}

		if file.IsSCCS(info.Name()) {
			return filepath.SkipDir
		}
		if rel[0] == '.' || info.IsDir() {
			return nil
		}

		writer, err := zipFileWriter.Create(rel)
		if err != nil {
			return err
		}
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(writer, file)
		if err != nil {
			return err
		}

		return nil
	}

	err = filepath.Walk(rootFolder, addEntryToZip)
	if err != nil {
		return err
	}
	return nil
}
