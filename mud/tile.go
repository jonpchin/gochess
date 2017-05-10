package mud

// A tile is the smallest unit on the map
// A single movement such as going north, east, south, west, up, down
// will cause a player to move to a different tile
type Tile struct {
	Row         int      // x coordinate
	Col         int      // y coordinate
	Name        string   // Name of the room an adventurer will see when they enter the room
	Description string   // Description of room will adventurer will see when they enter the room
	Floor       int      // The floor the tile is located
	Area        Area     // The area the tile is located
	Room        Room     // The room the tile is located
	TileType    TileChar // The type of tile such as floor, wall, openDoor, etc
}

// An area may consist of many rooms and is one way to identify the general
// location of a player, for example when casting the spell locate
type Area struct {
	Name string
}

type TileChar string

// asterisk character is reserved for adventurer (self)
const (
	UNUSED     TileChar = " "
	FLOOR      TileChar = "."
	CORRIDOR   TileChar = ","
	WALL       TileChar = "#"
	CLOSEDOOR  TileChar = "+"
	OPENDOOR   TileChar = "-"
	UPSTAIRS   TileChar = "<"
	DOWNSTAIRS TileChar = ">"
	FOREST     TileChar = "$"
	WATER      TileChar = "%"
	CLOUD      TileChar = "@"
	MOUNTAIN   TileChar = "^"
	WHIRLPOOL  TileChar = "!"
)

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
