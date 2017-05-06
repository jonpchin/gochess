package mud

// A dungeon consists of floors which can be transversed through stairs
type World struct {
	Floors []Floor
}

var world World

func CreateWorld() {
	for i := 0; i < 10; i += 1 {
		var floor Floor
		floor.width = 10
		floor.length = 10
		world.Floors[i] = floor
	}

}
