package gostuff

import (
	"fmt"
	"image"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

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

	// TODO: Read these values from a JSON file to avoid having to recompile every time these values
	// are modified
	var gopherPath = "./img/gophers/gopher.png"
	var gamePath = "./img/screenshots/game.png"
	var gameTargetPath = "./img/screenshots/gameResize.png"
	var lobbyPath = "./img/screenshots/lobby.png"
	var lobbyTargetPath = "./img/screenshots/lobbyResize.png"
	var profilePath = "./img/screenshots/profile.png"
	var profileTargetPath = "./img/screenshots/profileResize.png"
	var settingsPath = "./img/screenshots/settings.png"
	var settingsTargetPath = "./img/screenshots/settingsResize.png"
	var screenShotWidth = 600
	var screenShotHeight = 337

	resizeImage(gopherPath, gopherPath, 800, 409)
	resizeImage(gamePath, gameTargetPath, screenShotWidth, screenShotHeight)
	resizeImage(lobbyPath, lobbyTargetPath, screenShotWidth, screenShotHeight)
	resizeImage(profilePath, profileTargetPath, screenShotWidth, screenShotHeight)
	resizeImage(settingsPath, settingsTargetPath, screenShotWidth, screenShotHeight)
}

// resizes image by passing in path to image, the path to the image, the desired width and desired height
// of the target image
// TODO: Have images be deleted if target path already exists
func resizeImage(path string, targetPath string, desiredWidth int, desiredHeight int) {
	problems, _ := os.OpenFile("logs/errors.txt", os.O_APPEND|os.O_WRONLY, 0666)
	defer problems.Close()
	log.SetOutput(problems)

	image, err := imaging.Open(path)
	if err != nil {
		log.Println(err)
		return
	}
	// we want to get the sized of the shrunken image to compare with the desired measurements
	var width int
	var height int
	// if target image does not exist use the image that hasn't been resized
	exists, err := isDirOrFileExists(targetPath)
	if err != nil {
		log.Println(err)
	}
	if exists {
		width, height = getImageDimensions(targetPath)
	} else {
		//renames the extension PNG to png
		os.Rename(path, path)
		width, height = getImageDimensions(path)
	}

	if width != desiredWidth || height != desiredHeight {
		resizedImage := imaging.Resize(image, desiredWidth, desiredHeight, imaging.Lanczos)
		err := imaging.Save(resizedImage, targetPath)
		if err != nil {
			log.Println(err)
			return
		}
	}
	fmt.Println("Images successfully resized")
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

// returns a http client to time out requests that take too long
// @seconds number of seconds for the request before it times out
func TimeOutHttp(seconds time.Duration) http.Client {
	return http.Client{
		Timeout: time.Duration(seconds * time.Second),
	}
}

// Returns true if floats are equal withing the tolerance level set by EPSILON
func IsFloatEqual(a, b float64) bool {
	var EPSILON float64 = 0.00000001
	if (a-b) < EPSILON && (b-a) < EPSILON {
		return true
	}
	return false
}

func timeTrack(start time.Time, name string) {
	log := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}
