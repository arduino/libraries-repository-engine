package cron

import "crypto/sha256"
import "encoding/hex"
import "io"
import "net/http"
import "os"

/*
   Check for missing size and checksum field and fills them
   by downloading a copy of the file.
*/
func FillMissingChecksumsForDownloadArchives(URL string, filename string) (int64, string, error) {
	size, err := download(URL, filename)
	if err != nil {
		return 0, "", err
	}

	hash, err := checksum(filename)
	if err != nil {
		os.Remove(filename)
		return 0, "", err
	}

	return size, hash, nil
}

// Calculate hash for the file
func checksum(filename string) (string, error) {
	hasher := sha256.New()
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}
	return "SHA-256:" + hex.EncodeToString(hasher.Sum(nil)), nil
}

func download(URL string, filename string) (int64, error) {
	out, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		return 0, err
	}
	defer out.Close()

	resp, err := http.Get(URL)
	if err != nil {
		defer os.Remove(out.Name())
		return 0, err
	}
	defer resp.Body.Close()

	size, err := io.Copy(out, resp.Body)
	if err != nil {
		defer os.Remove(out.Name())
		return 0, err
	}

	return size, nil
}

// vi:ts=2
