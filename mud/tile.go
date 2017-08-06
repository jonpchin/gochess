package mud

import (
	"math/rand"
	"time"
)

// A tile is the smallest unit on the map
// A single movement such as going north, east, south, west, up, down
// will cause a player to move to a different tile
type Tile struct {
	Coordinate                // Contains X (row) and Y (col) coordinates of tile
	Name        string        // Name of the tile an adventurer will see when they enter the room
	Description string        // Description of tile adventurer will see when they enter the room
	Floor       int           // The floor the tile is located
	Area        Area          // The area the tile is located
	Room        Room          // The room the tile is located
	TileType    string        // The type of tile such as floor, wall, openDoor, etc
	Items       []interface{} // List of items or objects in the tile
	Players     []Player      // Adventurers in the tile
	Monsters    []Monster     // NPC in the tile
}

// An area may consist of many rooms and is one way to identify the general
// location of a player, for example when casting the spell locate
type Area struct {
	Name string
}

// asterisk character is reserved for adventurer (self), characters shoulds should be no longer then length one
const (
	UNUSED     = iota // " "
	FLOOR             // "."
	CORRIDOR          // "="
	WALL              // "#"
	CLOSEDOOR         // "+"
	OPENDOOR          // "-"
	UPSTAIRS          // "<"
	DOWNSTAIRS        // ">"
	FOREST            // "$"
	WATER             // "%"
	CLOUD             // "@"
	MOUNTAIN          // "^"
	WHIRLPOOL         // "!"
	UNKNOWN           // ","
)

var tileChars = []string{
	" ",
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

// These types of terrain are most common
var commonTerrainTypes = []string{
	tileChars[FLOOR],
	//tileChars[CORRIDOR],
	tileChars[FOREST],
	//tileChars[WATER],
	//tileChars[CLOUD],
}

type Direction int

const (
	NORTH Direction = iota
	EAST
	SOUTH
	WEST
)

func getTile() {

}

func setTile() {

}

func getRandomTileName() string {
	return "Default Tile Name" // Replace this later
}

func getRandomTileDescription() string {
	return "Default Tile Description" // Replace this later
}
func getRandomArea() Area {
	var area Area
	area.Name = "Default Area"
	return area // Replace this later
}

func getRandomTileChar() string {

	rand.Seed(time.Now().UnixNano())
	randNum := rand.Intn(len(tileChars) - 1)
	return tileChars[randNum]
}

func getCommonTerrainType() string {

	rand.Seed(time.Now().UnixNano())
	randNum := rand.Intn(len(commonTerrainTypes) - 1)
	return commonTerrainTypes[randNum]
}

// Create tile with all its meata data such as name, description, x, y etc
func (tile *Tile) createTile(floorLevel int, area Area, tileCharType string, coordinate Coordinate) {
	tile.Name = getRandomTileName()
	tile.Description = getRandomTileDescription()
	tile.Area = area
	tile.Row = coordinate.Row
	tile.Col = coordinate.Col
	tile.Floor = floorLevel
	//tile.Room =
	// TODO Randomly pick a TileChar but usually its a common type such as floor or trees
	// Need to make the edges walls and not override the door
	tile.TileType = tileCharType
}
