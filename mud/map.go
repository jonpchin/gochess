package mud

import "math"

// Converts map string to 2D array of string
func convertMapToTiles(mapString string) [][]string {

	length := math.Sqrt(float64(len(mapString)))
	intLength := int(length)
	var tiles = make([][]string, intLength)

	var temp []string

	for _, r := range mapString {
		temp = append(temp, string(r))
	}

	for i := 0; i < intLength; i += 1 {
		tiles[i] = make([]string, intLength)

		for j := 0; j < intLength; j += 1 {
			tiles[i][j] = temp[(i*intLength)+j]
		}
	}
	return tiles
}
