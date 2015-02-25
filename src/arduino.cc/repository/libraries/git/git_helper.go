package git

import (
	"os/exec"
	"strings"
	"errors"
)

func Clone(repoUrl string, folderName string) error {
	cmd := exec.Command("git", "clone", repoUrl, folderName)
	_, err := cmd.CombinedOutput()
	return err
}

func Fetch(folderName string) error {
	cmd := exec.Command("git", "fetch", "--all")
	cmd.Dir = folderName
	_, err := cmd.CombinedOutput()
	return err
}

func LastTag(folderName string) (string, error) {
	cmd := exec.Command("git", "tag")
	cmd.Dir = folderName
	bytes, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	output := string(bytes)
	rows := filterEmpty(strings.Split(output, "\n"))
	if len(rows) == 0 {
		return "", errors.New("No tags")
	}
	return rows[len(rows)-1], nil
}

func filterEmpty(rows []string) []string {
	var newRows []string
	for _, row := range rows {
		row = strings.TrimSpace(row)
		if len(row) > 0 {
			newRows = append(newRows, row)
		}
	}
	return newRows
}

func CheckoutTag(folderName string, tag string) error {
	cmd := exec.Command("git", "checkout", "refs/tags/"+tag)
	cmd.Dir = folderName
	_, err := cmd.CombinedOutput()
	return err
}
