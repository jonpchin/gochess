package mud

// A floor is a 2D plane which consists of many rooms connected to each other
// A floor should have at least one stairway or portal leading to another floor
type Floor struct {
	Width  int          // Number of tiles wide
	Length int          // Number of tiles vertically
	Rooms  []Room       // List of rooms on the floor in no particular order
	Plan   [][]TileChar // 2D ASCI wilderness map
}

// ----------> Width
// |
// |
// |
// V Length
// Initializes floor unused tile characters
func (floor *Floor) initFloor() {
	for i := 0; i < floor.Length; i += 1 {
		for j := 0; j < floor.Width; j += 1 {
			floor.Plan[i][j] = unused
		}
	}
}

func generateFloor() {

}
