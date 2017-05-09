package mud

import "fmt"

type Room struct {
	Tiles [][]Tile // A room has x by y tiles
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

	// x is top left and bottom bottomRight
	// O is where x and y coordinate is located
	//
	// x
	// .
	// .
	// . . . C . ..x
	//       O

	if dir == north {

		if floor.isRoomUsed((row-1)-length, col-(width/2), row-1, col+(width/2)) {

		}

	} else if dir == east {

	} else if dir == south {

	} else if dir == west {

	} else {
		fmt.Println("Error invalid direction for makeRoom", row, col, dir, firstRoom)
	}
}

// If the room is already occupied return true
func (floor *Floor) isRoomUsed(topLeftX, topLeftY, bottomRightX, bottomRightY int) bool {
	for i := topLeftY; i < bottomRightY; i += 1 {
		for j := topLeftX; j < bottomRightX; j += 1 {
			if floor.isValidCoordinate(i, j) && floor.Plan[i][j] != unused {
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
func createTilesInRoom() {

}

// Create tile with all its meata data such as name, description, x, y etc
func (tile *Tile) createTile() {

}

// Selects a random tile on the wall of a room
func (room *Room) selectTileOnWall() {

}
