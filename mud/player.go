package mud

import (
	"log"
	"strings"
)

type Player struct {
	Type       string // Message type
	Username   string // Go Play Chess account
	Name       string // Mud account name
	Class      string
	Race       string
	Gender     string
	Inventory  []interface{} // What the player is carrying
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

	var status string

	err := db.QueryRow(`SELECT name, class, race, gender, status, level, experience 
		FROM mud where username=?`, username).Scan(&player.Name, &player.Class, &player.Race, &player.Gender,
		&status, &player.Level, &player.Experience)
	if err != nil {
		log.Println(err)
		return
	}
}
