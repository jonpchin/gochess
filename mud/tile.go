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
	unused     TileChar = " "
	floor      TileChar = "."
	corridor   TileChar = ","
	wall       TileChar = "#"
	closedDoor TileChar = "+"
	openDoor   TileChar = "-"
	upStairs   TileChar = "<"
	downStairs TileChar = ">"
	forest     TileChar = "$"
	water      TileChar = "%"
	cloud      TileChar = "@"
	mountain   TileChar = "^"
	whirlpool  TileChar = "!"
)

type Direction int

const (
	north Direction = iota
	east
	south
	west
)

func getTile() {

}

func setTile() {

}
