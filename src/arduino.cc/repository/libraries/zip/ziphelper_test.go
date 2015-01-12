package zip

import (
	"testing"
	"io/ioutil"
	"github.com/stretchr/testify/require"
	"archive/zip"
)

func TestZip(t *testing.T) {
	zipFile, err := ioutil.TempFile("", "ziphelper")
	require.NoError(t, err)

	err = ZipDirectory("./testzip", zipFile)
	require.NoError(t, err)

	zipFileReader, err := zip.OpenReader(zipFile.Name())
	require.NoError(t, err)

	defer zipFileReader.Close()

	require.Equal(t, 2, len(zipFileReader.File))
	require.Equal(t, "testfile.txt", zipFileReader.File[0].Name)
	require.Equal(t, "testfolder/testfileinfolder.txt", zipFileReader.File[1].Name)
}
