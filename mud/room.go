package mud

import "fmt"

type Room struct {
	Tiles []Tile // List of all tiles in room
	Wall  []Tile // List of tiles that make the wall of the room
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
	floor.makeRoom(floor.Width/2, floor.Length/2, getRandomDirection(), tileChars[FLOOR])
	// TODO: If failed to make a room decrement i
	for i := 1; i < numRooms; i += 1 {
		tile := floor.getRandomTileOnWall()
		floor.makeRoom(tile.Row/2, tile.Col/2, getRandomDirection(), getCommonTerrainType())
	}
}

// x and y is the tile location which connects the current rooms
// to the next room. Direction will be used to check, terrainType is a tileChar
// such as floor, cloud, moutain, etc, except special terrain such as unused, door and whirlpool
// Returns true if a room is succesfully created
func (floor *Floor) makeRoom(row int, col int, dir Direction, terrainType string) bool {

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

		//       O
		// . . . C . ..x
		// .
		// x
		// .
		// .
	} else if dir == SOUTH {

		topLeft.Row = row - 1
		topLeft.Col = col - (width / 2)
		bottomRight.Row = row + 1 + length
		bottomRight.Col = col + (width / 2)

		//                  .
		//                  .
		//                  .
		//                  .
		//                  X O
		//                  .
		//                  .
		//  . . . . x . . . .
	} else if dir == WEST {

		topLeft.Row = row - (length / 2)
		topLeft.Col = col - 1 - (width / 2)
		bottomRight.Row = row + (length / 2)
		bottomRight.Col = col - 1

	} else {
		fmt.Println("Error invalid direction for makeRoom", row, col, dir)
	}

	isCreated := false

	if floor.isRoomUsed(topLeft, bottomRight) {

	} else {
		floor.createTilesInRoom(topLeft, bottomRight, terrainType)
		isCreated = true
	}
	return isCreated
}

// If the room is already occupied return true
func (floor *Floor) isRoomUsed(topLeft, bottomRight Coordinate) bool {
	for i := topLeft.Row; i < bottomRight.Row; i += 1 {
		for j := topLeft.Col; j < bottomRight.Col; j += 1 {
			if floor.isValidCoordinate(i, j) && floor.Plan[i][j].TileType != tileChars[UNUSED] {
				return true
			}
		}
	}
	return false
}

func (floor *Floor) isValidCoordinate(row, col int) bool {
	if row < 0 || row >= floor.Length {
		return false
	}
	if col < 0 || col >= floor.Width {
		return false
	}
	return true
}

// Builds all the tiles in the room
func (floor *Floor) createTilesInRoom(topLeft, bottomRight Coordinate, terrainType string) {
	var area Area
	area = getRandomArea()
	var room Room
	room.Wall = make([]Tile, 1)
	for i := topLeft.Row; i <= bottomRight.Row; i += 1 {
		for j := topLeft.Col; j <= bottomRight.Col; j += 1 {
			if floor.isValidCoordinate(i, j) {

				// Adds the edge tiles to the wall list of each room
				if i == topLeft.Row || i == bottomRight.Row ||
					j == topLeft.Col || j == bottomRight.Col {
					room.Wall = append(room.Wall, floor.Plan[i][j])
					floor.Plan[i][j].createTile(floor.Level, area, tileChars[WALL])
				} else {
					floor.Plan[i][j].createTile(floor.Level, area, terrainType)
				}

				// Adds all the tiles in the room
				room.Tiles = append(room.Tiles, floor.Plan[i][j])
			}
		}
	}
	floor.Rooms = append(floor.Rooms, room)
}

// Create tile with all its meata data such as name, description, x, y etc
func (tile *Tile) createTile(floorLevel int, area Area, tileCharType string) {
	tile.Name = getRandomTileName()
	tile.Description = getRandomTileDescription()
	tile.Floor = floorLevel
	tile.Area = area
	//tile.Room =
	// TODO Randomly pick a TileChar but usually its a common type such as floor or trees
	// Need to make the edges walls and not override the door
	tile.TileType = getRandomTileChar()
}
