package gostuff

import (
	"fmt"
	"image"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/disintegration/imaging"
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
	err := archiver.Zip.Make(destination, source)
	if err != nil {
		fmt.Println("util.go compress There was an error in compressing the file", err)
		return false
	}
	return true
}

// Unzips an archive and places it in a destination folder
func unzip(source string, destination string) bool {
	err := archiver.Zip.Open(source, destination)
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

// check if images need to be resized, if they do then resize them
func ResizeImages() {
	problems, _ := os.OpenFile("logs/errors.txt", os.O_APPEND|os.O_WRONLY, 0666)
	defer problems.Close()
	log.SetOutput(problems)

	// TODO: Read these values from a JSON file to avoid having to recompile every time these values
	// are modified
	var gopherPath = "./img/gophers/gopher.png"
	var expectedGopherWidth = 800
	var expectedGopherHeight = 409

	gopherImage, err := imaging.Open(gopherPath)
	if err != nil {
		log.Println(err)
		return
	}

	gopherWidth, gopherHeight := getImageDimensions(gopherPath)
	if gopherWidth != expectedGopherWidth || gopherHeight != expectedGopherHeight {
		resizedGopher := imaging.Resize(gopherImage, expectedGopherWidth, expectedGopherHeight, imaging.Lanczos)
		err := imaging.Save(resizedGopher, "./img/gophers/gopher.png")
		if err != nil {
			log.Println(err)
			return
		}
	}
}

// pass in file path to image, returns width and heigh of image
func getImageDimensions(imagePath string) (int, int) {

	problems, _ := os.OpenFile("logs/errors.txt", os.O_APPEND|os.O_WRONLY, 0666)
	defer problems.Close()
	log.SetOutput(problems)

	file, err := os.Open(imagePath)
	defer file.Close()
	if err != nil {
		log.Println(err)
	}

	image, _, err := image.DecodeConfig(file)
	if err != nil {
		log.Println(err)
	}
	return image.Width, image.Height
}
