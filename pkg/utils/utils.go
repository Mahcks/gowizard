package utils

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/mgutz/ansi"
)

func PrintError(msg string, args ...any) {
	fmt.Println(ansi.Color("[✗] Error:", "red"), ansi.Color(fmt.Sprintf(msg, args...), "red"), ansi.ColorCode("reset"))
}

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
