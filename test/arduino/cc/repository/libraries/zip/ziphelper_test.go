package zip

import (
	"archive/zip"
	arduinozip "arduino.cc/repository/libraries/zip"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"testing"
)

func TestZip(t *testing.T) {
	zipFile, err := ioutil.TempFile("", "ziphelper")
	require.NoError(t, err)

	err = arduinozip.ZipDirectory("./testzip", zipFile)
	require.NoError(t, err)

	zipFileReader, err := zip.OpenReader(zipFile.Name())
	require.NoError(t, err)

	defer zipFileReader.Close()

	require.Equal(t, 2, len(zipFileReader.File))
	require.Equal(t, "testfile.txt", zipFileReader.File[0].Name)
	require.Equal(t, "testfolder/testfileinfolder.txt", zipFileReader.File[1].Name)
}
