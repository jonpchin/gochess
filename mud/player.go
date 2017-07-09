package mud

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

// Default location new players start is their home base
var HOME_BASE = Coordinate{Row: 10, Col: 10, Level: 5}

type Player struct {
	Type       string // Message type
	Username   string // Go Play Chess account
	Name       string // Mud account name
	Class      string
	Race       string
	Gender     string
	Inventory  []string // What the player is carrying
	Equipment  Equipment
	Stats      PlayerStats
	Status     []string // List of afflictions or buffs affecting player
	Bleed      int      // Amount of health the player will lose every tick
	Level      int
	Experience int
	Location   Coordinate
	Area       Area
}

type PlayerStats struct {
	Name         string
	Health       int
	Mana         int
	Energy       int
	Strength     int
	Speed        int
	Dexterity    int
	Intelligence int
	Wisdom       int
}

// Checks if class the person entered is a substring of a valid class and returns the full class name
// If not a valid class it will return blank string
func isValidClass(class string) (bool, string) {
	var classes = []string{"warrior", "barbarian", "monk", "mage", "thief", "ranger", "swordmaster", "illusionist",
		"priest", "necromancer", "witch", "paladin", "alchemist", "jester"}
	class = strings.ToLower(class)
	for _, value := range classes {
		if strings.Contains(class, value) {
			return true, value
		}
	}
	return false, ""
}

func (player *Player) loadPlayerData(username string) {

	// status will be a csv which will be stored in an array
	var status string

	err := db.QueryRow(`SELECT name, class, race, gender, status, level, experience 
		FROM mud where username=?`, username).Scan(&player.Name, &player.Class, &player.Race, &player.Gender,
		&status, &player.Level, &player.Experience)
	if err != nil {
		log.Println(err)
		return
	}
}

func (player *Player) save() {

	stmt, err := db.Prepare("UPDATE mud SET name=?, class=?, race=?, gender=?, status=?, level=?, experience=? WHERE username=?")
	defer stmt.Close()

	status := ""

	for _, value := range player.Status {
		status = status + "," + value
	}

	// Removing leading and trailing commas
	strings.Trim(status, ",")

	_, err = stmt.Exec(&player.Name, &player.Class, &player.Race, &player.Gender,
		&status, &player.Level, &player.Experience, &player.Username)
	if err != nil {
		fmt.Println("save 1", err)
	}
}

// Update player stats based on race and class
func (player *Player) updateByRaceClass(jsonFile string) {
	log := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
	message, err := ioutil.ReadFile(jsonFile)
	fmt.Println(message)
	if err != nil {
		log.Println(err)
	}
}

/*
func (equipment *Equipment) UnmarshalJSON(data []byte) error {

	var v []interface{}
    if err := json.Unmarshal(data, &v); err != nil {
        fmt.Printf("Error whilde decoding %v\n", err)
        return err
    }
    equipment.Timestamp = int64(v[0].(float64))
    tp.Latitude, _ = strconv.ParseFloat(v[1].(string), 64)
    tp.Longitude, _ = strconv.ParseFloat(v[2].(string), 64)
    tp.Altitude = int(v[3].(float64))
    tp.Value1 = v[4].(float64)
    tp.Value2 = int16(v[5].(float64))
    tp.Value3 = int16(v[6].(float64))

    return nil
}

func (playerStats *PlayerStats) UnmarshalJSON(data []byte) error {
}

func (coordinate *Coordinate) UnmarshalJSON(data []byte) error {
}

func (area *Area) UnmarshalJSON(data []byte) error {
}
*/
