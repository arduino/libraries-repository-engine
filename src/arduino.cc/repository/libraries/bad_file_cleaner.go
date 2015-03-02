package libraries

import "os/exec"

var PATTERNS = []string{"*.exe"}

func RemoveUndesiderFiles(folder string) error {
	for _, pattern := range PATTERNS {
		cmd := exec.Command("find", folder, "-type", "f", "-name", pattern, "-delete")
		err := cmd.Run()
		if err != nil {
			return err
		}
	}
	return nil
}
