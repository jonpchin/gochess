package mud

import (
	"bufio"
	crypto "crypto/rand"
	"errors"
	"fmt"
	"log"
	"math/big"
	"math/rand"
	"os"
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

	if max < 0 {
		return 0, errors.New("maxInclusive is less then zero")
	}
	maxInt := big.NewInt(max)

	result, err := crypto.Int(crypto.Reader, maxInt)
	if err != nil {
		return 0, err
	}

	return result.Int64(), nil
}

// Generates a randomly secure integer (int64) from minInclusive to maxInclusive
// Takes in an int64
// Returns 0 for the integer if there was an error as well as the error
func secureRandomIntRange(min, max int64) (int64, error) {

	if max < 0 {
		return 0, errors.New("maxInclusive is less then zero")
	}

	maxInt := big.NewInt(max - min)

	// Sets maxInt = maxInt + min
	maxInt.Add(maxInt, big.NewInt(min))
	result, err := crypto.Int(crypto.Reader, maxInt)

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
		return NORTH
	case 1:
		return EAST
	case 2:
		return SOUTH
	case 3:
		return WEST
	default:
		fmt.Println("Invalid direction, this should be impossible")
	}
	return NORTH
}

func (floor *Floor) getRandomRoomOnFloor() Room {
	log := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
	randNum, err := secureRandomInt(int64(len(floor.Rooms) - 1))
	if err != nil {
		log.Println(err)
	}
	return floor.Rooms[randNum]
}

// Selects a random tile on the wall of a room
func (floor *Floor) getRandomTileOnWall() Tile {
	log := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

	room := floor.getRandomRoomOnFloor()
	randNum, err := secureRandomIntRange(0, int64(len(room.Tiles)-1))
	if err != nil {
		log.Println(err)
	}
	return room.Wall[randNum]
}

// Returns a random dagger name
func GetRandomDaggerName() string {

	const daggerNamePath = "mud/equipment/generated/weapons/daggers.txt"
	daggerName, err := os.Open(daggerNamePath)
	defer daggerName.Close()

	if err != nil {
		fmt.Println("random.go getRandomDaggerName 0", err)
	}

	scanner := bufio.NewScanner(daggerName)
	var counter int64
	counter = 0

	for scanner.Scan() {
		counter++
	}
	maxNum, err := secureRandomInt(counter)

	if err != nil {
		fmt.Println("random.go getRandomDaggerName 0", err)
	}

	_, err = daggerName.Seek(0, 0)
	if err != nil {
		fmt.Println(err)
	}

	scanner = bufio.NewScanner(daggerName)
	counter = 0
	dagger := ""

	for scanner.Scan() {
		counter++
		if counter == maxNum {

			dagger = scanner.Text()
		}
	}
	return dagger
}
