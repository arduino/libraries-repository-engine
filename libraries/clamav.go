// This file is part of libraries-repository-engine.
//
// Copyright 2021 ARDUINO SA (http://www.arduino.cc/)
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published
// by the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.
//
// You can be released from the requirements of the above licenses by purchasing
// a commercial license. Buying such a license is mandatory if you want to
// modify or otherwise use the software for commercial activities involving the
// Arduino software without disclosing the source code of your own applications.
// To purchase a commercial license, send an email to license@arduino.cc.

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
	cmd := exec.Command("clamdscan", "--fdpass", "-i", folder)
	cmd.Env = modifyEnv(os.Environ(), "LANG", "en")

	out, err := cmd.CombinedOutput()
	if err != nil {
		return out, err
	}

	output := string(out)
	if strings.Index(output, "Infected files: 0") == -1 {
		return out, errors.New("Infected files found")
	}

	return out, nil
}
