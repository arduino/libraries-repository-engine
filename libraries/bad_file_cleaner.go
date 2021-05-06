package libraries

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
)

// FailIfHasUndesiredFiles returns an error if the folder contains any undesired files.
func FailIfHasUndesiredFiles(folder string) error {
	err := failIfContainsForbiddenFileInRoot(folder)
	if err != nil {
		return err
	}
	return failIfContainsExes(folder)
}

// ForbiddenFiles is the names of the forbidden files.
var ForbiddenFiles = []string{".development"}

func failIfContainsForbiddenFileInRoot(folder string) error {
	for _, file := range ForbiddenFiles {
		if _, err := os.Stat(filepath.Join(folder, file)); err == nil {
			return errors.New(file + " file found, skipping")
		}
	}

	return nil
}

// Patterns is the file patterns of executables.
var Patterns = []string{"*.exe"}

func failIfContainsExes(folder string) error {
	for _, pattern := range Patterns {
		cmd := exec.Command("find", folder, "-type", "f", "-name", pattern)
		output, err := cmd.CombinedOutput()
		if err != nil {
			return err
		}
		if len(string(output)) > 0 {
			return errors.New(pattern + " files found, skipping")
		}
	}
	return nil
}
