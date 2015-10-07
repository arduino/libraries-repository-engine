package libraries

import (
	"errors"
	"os"
	"os/exec"
	"path"
)

func FailIfHasUndesiredFiles(folder string) error {
	err := failIfContainsForbiddenFileInRoot(folder)
	if err != nil {
		return err
	}
	return failIfContainsExes(folder)
}

var FORBIDDEN_FILES = []string{".development"}

func failIfContainsForbiddenFileInRoot(folder string) error {
	for _, file := range FORBIDDEN_FILES {
		if _, err := os.Stat(path.Join(folder, file)); err == nil {
			return errors.New("... ... " + file + " file found, skipping")
		}
	}

	return nil
}

var PATTERNS = []string{"*.exe"}

func failIfContainsExes(folder string) error {
	for _, pattern := range PATTERNS {
		cmd := exec.Command("find", folder, "-type", "f", "-name", pattern)
		output, err := cmd.CombinedOutput()
		if err != nil {
			return err
		}
		if len(string(output)) > 0 {
			return errors.New("... ... " + pattern + " files found, skipping")
		}
	}
	return nil
}
