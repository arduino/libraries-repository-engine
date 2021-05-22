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

package metadata

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	testTables := []struct {
		testName                 string
		propertiesData           []byte
		libraryMetadataAssertion *LibraryMetadata
		errorAssertion           assert.ErrorAssertionFunc
	}{
		{
			testName:                 "Invalid",
			propertiesData:           []byte(`broken`),
			libraryMetadataAssertion: nil,
			errorAssertion:           assert.Error,
		},
		{
			testName: "Compliant",
			propertiesData: []byte(`
name=WebServer
version=1.0.0
author=Cristian Maglie <c.maglie@example.com>, Pippo Pluto <pippo@example.com>
maintainer=Cristian Maglie <c.maglie@example.com>
sentence=A library that makes coding a Webserver a breeze.
paragraph=Supports HTTP1.1 and you can do GET and POST.
category=Communication
url=http://example.com/
architectures=avr
includes=WebServer.h
depends=ArduinoHttpClient
			`),
			libraryMetadataAssertion: &LibraryMetadata{
				Name:          "WebServer",
				Version:       "1.0.0",
				Author:        "Cristian Maglie <c.maglie@example.com>, Pippo Pluto <pippo@example.com>",
				Maintainer:    "Cristian Maglie <c.maglie@example.com>",
				License:       "",
				Sentence:      "A library that makes coding a Webserver a breeze.",
				Paragraph:     "Supports HTTP1.1 and you can do GET and POST.",
				URL:           "http://example.com/",
				Architectures: "avr",
				Category:      "Communication",
				Types:         nil,
				Includes:      "WebServer.h",
				Depends:       "ArduinoHttpClient",
			},
			errorAssertion: assert.NoError,
		},
		{
			testName: "Invalid version",
			propertiesData: []byte(`
name=WebServer
version=foo
author=Cristian Maglie <c.maglie@example.com>, Pippo Pluto <pippo@example.com>
maintainer=Cristian Maglie <c.maglie@example.com>
sentence=A library that makes coding a Webserver a breeze.
paragraph=Supports HTTP1.1 and you can do GET and POST.
category=Communication
url=http://example.com/
architectures=avr
includes=WebServer.h
depends=ArduinoHttpClient
			`),
			libraryMetadataAssertion: &LibraryMetadata{
				Name:          "WebServer",
				Version:       "foo",
				Author:        "Cristian Maglie <c.maglie@example.com>, Pippo Pluto <pippo@example.com>",
				Maintainer:    "Cristian Maglie <c.maglie@example.com>",
				License:       "",
				Sentence:      "A library that makes coding a Webserver a breeze.",
				Paragraph:     "Supports HTTP1.1 and you can do GET and POST.",
				URL:           "http://example.com/",
				Architectures: "avr",
				Category:      "Communication",
				Types:         nil,
				Includes:      "WebServer.h",
				Depends:       "ArduinoHttpClient",
			},
			errorAssertion: assert.NoError,
		},
		{
			testName: "Non-semver version",
			propertiesData: []byte(`
name=WebServer
version=1.0
author=Cristian Maglie <c.maglie@example.com>, Pippo Pluto <pippo@example.com>
maintainer=Cristian Maglie <c.maglie@example.com>
sentence=A library that makes coding a Webserver a breeze.
paragraph=Supports HTTP1.1 and you can do GET and POST.
category=Communication
url=http://example.com/
architectures=avr
includes=WebServer.h
depends=ArduinoHttpClient
			`),
			libraryMetadataAssertion: &LibraryMetadata{
				Name:          "WebServer",
				Version:       "1.0.0",
				Author:        "Cristian Maglie <c.maglie@example.com>, Pippo Pluto <pippo@example.com>",
				Maintainer:    "Cristian Maglie <c.maglie@example.com>",
				License:       "",
				Sentence:      "A library that makes coding a Webserver a breeze.",
				Paragraph:     "Supports HTTP1.1 and you can do GET and POST.",
				URL:           "http://example.com/",
				Architectures: "avr",
				Category:      "Communication",
				Types:         nil,
				Includes:      "WebServer.h",
				Depends:       "ArduinoHttpClient",
			},
			errorAssertion: assert.NoError,
		},
		{
			testName: "Invalid category",
			propertiesData: []byte(`
name=WebServer
version=1.0.0
author=Cristian Maglie <c.maglie@example.com>, Pippo Pluto <pippo@example.com>
maintainer=Cristian Maglie <c.maglie@example.com>
sentence=A library that makes coding a Webserver a breeze.
paragraph=Supports HTTP1.1 and you can do GET and POST.
category=foo
url=http://example.com/
architectures=avr
includes=WebServer.h
depends=ArduinoHttpClient
			`),
			libraryMetadataAssertion: &LibraryMetadata{
				Name:          "WebServer",
				Version:       "1.0.0",
				Author:        "Cristian Maglie <c.maglie@example.com>, Pippo Pluto <pippo@example.com>",
				Maintainer:    "Cristian Maglie <c.maglie@example.com>",
				License:       "",
				Sentence:      "A library that makes coding a Webserver a breeze.",
				Paragraph:     "Supports HTTP1.1 and you can do GET and POST.",
				URL:           "http://example.com/",
				Architectures: "avr",
				Category:      "Uncategorized",
				Types:         nil,
				Includes:      "WebServer.h",
				Depends:       "ArduinoHttpClient",
			},
			errorAssertion: assert.NoError,
		},
	}

	for _, testTable := range testTables {
		metadata, err := Parse(testTable.propertiesData)
		testTable.errorAssertion(t, err, fmt.Sprintf("%s error", testTable.testName))
		if err == nil {
			assert.Equal(t, testTable.libraryMetadataAssertion, metadata, fmt.Sprintf("%s metadata", testTable.testName))
		}
	}
}
