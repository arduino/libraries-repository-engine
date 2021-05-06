package libraries

import (
	"errors"
	"os"
	"os/exec"
	"strings"
)

func envSliceToMap(env []string) map[string]string {
	envMap := make(map[string]string)
	for _, value := range env {
		key := value[:strings.Index(value, "=")]
		value = value[strings.Index(value, "=")+1:]
		envMap[key] = value
	}
	return envMap
}

func envMapToSlice(envMap map[string]string) []string {
	var env []string
	for key, value := range envMap {
		env = append(env, key+"="+value)
	}
	return env
}

func modifyEnv(env []string, key, value string) []string {
	envMap := envSliceToMap(env)
	envMap[key] = value
	return envMapToSlice(envMap)
}

// RunAntiVirus scans the folder for viruses.
func RunAntiVirus(folder string) ([]byte, error) {
	cmd := exec.Command("clamdscan", "-i", folder)
	cmd.Env = modifyEnv(os.Environ(), "LANG", "en")

	out, err := cmd.CombinedOutput()
	if err != nil {
		return out, err
	}

	output := string(out)
	if strings.Index(output, "Infected files: 0") == -1 {
		return out, errors.New("Infected files found!")
	}

	return out, nil
}
