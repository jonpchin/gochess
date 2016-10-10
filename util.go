package gostuff

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/mholt/archiver"
)

// Returns true if a given file or directory exist in the path
func isDirOrFileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

// Archives a group of files as an array of string into a destination zip file
func compress(destination string, source []string) bool {
	err := archiver.Zip(destination, source)
	if err != nil {
		fmt.Println("util.go compress There was an error in compressing the file", err)
		return false
	}
	return true
}

// Unzips an archive and places it in a destination folder
func unzip(source string, destination string) bool {
	err := archiver.Unzip(source, destination)
	if err != nil {
		fmt.Println("util.go unzip There was an error in archiving the file", err)
		return false
	}
	return true
}

// Replaces target string in file with desired string in the file path
// Returns false if there was an error in the operation
func ReplaceString(target, desired, source string) bool {

	input, err := ioutil.ReadFile(source)
	if err != nil {
		fmt.Println("Failed to replaceString 1 util.go", err)
		return false
	}

	lines := strings.Split(string(input), "\n")

	for i, line := range lines {
		if strings.Contains(line, target) {
			lines[i] = strings.Replace(lines[i], target, desired, -1)
		}
	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(source, []byte(output), 0644)
	if err != nil {
		fmt.Println("Failed to replaceString 2 util.go", err)
		return false
	}
	return true
}
