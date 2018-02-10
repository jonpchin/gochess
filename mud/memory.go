package mud

import(
	"io/ioutil"
	"fmt"
	"bufio"
	"os"
	"strings"
	"strconv"
)

// Loads the world map by ID into memory in /world directory
func LoadMapsIntoMemory(id string){
	worldFolder := "./world/" + id
	
	files, err := ioutil.ReadDir(worldFolder)
    if err != nil {
        fmt.Println(err)
		return
	}
	
	world.Floors = make([]Floor, len(files))

    for _, f := range files {
		fileName := worldFolder + "/" + f.Name()
		file, err := os.Open(fileName)
		defer file.Close()

		if err != nil {
			fmt.Println(err)
			return
		}

		fscanner := bufio.NewScanner(file)
		var floorLevel int
		floorLevel = 0

		fileTokens := strings.Split(f.Name(), "_")
		if len(fileTokens) > 2{
			floorLevel, err = strconv.Atoi(fileTokens[1]) 
			if err != nil{
				fmt.Println("Can't convert floor level string to int", err)
				return
			}
		}else{
			fmt.Println("Error parsing file token when loading map into memory")
			return
		}
		
		i := 0
		j := 0
		linesInFile := getLinesInFile(fileName)
		world.Floors[floorLevel].Plan = make([][]Tile, linesInFile)

		for fscanner.Scan() {
			result := fscanner.Text()
			
			world.Floors[floorLevel].Plan[i] = make([]Tile, len(result))
			for _, value := range result{
				//TODO: Need to test
				world.Floors[floorLevel].Plan[i][j].TileType = string(value)
				// TODO: Read metadata of tile from tile metadata
				j+=1
			}
			i+=1
		}
    }
}

// Returns number of lines in a file
func getLinesInFile(filepath string) int{
	file, err := os.Open(filepath)
	defer file.Close()

	if err != nil {
		fmt.Println(err)
		return 0
	}

	fscanner := bufio.NewScanner(file)
	count := 0
	for fscanner.Scan() {
		count+=1
	}

	return count
}