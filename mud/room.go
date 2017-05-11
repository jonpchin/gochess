package mud

import "fmt"

type Room struct {
	Tiles [][]Tile // A room has x by y tiles
}

type Coordinate struct {
	Row   int // X
	Col   int // Y
	Level int // Z or level of floor
}

func (floor *Floor) makeRooms(floorLevel int) {

	const (
		roomsLow  = 50
		roomsHigh = 200
	)

	numRooms := getRandomIntRange(roomsLow, roomsHigh)

	// Place the first room in the center
	floor.makeRoom(floor.Width/2, floor.Length/2, getRandomDirection(), true)

	for i := 1; i < numRooms; i += 1 {

	}

}

// x and y is the tile location which connects the current rooms
// to the next room. Direction will be used to check
func (floor *Floor) makeRoom(row int, col int, dir Direction, firstRoom bool) {

	const (
		roomDimensionLow   = 7
		roomsDimensionHigh = 40
	)

	width := getRandomIntRange(roomDimensionLow, roomsDimensionHigh)
	length := getRandomIntRange(roomDimensionLow, roomsDimensionHigh)

	var topLeft Coordinate
	var bottomRight Coordinate

	// x is top left and bottom bottomRight
	// O is where x and y coordinate is located
	//
	// x
	// .
	// .
	// . . . C . ..x
	//       O
	if dir == NORTH {

		topLeft.Row = (row - 1) - length
		topLeft.Col = col - (width / 2)
		bottomRight.Row = row - 1
		bottomRight.Col = col + (width / 2)

		if floor.isRoomUsed(topLeft, bottomRight) {

		} else {
			floor.createTilesInRoom(topLeft, bottomRight)
		}

		//   x
		//   .
		//   .
		// O C
		//   .
		//   .
		//   . . . . . . .x
	} else if dir == EAST {

		topLeft.Row = row - (length / 2)
		topLeft.Col = col + 1
		bottomRight.Row = row + (length / 2)
		bottomRight.Col = col + 1 + width

		if floor.isRoomUsed(topLeft, bottomRight) {

		}
	} else if dir == SOUTH {

	} else if dir == WEST {

	} else {
		fmt.Println("Error invalid direction for makeRoom", row, col, dir, firstRoom)
	}
}

// If the room is already occupied return true
func (floor *Floor) isRoomUsed(topLeft, bottomRight Coordinate) bool {
	for i := topLeft.Row; i < bottomRight.Row; i += 1 {
		for j := topLeft.Col; j < bottomRight.Col; j += 1 {
			if floor.isValidCoordinate(i, j) && floor.Plan[i][j].TileType != UNUSED {
				return true
			}
		}
	}
	return false
}

func (floor *Floor) isValidCoordinate(row, col int) bool {
	if row < 0 || row > floor.Length {
		return false
	}
	if col < 0 || col > floor.Width {
		return false
	}
	return true
}

// Builds all the tiles in the room
func (floor *Floor) createTilesInRoom(topLeft, bottomRight Coordinate) {
	var area Area
	area = getRandomArea()
	for i := topLeft.Row; i < bottomRight.Row; i += 1 {
		for j := topLeft.Col; j < bottomRight.Col; j += 1 {
			if floor.isValidCoordinate(i, j) {
				floor.Plan[i][j].createTile(floor.Level, area)
			}
		}
	}
}

// Create tile with all its meata data such as name, description, x, y etc
func (tile *Tile) createTile(floorLevel int, area Area) {
	tile.Name = getRandomTileName()
	tile.Description = getRandomTileDescription()
	tile.Floor = floorLevel
	tile.Area = area
	//tile.Room =
	// TODO Randomly pick a TileChar but usually its a common type such as floor or trees
	// Need to make the edges walls and not override the door
	tile.TileType = FLOOR
}

// Selects a random tile on the wall of a room
func (room *Room) selectTileOnWall() {

}
