package zip

import (
	"archive/zip"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestZip(t *testing.T) {
	zipFile, err := ioutil.TempFile("", "ziphelper*.zip")
	require.NoError(t, err)
	require.NotNil(t, zipFile)
	zipFileName := zipFile.Name()
	require.NoError(t, zipFile.Close())
	require.NoError(t, os.Remove(zipFileName))
	defer os.RemoveAll(zipFileName)

	err = Directory("./testzip", "a_zip", zipFileName)
	require.NoError(t, err)

	zipFileReader, err := zip.OpenReader(zipFileName)
	require.NoError(t, err)

	defer zipFileReader.Close()

	require.Equal(t, 4, len(zipFileReader.File))

	containsName := func(name string) bool {
		for _, file := range zipFileReader.File {
			if file.Name == name {
				return true
			}
		}

		return false
	}
	require.True(t, containsName("a_zip/"))
	require.True(t, containsName("a_zip/testfile.txt"))
	require.True(t, containsName("a_zip/testfolder/"))
	require.True(t, containsName("a_zip/testfolder/testfileinfolder.txt"))
}
