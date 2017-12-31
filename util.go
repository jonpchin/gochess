package gostuff

import (
	"fmt"
	"image"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
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

// Returns true if the number of seconds is greater then the difference of targetTime - (this moment)
// Also returns time difference of targetTime and now, returns zero could mean there was an error
// timeFormat is the time format targetTime is in
// useHour is true to get difference in hours for rating history and false for forum spam control
// timeCompare can be in hours or seconds depending on useHour
func HasTimeElapsed(targetTime string, timeCompare int, timeFormat string, useHour bool) (bool, int) {

	log := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

	// If not valid that means there is no existing timestamp in the database
	then, err := time.Parse(timeFormat, targetTime)
	if err != nil {
		log.Println(err)
		return false, 0
	}

	duration := time.Now().Sub(then)

	var timeZoneDiff float64

	if isWindows {
		// UTC-5 is Eastern US time
		timeZoneDiff = 14400.0
	} else {
		timeZoneDiff = 0
	}

	var timeDiff int
	// If greater then 120 seconds its not forum control but rating history being managed
	if useHour {
		timeDiff = int(duration.Hours())
	} else {
		timeDiff = int(duration.Seconds() - timeZoneDiff)
	}

	// A certain number of seconds need to pass since target time
	if timeDiff < timeCompare {
		return false, timeDiff
	}
	return true, 0
}

func timeTrack(start time.Time, name string) {
	log := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}

// Validates all json files in project
func ValidateJSONFiles() {
	// windows use double backslash for path
	RecurseDirectory("data", validateJSONFile, "*.json")
	RecurseDirectory("mud/equipment", validateJSONFile, "*.json")
}

// Checks if a single JSON file is valid and if its not print the error message to console
func validateJSONFile(file string) {

	log := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
	var output []byte
	var err error

	if runtime.GOOS == "windows" {
		output, err = exec.Command("jsonlint", "-q", "-c", file).CombinedOutput()
	} else {
		output, err = exec.Command("/bin/bash", "-c", "jsonlint -q -c "+file).Output()
	}

	if err != nil {
		log.Println(err)
	}
	if string(output) != "" {
		log.Println(string(output))
	}
}

// Recursives through directory and calls the function pointer fp on each file
// Pattern is the file pattern to apply the function to for example
// *json would apply the fp on all files recursively in the directory
func RecurseDirectory(searchDir string, fp func(string), pattern string) {

	log := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

	fileList := []string{}
	err := filepath.Walk(searchDir, func(path string, f os.FileInfo, err error) error {
		fileList = append(fileList, path)
		return nil
	})
	if err != nil {
		log.Println(err)
	}

	for _, file := range fileList {
		isDir := isDirectory(file)
		if isDir == false {
			isMatch, err := filepath.Match(pattern, filepath.Base(file))
			if err != nil {
				log.Println(err)
			} else if isMatch {
				fp(file)
			}
		}
	}
}

//return trues if path is a directory
func isDirectory(path string) bool {
	f, err := os.Open(path)
	if err != nil {
		fmt.Println("isDirectory 1", err)
		return false
	}
	stat, err := f.Stat()
	if err != nil {
		fmt.Println("isDirectory 2", err)
	}
	return stat.IsDir()
}

// Print out memory usage
func PrintMemoryStats() {

	var mem runtime.MemStats

	log := log.New(os.Stdout, "", log.LstdFlags)
	log.Println("Printing memory stats...")

	runtime.ReadMemStats(&mem)
	//bytes of allocated heap objects
	log.Println("Alloc: ", mem.Alloc)
	//cumulative bytes allocated for heap objects
	log.Println("Total alloc: ", mem.TotalAlloc)
	log.Println("Heap alloc:", mem.HeapAlloc)
	log.Println("Heap system:", mem.HeapSys)
}
