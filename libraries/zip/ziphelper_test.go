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
	require.Equal(t, "a_zip/", zipFileReader.File[0].Name)
	require.Equal(t, "a_zip/testfile.txt", zipFileReader.File[1].Name)
	require.Equal(t, "a_zip/testfolder/", zipFileReader.File[2].Name)
	require.Equal(t, "a_zip/testfolder/testfileinfolder.txt", zipFileReader.File[3].Name)
}
