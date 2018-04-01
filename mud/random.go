package mud

import (
	"bufio"
	crypto "crypto/rand"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"time"
)

// Used to keep track of the currently used line in a file for generating random data
type LineTracker struct {
	totalLines  int
	currentLine int
}

func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

// Calling getRandomInt(100) will return a random number 0 to 100 inclusive
// If max is less then zero, then zero will be returned
func getRandomInt(max int) int {
	rand.Seed(time.Now().UnixNano())
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

	if max <= 0 {
		return 0, nil
	}
	maxInt := big.NewInt(max)

	result, err := crypto.Int(crypto.Reader, maxInt)
	if err != nil {
		return 0, err
	}

	return result.Int64(), nil
}

// Generates a randomly secure integer (int64) from minInclusive to maxInclusive
// Takes in an int64, must be greater then zero or panic will occur
// Returns 0 for the integer if there was an error as well as the error
func secureRandomIntRange(min, max int64) (int64, error) {

	if max <= 0 {
		return 0, nil
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

	rand.Seed(time.Now().UnixNano())
	result := rand.Intn(4)

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

	// Crypto int cannot take a max int of zero
	rand.Seed(time.Now().UnixNano())
	randNum := rand.Intn(len(floor.Rooms))

	return floor.Rooms[randNum]
}

// Selects a random tile on the wall of a room
func (room *Room) getRandomTileOnWall() Tile {

	rand.Seed(time.Now().UnixNano())
	randNum := rand.Intn(len(room.Wall))
	return room.Wall[randNum]
}

func (lineTracker LineTracker) getRandomItemFromPath(file *os.File) string {
	scanner := bufio.NewScanner(file)
	counter := 0
	item := ""

	for scanner.Scan() {
		item = scanner.Text()
		if counter == lineTracker.currentLine {
			break
		}
		counter++
	}
	return item
}

func getRandomItemFromPath(file *os.File) string {

	scanner := bufio.NewScanner(file)
	var counter int
	counter = 0

	for scanner.Scan() {
		counter++
	}
	maxNum := getRandomInt(counter)

	_, err := file.Seek(0, 0)
	if err != nil {
		fmt.Println(err)
	}

	scanner = bufio.NewScanner(file)
	counter = 0
	item := ""

	for scanner.Scan() {
		counter++
		if counter == maxNum {

			item = scanner.Text()
		}
	}
	return item
}

func getRandomRoomName() string {
	const roomPath = "mud/story/names.txt"
	room, err := os.Open(roomPath)
	defer room.Close()

	if err != nil {
		fmt.Println("random.go GetRandomRoomName 0", err)
	}
	return getRandomItemFromPath(room)
}

func getRandomTileDescription() string {
	const descriptionPath = "mud/story/descriptions.txt"
	description, err := os.Open(descriptionPath)
	defer description.Close()

	if err != nil {
		fmt.Println("random.go getRandomTileDescription 0", err)
	}
	return getRandomItemFromPath(description)
}
func getRandomArea() Area {
	var area Area
	const areaPath = "mud/story/areas.txt"
	areaTemp, err := os.Open(areaPath)
	defer areaTemp.Close()

	if err != nil {
		fmt.Println("random.go getRandomArea 0", err)
	}

	area.Name = getRandomItemFromPath(areaTemp)
	return area // Replace this later
}

func getRandomTileChar() string {

	rand.Seed(time.Now().UnixNano())
	randNum := rand.Intn(len(tileChars) - 1)
	return tileChars[randNum]
}

func getCommonTerrainType() string {

	rand.Seed(time.Now().UnixNano())
	randNum := rand.Intn(len(commonTerrainTypes) - 1)
	return commonTerrainTypes[randNum]
}

// Returns a random dagger name
func getRandomDaggerName() string {

	const daggerPath = "mud/equipment/generated/weapons/daggers.txt"
	dagger, err := os.Open(daggerPath)
	defer dagger.Close()

	if err != nil {
		fmt.Println("random.go getRandomDaggerName 0", err)
	}
	return getRandomItemFromPath(dagger)
}

func getRandomBeltsName() string {

	const beltsPath = "mud/equipment/generated/armor/belts.txt"
	belts, err := os.Open(beltsPath)
	defer belts.Close()

	if err != nil {
		fmt.Println("random.go getRandomBeltsName 0", err)
	}
	return getRandomItemFromPath(belts)
}

func getRandomBootsName() string {

	const bootsPath = "mud/equipment/generated/armor/boots.txt"
	boots, err := os.Open(bootsPath)
	defer boots.Close()

	if err != nil {
		fmt.Println("random.go getRandomBeltsName 0", err)
	}
	return getRandomItemFromPath(boots)
}

func getRandomLegsName() string {

	const legsPath = "mud/equipment/generated/armor/legs.txt"
	legs, err := os.Open(legsPath)
	defer legs.Close()

	if err != nil {
		fmt.Println("random.go getRandomLegsName 0", err)
	}
	return getRandomItemFromPath(legs)
}

func getRandomShieldsName() string {

	const shieldsPath = "mud/equipment/generated/armor/shields.txt"
	shields, err := os.Open(shieldsPath)
	defer shields.Close()

	if err != nil {
		fmt.Println("random.go getRandomShieldsName 0", err)
	}
	return getRandomItemFromPath(shields)
}

func getRandomTorsosName() string {

	const torsoPath = "mud/equipment/generated/armor/torso.txt"
	torso, err := os.Open(torsoPath)
	defer torso.Close()

	if err != nil {
		fmt.Println("random.go getRandomTorsosName 0", err)
	}
	return getRandomItemFromPath(torso)
}
