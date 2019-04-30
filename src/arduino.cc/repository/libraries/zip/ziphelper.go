package zip

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"arduino.cc/repository/libraries/file"
)

// ZipDirectory creates a new zip archive that contains a copy of "rootFolder" into "zipFile".
// Inside the archive "rootFolder" will be renamed to "zipRootFolderName".
func ZipDirectory(rootFolder string, zipRootFolderName string, zipFile string) error {
	checks := func(path string, info os.FileInfo, err error) error {
		info, err = os.Lstat(path)
		if err != nil {
			return err
		}
		if (info.Mode() & os.ModeSymlink) != 0 {
			dest, _ := os.Readlink(path)
			return fmt.Errorf("Symlink not allowed: %s -> %s", path, dest)
		}
		if file.IsSCCS(info.Name()) {
			return filepath.SkipDir
		}
		return nil
	}

	rootFolder, err := filepath.Abs(rootFolder)
	if err != nil {
		return err
	}

	if err := filepath.Walk(rootFolder, checks); err != nil {
		return err
	}

	tmpdir, err := ioutil.TempDir("", "ziphelper")
	if err != nil {
		return fmt.Errorf("creating temp dir for zip archive: %s", err)
	}
	defer os.RemoveAll(tmpdir)
	if err := os.Symlink(rootFolder, filepath.Join(tmpdir, zipRootFolderName)); err != nil {
		return fmt.Errorf("creating temp dir for zip archive: %s", err)
	}

	args := []string{"-r", zipFile, zipRootFolderName, "-x", ".*", "-x", "*/.*"}
	for sccs := range file.SCCSFiles {
		args = append(args, "-x")
		args = append(args, "*/"+sccs+"/*")
		args = append(args, "-x")
		args = append(args, sccs+"/*")
	}

	zipCmd := exec.Command("zip", args...)
	zipCmd.Dir = tmpdir
	if _, err := zipCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("archiving into zip: %s", err)
	}
	return nil
}
