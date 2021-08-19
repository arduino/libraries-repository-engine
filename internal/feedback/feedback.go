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

// Package feedback provides feedback to the user.
package feedback

import (
	"fmt"
	"log"
	"os"
)

// LogError logs non-nil errors and returns whether the error was nil.
func LogError(err error) bool {
	if err != nil {
		log.Println(err)
		return true
	}
	return false
}

// Warningf behaves like fmt.Printf but adds a prefix and newline.
func Warningf(format string, v ...interface{}) {
	Warning(fmt.Sprintf(format, v...))
}

// Warning behaves like fmt.Println but adds a prefix.
func Warning(v ...interface{}) {
	fmt.Fprint(os.Stderr, "warning: ")
	fmt.Fprintln(os.Stderr, v...)
}

// Errorf behaves like fmt.Printf but adds a prefix and newline.
func Errorf(format string, v ...interface{}) {
	Error(fmt.Sprintf(format, v...))
}

// Error behaves like fmt.Println but adds a prefix.
func Error(v ...interface{}) {
	fmt.Fprint(os.Stderr, "error: ")
	fmt.Fprintln(os.Stderr, v...)
}
