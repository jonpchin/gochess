package mud

import (
	"fmt"
	"math/rand"
	"time"
)

type Room struct {
	Name  string // name of room
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
		roomsLow  = 800
		roomsHigh = 1000

		roomDimensionLow   = 7
		roomsDimensionHigh = 40

		lowDistance  = 5
		highDistance = 10
	)

	rand.Seed(time.Now().UnixNano())
	distanceBetweenRooms := rand.Intn(highDistance) + lowDistance

	numRooms := getRandomIntRange(roomsLow, roomsHigh)

	width := getRandomIntRange(roomDimensionLow, roomsDimensionHigh)
	length := getRandomIntRange(roomDimensionLow, roomsDimensionHigh)

	// Place the first room in the center
	floor.makeRoom(floor.Width/2, floor.Length/2, getRandomDirection(),
		width, length, distanceBetweenRooms)

	// TODO: If failed to make a room decrement i
	for i := 1; i < numRooms; i += 1 {
		room := floor.getRandomRoomOnFloor()
		tile := room.getRandomTileOnWall()
		rand.Seed(time.Now().UnixNano())
		width = rand.Intn(roomsDimensionHigh) + roomDimensionLow
		length = rand.Intn(roomsDimensionHigh) + roomDimensionLow

		if floor.isRoomWithinFloorDimensions(tile.Row, tile.Col, width, length) {
			direction := getRandomDirection()
			if floor.makeRoom(tile.Row, tile.Col, direction,
				width, length, distanceBetweenRooms) {
				tile.createWallFeature()
				// For now just assign tileType later use entire tile when it has more metadata
				floor.Plan[tile.Row][tile.Col].TileType = tile.TileType
				floor.createCorridorInDirection(tile, direction, distanceBetweenRooms)
			}
		}
	}
}

// Prunes floor plan in memory, this is done to match the file pruning, this is experimental and not working
func (floor *Floor) pruneFloorPlan() {
	left := floor.Width
	right := 0
	top := 0
	bottom := floor.Length
	foundTile := false

	for rowIndex, row := range floor.Plan {
		for colIndex := 0; colIndex < len(row); colIndex++ {
			if row[colIndex].TileType != tileChars[UNUSED] {
				if colIndex < left {
					left = colIndex
				}

				if foundTile == false {
					top = rowIndex
				}
				bottom = rowIndex
				foundTile = true
				break
			}
		}

		for colIndex := len(row) - 1; colIndex >= 0; colIndex-- {
			if row[colIndex].TileType != tileChars[UNUSED] {
				if colIndex > right {
					right = colIndex
				}
			}
		}
	}

	newColTotal := right - left
	newRowTotal := bottom - top
	fmt.Println("right", right, "left", left)
	fmt.Println("top", top, "bottom", bottom)
	fmt.Println(newColTotal)
	fmt.Println(newRowTotal)
	newPlan := make([][]Tile, newRowTotal)

	for planIndex, _ := range newPlan {
		newPlan[planIndex] = make([]Tile, newColTotal)
	}
	i := 0
	j := 0

	type Tile struct {
		Coordinate                // Contains X (row) and Y (col) coordinates of tile
		Name        string        // Name of the tile an adventurer will see when they enter the room
		Description string        // Description of tile adventurer will see when they enter the room
		Area        Area          // The area the tile is located
		Room        Room          // The room the tile is located
		TileType    string        // The type of tile such as floor, wall, openDoor, etc
		Items       []interface{} // List of items or objects in the tile
		Players     []Player      // Adventurers in the tile
		Monsters    []Monster     // NPC in the tile
	}

	for floorRow := bottom; floorRow < top; floorRow++ {
		i++
		for floorCol := left; floorCol < right; floorCol++ {
			newPlan[i][j].Coordinate = floor.Plan[floorRow][floorCol].Coordinate
			newPlan[i][j].Name = floor.Plan[floorRow][floorCol].Name
			newPlan[i][j].Description = floor.Plan[floorRow][floorCol].Description
			newPlan[i][j].Area = floor.Plan[floorRow][floorCol].Area
			newPlan[i][j].Room = floor.Plan[floorRow][floorCol].Room
			newPlan[i][j].TileType = floor.Plan[floorRow][floorCol].TileType
			newPlan[i][j].Items = floor.Plan[floorRow][floorCol].Items
			newPlan[i][j].Players = floor.Plan[floorRow][floorCol].Players
			newPlan[i][j].Monsters = floor.Plan[floorRow][floorCol].Monsters
			j++
		}
	}

	floor.Plan = newPlan
}

// Returns true if room is within floor dimensions
func (floor *Floor) isRoomWithinFloorDimensions(row, col, width, length int) bool {

	if row-width < 0 || col-length < 0 {
		return false
	}

	if row+width > floor.Width || col+length > floor.Length {
		return false
	}
	return true
}

// x and y is the tile location which connects the current rooms
// to the next room. Direction will be used to check
// Returns true if a room is succesfully created
func (floor *Floor) makeRoom(row int, col int, dir Direction,
	width int, length int, distanceBetweenRooms int) bool {

	var topLeft Coordinate
	var bottomRight Coordinate
	topLeft.Level = floor.Level
	bottomRight.Level = floor.Level

	// x is top left and bottom bottomRight
	// O is where x and y coordinate is located
	//
	// x
	// .
	// .
	// . . . C . ..x
	//       O
	if dir == NORTH {

		topLeft.Row = (row - distanceBetweenRooms) - length
		topLeft.Col = col - (width / 2)
		bottomRight.Row = row - distanceBetweenRooms
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
		topLeft.Col = col + distanceBetweenRooms
		bottomRight.Row = row + (length / 2)
		bottomRight.Col = col + distanceBetweenRooms + width

		//       O
		// . . . C . ..x
		// .
		// x
		// .
		// .
	} else if dir == SOUTH {

		topLeft.Row = row - distanceBetweenRooms
		topLeft.Col = col - (width / 2)
		bottomRight.Row = row + distanceBetweenRooms + length
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
		topLeft.Col = col - distanceBetweenRooms - width
		bottomRight.Row = row + (length / 2)
		bottomRight.Col = col - distanceBetweenRooms

	} else {
		fmt.Println("Error invalid direction for makeRoom", row, col, dir)
	}

	if floor.isRoomUsed(topLeft, bottomRight) {
		return false
	}

	return floor.createTilesInRoom(topLeft, bottomRight, getCommonTerrainType())
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

// Builds all the tiles in the room, returns true if sucessfully built
// terrainType is a tileChar such as floor, cloud, moutain, etc, except
// special terrain such as unused, door and whirlpool
func (floor *Floor) createTilesInRoom(topLeft, bottomRight Coordinate, terrainType string) bool {

	area := getRandomArea()

	var room Room

	for i := topLeft.Row; i <= bottomRight.Row; i += 1 {
		for j := topLeft.Col; j <= bottomRight.Col; j += 1 {
			if floor.isValidCoordinate(i, j) {

				var tile Tile
				tile.Area = area

				// Adds the edge tiles to the wall list of each room
				if i == topLeft.Row || i == bottomRight.Row ||
					j == topLeft.Col || j == bottomRight.Col {

					tile.Name = "Wall of " + getRandomRoomName()
					tile.Description = "A wall is here."
					tile.TileType = tileChars[WALL]
					floor.Plan[i][j].createTile(tile)
					room.Wall = append(room.Wall, floor.Plan[i][j])

				} else {

					tile.Name = getRandomRoomName()
					tile.Description = getRandomTileDescription()
					// TODO Randomly pick a TileChar but usually its a common type such as floor or trees
					// Need to make the edges walls and not override the door
					tile.TileType = terrainType
					floor.Plan[i][j].createTile(tile)
				}

				// Adds all the tiles in the room
				room.Tiles = append(room.Tiles, floor.Plan[i][j])
			}
		}
	}

	if len(room.Wall) < 8 || len(room.Tiles) < 9 {
		return false
	}

	floor.Rooms = append(floor.Rooms, room)
	return true
}
