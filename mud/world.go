package mud

// A dungeon consists of floors which can be transversed through stairs
type World struct {
	Floors []Floor
}

var world World

func CreateWorld() {

	const (
		low       = 3
		high      = 10
		floorLow  = 3
		floorHigh = 30
	)
	numOfFloors := getRandomIntRange(low, high)

	for i := 0; i < numOfFloors; i += 1 {
		var floor Floor
		floor.Width = getRandomIntRange(floorLow, floorHigh)
		floor.Length = getRandomIntRange(floorLow, floorHigh)
		floor.initFloorTileType()
		floor.makeRooms(i)
		world.Floors[i] = floor
	}

}
