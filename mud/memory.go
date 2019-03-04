package mud

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// Loads the world map by ID into memory in /world directory
// Creates tile meta data if no existing metadata exist otherwise load it
func LoadMapsIntoMemory(id string) {
	worldFile := "mud/tile_metadata/" + id + ".json"

	data, err := ioutil.ReadFile(worldFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = json.Unmarshal(data, &world)
	if err != nil {
		fmt.Println("LoadMapsIntoMemory 0", err)
	}

	fmt.Println("Map loaded into memory")
}

// Returns number of lines in a file
func getLinesInFile(filepath string) int {
	file, err := os.Open(filepath)
	defer file.Close()

	if err != nil {
		fmt.Println(err)
		return 0
	}

	fscanner := bufio.NewScanner(file)
	count := 0
	for fscanner.Scan() {
		count++
	}

	return count
}
