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

	for i := 0; i < intLength; i++ {
		tiles[i] = make([]string, intLength)

		for j := 0; j < intLength; j++ {
			tiles[i][j] = temp[(i*intLength)+j]
		}
	}
	return tiles
}

// Trims newlines from MUD map file and also removes leading and ending space for the sides
func trimNewlinesAndSides(filename string) {

	// "data/floor_1.txt"
	file, err := os.Open(filename)
	defer file.Close()

	if err != nil {
		fmt.Println(err)
		return
	}

	result := ""
	finalOutput := ""
	start := 9999999
	back := 9999999

	fscanner := bufio.NewScanner(file)

	for fscanner.Scan() {

		result = fscanner.Text()

		if doesStringHaveMapChar(result) {
			for pos, char := range result {
				if char != ' ' {
					if pos < start {
						start = pos
					}
					break
				}
			}
			for pos, _ := range result {
				// Start checking chars from back to see where spaces end
				backIndex := len(result) - pos - 1
				backChar := result[backIndex]

				if backChar != ' ' {
					if pos < back {
						back = pos
					}
					break
				}
			}
			finalOutput += (result + "\n")
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

	trimLeadAndEnd(start, back, filename)
}

func doesStringHaveMapChar(result string) bool {

	// tileChar is tile.go and is array of all possible tiles strings
	mapChars := excludeSpaceInTiles(tileChars)

	for _, mapChar := range mapChars {
		if strings.ContainsAny(result, mapChar) {
			return true
		}
	}
	return false
}

//Returns array exxcluding space in tileChar array
func excludeSpaceInTiles(tiles []string) []string {
	var result []string
	for _, value := range tiles {
		if value != " " {
			result = append(result, value)
		}
	}
	return result
}

// Trims leading and ending spaces of each line in a file
func trimLeadAndEnd(start int, back int, filename string) {

	file, err := os.Open(filename)
	defer file.Close()

	if err != nil {
		fmt.Println(err)
		return
	}

	fscanner := bufio.NewScanner(file)

	result := ""
	finalOutput := ""

	for fscanner.Scan() {

		result = fscanner.Text()

		if len(result) != 0 {
			temp := result[start:(len(result) - back)]
			finalOutput += (temp + "\n")

			if err != nil {
				fmt.Println(err)
				break
			}
		}
	}
	err = ioutil.WriteFile(filename, []byte(finalOutput), 0644)
	if err != nil {
		fmt.Println(err)
	}
}
