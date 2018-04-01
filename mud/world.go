package mud

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
)

// A dungeon consists of floors which can be transversed through stairs
// Order of Floor in Floors determines level, zero index is first floor,
// 1st index is 2nd floor, etc
type World struct {
	Floors []Floor
}

var world World

const (
	LOAD_PLAYER = true
	SKIP_LOAD   = false
)

func (player *Player) enterWorld(loadPlayer bool, connection *MudConnection) {

	if loadPlayer {
		player.loadPlayerData()
		player.loadMap()
	}
	// Send player data to client
	connection.sendJSONWebSocket(&player)
}

func CreateWorld() {

	const (
		low       = 2
		high      = 4
		floorLow  = 50
		floorHigh = 200
	)
	numOfFloors := getRandomIntRange(low, high)
	world.Floors = make([]Floor, numOfFloors)

	for i := 0; i < numOfFloors; i += 1 {
		var floor Floor
		floor.Width = getRandomIntRange(floorLow, floorHigh)
		floor.Length = getRandomIntRange(floorLow, floorHigh)
		floor.Level = i
		floor.initFloorTileType()
		floor.makeRooms(i)
		world.Floors[i] = floor
	}
}

// Loads entire world from file to memory
func loadWorldFile() {

}

func SaveMetadataToFile(id string) {
	jsonWorld, _ := json.Marshal(world)
	err := ioutil.WriteFile("mud/tile_metadata/"+id+".json", jsonWorld, 0644)
	if err != nil {
		fmt.Println("SavedMetadataToFile", err)
	}
}

// Prints the world with each floor as floor_#.txt in ASCII format
func PrintWorldToFile(worldNumber string) {
	for index, floor := range world.Floors {
		floor.writeFloorToFile(index, worldNumber)
	}
}

func (floor *Floor) writeFloorToFile(index int, worldNumber string) {

	floorData := ""
	for i := 0; i < len(floor.Plan); i += 1 {
		for j := 0; j < len(floor.Plan[i]); j += 1 {
			floorData += floor.Plan[i][j].TileType
		}
		floorData += "\n"
	}
	filename := "mud/world/" + worldNumber + "/floor_" + strconv.Itoa(index) + ".txt"
	ioutil.WriteFile(filename, []byte(floorData), 0666)
	trimNewlinesAndSides(filename)
}
