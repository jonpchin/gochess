package gostuff

import(
	"os"
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