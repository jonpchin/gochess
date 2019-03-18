package mud

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

// Default location new players start is their home base
var HOME_BASE = Coordinate{Row: 47, Col: 35, Level: 0}

// Same as Player but without sensitive sessionID
// This struct will be used to send to other players
type PlayerPublic struct {
	Type       string // Message type
	Name       string // Mud account name
	Class      string
	Race       string
	Gender     string   // M or F
	Inventory  []string // What the player is carrying
	Equipment  Equipment
	Stats      PlayerStats
	Status     []string // List of afflictions or buffs affecting player
	Bleed      int      // Amount of health the player will lose every tick
	Map        string   // Shows the section of the map where a player's vision is limited too
	Level      int
	Experience int
	Location   Coordinate
	Area       Area
	Vision     int
}

// Map will always be square and limited to odd number so MapVision 5 is 11x11, MapVision 4 is 9x9
// Largest allowed MapVision will be 10 or 21x21 unless provided by a class skill, default is MapVision 5
type Player struct {
	Type       string // Message type
	Username   string // Go Play Chess account
	SessionID  string // Encrypted session token
	Name       string // Mud account name
	Class      string
	Race       string
	Gender     string   // M or F
	Inventory  []string // What the player is carrying
	Equipment  Equipment
	Stats      PlayerStats
	Status     []string // List of afflictions or buffs affecting player
	Bleed      int      // Amount of health the player will lose every tick
	Map        string   // Shows the section of the map where a player's vision is limited too
	Level      int
	Experience int
	Location   Coordinate
	Area       Area
	Vision     int
	Tile       Tile // The tile the adventurer is currently in
}

// Contains string of MapVision of player and the details of the room he is in
// Used when player "looks"
type PlayerMap struct {
	Type        string
	Map         string
	Coordinates Coordinate
	CurrentTile Tile
}

type PlayerStats struct {
	Name         string
	Health       int // Health points
	Mana         int // Mana points (skills that require magic)
	Energy       int // Move points (skils that require physical effort)
	Strength     int // Multiplier for physical attacks
	Speed        int // Mhance to dodge, block and parry physical attacks
	Dexterity    int // Multiplier to regain balance
	Intelligence int // Multipier to regain mental stabliity
	Wisdom       int // Multiplier for magic attacks
}

// Default actions for all players
type PlayerActions interface {
	kill(string)
	say(string)
	look(string)
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

func (player *Player) loadMap() {

	// only show a portion of the map to the user
	//mapView := ""
	// default view for map is 5
	//viewDistance := 5

	// Make sure the world is already set
	player.setMapVision(world)
}

func (player *Player) loadPlayerData() {

	// status will be a csv which will be stored in an array
	var status string
	log := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

	err := db.QueryRow(`SELECT name, class, race, gender, status, level, experience 
		FROM mud where username=?`, player.Username).Scan(&player.Name, &player.Class, &player.Race, &player.Gender,
		&status, &player.Level, &player.Experience)
	if err != nil {
		log.Println(err)
		return
	}

	err = db.QueryRow(`SELECT area, x, y, z FROM location where name=?`,
		player.Name).Scan(&player.Area.Name, &player.Location.Row, &player.Location.Col,
		&player.Location.Level)
	if err != nil {
		log.Println(err)
		return
	}
}

// Sets the player map to be sent to the client
func (player *Player) setPlayerMap() {

	log := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

	err := db.QueryRow(`SELECT x, y, z FROM location where name=?`,
		player.Username).Scan(&player.Location.Row, &player.Location.Col,
		&player.Location.Level)
	if err != nil {
		log.Println(err)
		return
	}

	//playerMap.CurrentTile = world.Floors[playerMap.Coordinates.Level].Plan[playerMap.Coordinates.Row][playerMap.Coordinates.Col]
}

func (player *Player) save() {

	log := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
	stmt, err := db.Prepare("UPDATE mud SET name=?, class=?, race=?, gender=?, status=?, level=?, experience=? WHERE username=?")
	if err != nil {
		log.Println(err)
	}
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
		log.Println(err)
	}

	stmt, err = db.Prepare("UPDATE location SET area=?, x=?, y=?, z=?, map=?, WHERE name=?")
	if err != nil {
		log.Println(err)
	}

	_, err = stmt.Exec(&player.Area.Name, &player.Location.Row, &player.Location.Col, &player.Location.Level,
		&player.Map, &player.Name)
	if err != nil {
		log.Println(err)
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

// Ensure player is not trying to impersonate someone else or change name without permission
// Returns true if player's username and sessionID are all valid
func (player *Player) isCredValid(username, sessionID string) bool {
	log := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

	if player.Username != username {
		log.Println("player.Username:", player.Username, "username:", username)
		return false
	} else if player.SessionID != sessionID {
		log.Println("player.SessionID:", player.SessionID, "sessionID:", sessionID)
		return false
	}
	return true
}

// First time logging only compares username and sessionID
func (player *Player) isCredValidFirstTime(username, sessionID string) bool {
	log := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
	if player.Username != username {
		log.Println("player.Username:", player.Username, "username:", username)
		return false
	} else if player.SessionID != sessionID {
		log.Println("player.SessionID:", player.SessionID, "sessionID:", sessionID)
		return false
	}
	return true
}

/*
func (equipment *Equipment) UnmarshalJSON(data []byte) error {

	var v []interface{}
    if err := json.Unmarshal(data, &v); err != nil {
        fmt.Printf("Error whilde decoding %v\n", err)
        return err
    }
    tp.Timestamp = int64(v[0].(float64))
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
