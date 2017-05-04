package mud

// A floor is a 2D plane which consists of many rooms connected to each other
// A floor should have at least one stairway or portal leading to another floor
type Floor struct {
	width  int // x
	length int // y
	rooms  []Room
}

func generateFloor() {

}
