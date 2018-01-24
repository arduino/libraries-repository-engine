package zip

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"arduino.cc/repository/libraries/file"
)

// Create a new zip archive that contains a copy of "rootFolder" into "zipFile".
// Inside the archive "rootFolder" will be renamed to "zipRootFolder".
func ZipDirectory(rootFolder string, zipRootFolderName string, zipFile *os.File) error {
	rootFolder, err := filepath.Abs(rootFolder)
	if err != nil {
		return err
	}

	zipFileWriter := zip.NewWriter(zipFile)
	defer zipFileWriter.Close()

	addEntryToZip := func(path string, info os.FileInfo, err error) error {
		info, err = os.Lstat(path)
		if err != nil {
			return err
		}
		if (info.Mode() & os.ModeSymlink) != 0 {
			dest, _ := os.Readlink(path)
			return fmt.Errorf("Symlink not allowed: %s -> %s", path, dest)
		}
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

		rel = filepath.Join(zipRootFolderName, rel)
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
