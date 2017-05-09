package mud

import (
	crypto "crypto/rand"
	"errors"
	"fmt"
	"io"
	"math/big"
	"math/rand"
	"time"
)

// Calling getRandomInt(100) will return a random number 0 to 100 inclusive
// If max is less then zero, then zero will be returned
func getRandomInt(max int) int {
	if max >= 0 {
		return rand.Intn(max)
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
func secureRandomInt(max int64) (int64, error) {
	var reader io.Reader
	if max < 0 {
		return 0, errors.New("maxInclusive is less then zero")
	}
	maxInt := big.NewInt(max)
	result, err := crypto.Int(reader, maxInt)
	if err != nil {
		return 0, err
	}
	return result.Int64(), nil
}

// Generates a randomly secure integer (int64) from minInclusive to maxInclusive
// Takes in an int64
// Returns 0 for the integer if there was an error as well as the error
func secureRandomIntRange(min, max int64) (int64, error) {
	var reader io.Reader
	if max < 0 {
		return 0, errors.New("maxInclusive is less then zero")
	}

	maxInt := big.NewInt(max - min)

	// Sets maxInt = maxInt + min
	maxInt.Add(maxInt, big.NewInt(min))
	result, err := crypto.Int(reader, maxInt)

	if err != nil {
		return 0, err
	}
	return result.Int64(), nil
}

// Returns a radom direction. Returns north if there was an error
func getRandomDirection() Direction {
	result, err := secureRandomInt(3)
	if err != nil {
		fmt.Println("Can't get random integer for direction", err)
	}
	switch result {
	case 0:
		return north
	case 1:
		return east
	case 2:
		return south
	case 3:
		return west
	default:
		fmt.Println("Invalid direction, this should be impossible")
	}
	return north
}
