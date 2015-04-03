package libraries

import (
	"os/exec"
	"errors"
)

var PATTERNS = []string{"*.exe"}

func FailIfHasUndesiredFiles(folder string) error {
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
