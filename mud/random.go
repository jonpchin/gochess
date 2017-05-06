package mud

import (
	crypto "crypto/rand"
	"io"
	"math/big"
	"math/rand"
	"time"
)

// Calling getRandomInt(100) will return a random number 0 to 100 inclusive
// If maxInclusive is less then zero, then zero will be returned
func getRandomInt(maxInclusive int) int {
	if maxInclusive >= 0 {
		return rand.Intn(maxInclusive)
	}
	return 0
}

func getRandomIntRange(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}

// Generates a randomly secure integer (int64) from 0 to maxInclusive
// Takes in an int64
// Returns 0 for the integer if there was an error as well as the error
func secureRandomInt(maxInclusive int64) (int64, error) {
	var reader io.Reader
	maxInt := big.NewInt(maxInclusive)
	result, err := crypto.Int(reader, maxInt)
	if err != nil {
		return 0, err
	}
	return result.Int64(), nil
}
