package mud

import (
	"log"
	"os"
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
	CORRIDOR          // ","
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
)

var tileChars = []string{
	" ",
	".",
	",",
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
}

// These types of terrain are most common
var commonTerrainTypes = []string{
	tileChars[FLOOR],
	tileChars[CORRIDOR],
	tileChars[FOREST],
	tileChars[WATER],
	tileChars[CLOUD],
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
	log := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
	randNum, err := secureRandomInt(int64(len(tileChars) - 1))
	if err != nil {
		log.Println(err)
	}
	return tileChars[randNum]
}

func getCommonTerrainType() string {
	log := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
	randNum, err := secureRandomInt(int64(len(commonTerrainTypes) - 1))
	if err != nil {
		log.Println(err)
	}
	return tileChars[randNum]
}
