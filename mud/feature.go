package mud

import (
	"math/rand"
	"time"
)

// These types of terrain are most common
var commonWallTypes = []string{
	tileChars[CLOSEDOOR],
	tileChars[OPENDOOR],
}

// Creates a feature on a wall
func (wall *Tile) createWallFeature() {

	rand.Seed(time.Now().UnixNano())
	randNum := rand.Intn(len(commonWallTypes))
	wall.TileType = commonWallTypes[randNum]
	//TODO: Have function to specify meta data for the wall
}
