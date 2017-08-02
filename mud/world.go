package mud

import (
	"io/ioutil"
	"strconv"
)

// A dungeon consists of floors which can be transversed through stairs
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
		player.loadPlayerData(player.Username)
	}
	// Send player data to client
	connection.sendJSONWebSocket(&player)
}

func CreateWorld() {

	const (
		low       = 3
		high      = 10
		floorLow  = 3
		floorHigh = 30
	)
	numOfFloors := getRandomIntRange(low, high)
	world.Floors = make([]Floor, numOfFloors)

	for i := 0; i < numOfFloors; i += 1 {
		var floor Floor
		floor.Width = getRandomIntRange(floorLow, floorHigh)
		floor.Length = getRandomIntRange(floorLow, floorHigh)
		floor.initFloorTileType()
		floor.makeRooms(i)
		world.Floors[i] = floor
	}
}

// Prints the world with each floor as floor_#.txt in ASCII format
func PrintWorldToFile() {
	for index, floor := range world.Floors {
		floor.writeFloorToFile(index)
	}
}

func (floor *Floor) writeFloorToFile(index int) {

	floorData := ""
	for i := 0; i < len(floor.Plan); i += 1 {
		for j := 0; j < len(floor.Plan[i]); j += 1 {
			floorData += floor.Plan[i][j].TileType
		}
		floorData += "\n"
	}

	ioutil.WriteFile("mud/world/floor_"+strconv.Itoa(index)+".txt", []byte(floorData), 0666)
}
