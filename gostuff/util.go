package gostuff

import(
	"os"
	"github.com/mholt/archiver"
	"fmt"
)

// Returns true if a given file or directory exist in the path
func isDirOrFileExists(path string) (bool, error){
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
func compress(destination string, source []string) bool{
	err := archiver.Zip(destination, source)
	if err != nil {
		fmt.Println("util.go compress There was an error in compressing the file", err)
		return false
	}
	return true
}

// Unzips an archive and places it in a destination folder
func unzip(source string, destination string) bool{
	err := archiver.Unzip(source, destination)
	if err != nil {
		fmt.Println("util.go unzip There was an error in archiving the file", err)
		return false
	}
	return true
}