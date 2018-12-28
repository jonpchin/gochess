package mud

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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
		//floor.pruneFloorPlan()
		world.Floors[i] = floor
	}
}

// Loads entire world from json file to memory
func LoadWorldFile(id string) {

	data, err := ioutil.ReadFile("mud/tile_metadata/" + id + ".json")
	if err != nil {
		log.Fatal(err)
	}

	var tempWorld World

	err = json.Unmarshal(data, &tempWorld)
	if err != nil {
		log.Fatal(err)
	}
	world = tempWorld
}

// Sets the map vision of the player based on his coordinate and vision (default is 5 vision)
func (player *Player) setMapVision(tempWorld World) {

	// TODO: Ensure player vision is not tampered
	player.Vision = 5

	lowestRow := player.Location.Row - player.Vision
	highestRow := player.Location.Row + player.Vision
	lowestCol := player.Location.Col - player.Vision
	highestCol := player.Location.Col + player.Vision

	if lowestRow < 0 {
		lowestRow = 0
	}

	if lowestCol < 0 {
		lowestCol = 0
	}

	if highestCol >= tempWorld.Floors[player.Location.Level].Width {
		highestCol = tempWorld.Floors[player.Location.Level].Width - 1
	}

	if highestRow >= tempWorld.Floors[player.Location.Level].Length {
		highestRow = tempWorld.Floors[player.Location.Level].Length - 1
	}

	player.Map = ""

	for i := lowestRow; i < highestRow; i++ {
		for j := lowestCol; j < highestCol; j++ {
			player.Map += tempWorld.Floors[player.Location.Level].Plan[player.Location.Row][player.Location.Col].TileType
		}
	}
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

	// Disabled for now until pruneFloorPlan is working
	//trimNewlinesAndSides(filename)
}
