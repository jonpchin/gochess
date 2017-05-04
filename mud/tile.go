package mud

// A tile is the smallest unit on the map
// A single movement such as going north, east, south, west, up, down
// will cause a player to move to a different tile
type Tile struct {
	X           int
	Y           int
	Name        string
	Description string
	Floor       int
	Area        Area
	Room        Room
}

// An area may consist of many rooms and is one way to identify the general
// location of a player, for example when casting the spell locate
type Area struct {
	Name string
}

type TileChar string

const (
	unused     TileChar = " "
	floor      TileChar = "."
	corridor   TileChar = ","
	wall       TileChar = "#"
	closedDoor TileChar = "+"
	openDoor   TileChar = "-"
	upStairs   TileChar = "<"
	downStairs TileChar = ">"
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
