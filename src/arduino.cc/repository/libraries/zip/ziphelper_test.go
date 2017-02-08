package zip

import (
	"archive/zip"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestZip(t *testing.T) {
	zipFile, err := ioutil.TempFile("", "ziphelper")
	require.NoError(t, err)
	defer os.RemoveAll(zipFile.Name())

	err = ZipDirectory("./testzip", "a_zip", zipFile)
	require.NoError(t, err)

	zipFileReader, err := zip.OpenReader(zipFile.Name())
	require.NoError(t, err)

	defer zipFileReader.Close()

	require.Equal(t, 2, len(zipFileReader.File))
	require.Equal(t, "a_zip/testfile.txt", zipFileReader.File[0].Name)
	require.Equal(t, "a_zip/testfolder/testfileinfolder.txt", zipFileReader.File[1].Name)
}
