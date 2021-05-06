package libraries

import (
	"arduino.cc/repository/libraries/hash"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

// GithubDownloadRelease downloads GitHub's archive of the release.
func GithubDownloadRelease(repoURL, version string) (string, int64, string, error) {
	tempfile, err := ioutil.TempFile("", "github")
	if err != nil {
		return "", -1, "", err
	}
	defer os.Remove(tempfile.Name())

	zipFileURL := strings.Replace(repoURL, ".git", "", 1) + "/archive/" + version + ".zip"

	err = saveURLIn(zipFileURL, tempfile)
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

	return zipFileURL, size, checksum, nil
}

func saveURLIn(url string, tempfile *os.File) error {
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
