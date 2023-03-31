package utils

import (
	"io/ioutil"
	"strings"
)

// IsDirEmpty checks if a directory is empty or not
func IsDirEmpty(path string) (bool, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return false, err
	}

	for _, file := range files {
		if !file.IsDir() && (strings.HasPrefix(file.Name(), ".") || strings.HasSuffix(file.Name(), ".DS_Store")) {
			// Skip hidden files or files with .DS_Store extension
			continue
		}

		return false, nil
	}

	return true, nil
}
