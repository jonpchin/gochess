package mud

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strings"
)

// Converts map string to 2D array of string
func convertMapToTiles(mapString string) [][]string {

	length := math.Sqrt(float64(len(mapString)))
	intLength := int(length)
	var tiles = make([][]string, intLength)

	var temp []string

	for _, r := range mapString {
		temp = append(temp, string(r))
	}

	for i := 0; i < intLength; i += 1 {
		tiles[i] = make([]string, intLength)

		for j := 0; j < intLength; j += 1 {
			tiles[i][j] = temp[(i*intLength)+j]
		}
	}
	return tiles
}

// Trims newlines from MUD map file
func trimNewlines(filename string) {

	// "data/floor_1.txt"
	file, err := os.Open(filename)
	defer file.Close()

	if err != nil {
		fmt.Println(err)
		return
	}

	// Start reading from the file with a reader.
	reader := bufio.NewReader(file)

	result := ""
	finalOutput := ""

	for {
		content, err := reader.ReadString('\n')

		result = string(content)

		if doesStringHaveMapChar(result) {
			finalOutput += result
		}

		if err != nil {
			fmt.Println(err)
			break
		}
	}

	err = ioutil.WriteFile(filename, []byte(finalOutput), 0644)
	if err != nil {
		fmt.Println(err)
	}
}

func doesStringHaveMapChar(result string) bool {

	var mapChars = []string{
		".",
		"=",
		"#",
		"+",
		"-",
		"<",
		">",
		"$",
		"%",
		"@",
		"^",
		"!",
		",",
	}

	for _, mapChar := range mapChars {
		if strings.ContainsAny(result, mapChar) {
			return true
		}
	}
	return false
}
