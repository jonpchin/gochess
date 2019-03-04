package mud

import "fmt"

// From a starting tile create a corridor in a certain direction of a certain length
func (floor *Floor) createCorridorInDirection(tile Tile, dir Direction, length int) {

	if dir == NORTH {
		i := 0
		for i = 1; i < length; i++ {
			floor.Plan[tile.Row-i][tile.Col].makeCorridor()
			floor.Plan[tile.Row-i][tile.Col].Row = tile.Row - i
			floor.Plan[tile.Row-i][tile.Col].Col = tile.Col
		}
		floor.Plan[tile.Row-i][tile.Col].createWallFeature()

	} else if dir == EAST {
		i := 0
		for i = 1; i < length; i++ {
			floor.Plan[tile.Row][tile.Col+i].makeCorridor()
			floor.Plan[tile.Row][tile.Col+i].Row = tile.Row
			floor.Plan[tile.Row][tile.Col+i].Col = tile.Col + i
		}
		floor.Plan[tile.Row][tile.Col+i].createWallFeature()
	} else if dir == SOUTH {
		i := 0
		for i = 1; i < length; i++ {
			floor.Plan[tile.Row+i][tile.Col].makeCorridor()
			floor.Plan[tile.Row+i][tile.Col].Row = tile.Row + i
			floor.Plan[tile.Row+i][tile.Col].Col = tile.Col
		}
		floor.Plan[tile.Row+i][tile.Col].createWallFeature()
	} else if dir == WEST {
		i := 0
		for i = 1; i < length; i++ {
			floor.Plan[tile.Row][tile.Col-i].makeCorridor()
			floor.Plan[tile.Row][tile.Col-i].Row = tile.Row
			floor.Plan[tile.Row][tile.Col-i].Col = tile.Col - i
		}
		floor.Plan[tile.Row][tile.Col-i].createWallFeature()
	} else {
		fmt.Println("Impossible corridor direction")
	}
}

// Fills out basic meta data for corridor
func (tile *Tile) makeCorridor() {

	tile.Name = "Corridor"
	tile.Description = "A long winding corridor looms before you."
	tile.TileType = tileChars[CORRIDOR]
}
