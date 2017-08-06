package mud

// A floor is a 2D plane which consists of many rooms connected to each other
// A floor should have at least one stairway or portal leading to another floor
type Floor struct {
	Width  int      // Number of tiles wide
	Length int      // Number of tiles vertically
	Rooms  []Room   // List of rooms on the floor in no particular order
	Plan   [][]Tile // 2D ASCI wilderness map
	Level  int
}

// ----------> Width
// |
// |
// |
// V Length
// Initializes floor tiletype to unused tile characters
func (floor *Floor) initFloorTileType() {

	floor.Plan = make([][]Tile, floor.Length)

	for i := 0; i < floor.Length; i += 1 {
		floor.Plan[i] = make([]Tile, floor.Width)

		for j := 0; j < floor.Width; j += 1 {
			floor.Plan[i][j].TileType = tileChars[UNUSED]
			floor.Plan[i][j].Row = i
			floor.Plan[i][j].Col = j
		}
	}
}

func generateFloor() {

}
