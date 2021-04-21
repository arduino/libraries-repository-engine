package libraries

import (
	"arduino.cc/repository/libraries/hash"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func GithubDownloadRelease(repoUrl, version string) (string, int64, string, error) {
	tempfile, err := ioutil.TempFile("", "github")
	if err != nil {
		return "", -1, "", err
	}
	defer os.Remove(tempfile.Name())

	zipFileUrl := strings.Replace(repoUrl, ".git", "", 1) + "/archive/" + version + ".zip"

	err = saveUrlIn(zipFileUrl, tempfile)
	if err != nil {
		return "", -1, "", err
	}

	info, err := os.Stat(tempfile.Name())
	if err != nil {
		return "", -1, "", err
	}
	size := info.Size()

	checksum, err := hash.Checksum(tempfile.Name())
	if err != nil {
		return "", -1, "", err
	}

	return zipFileUrl, size, checksum, nil
}

func saveUrlIn(url string, tempfile *os.File) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer tempfile.Close()
	defer resp.Body.Close()

	_, err = io.Copy(tempfile, resp.Body)
	if err != nil {
		return err
	}
	return nil
}
