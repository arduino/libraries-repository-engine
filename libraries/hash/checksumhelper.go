package hash

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
)

// Checksum calculates the hash for the file.
func Checksum(filename string) (string, error) {
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
