package mud

import (
	"database/sql"
	"log"
	"os"
	"strings"

	"github.com/jonpchin/gochess/gostuff"
)

var db *sql.DB

func ConnectDb() {
	db = gostuff.GetDb()
}

// Checks if Go Play Chess username is present in MUD table if not then add it
func lookupName(username string) error {
	var name string
	//checking if name exists
	err := db.QueryRow("SELECT username FROM mud WHERE username=?", username).Scan(&name)

	if err != nil {
		return err
	}

	return nil
}

// Checks if a player has an adventurer name
// Returns true if MUD name exists
func isNameExistForPlayer(username string) bool {
	var name string
	//checking if name exists
	_ = db.QueryRow("SELECT name FROM mud WHERE username=?", username).Scan(&name)
	if name == "" {
		return false
	}
	return true
}

// Checks if name is already isNameTaken
// Returns true if name is already taken
func isNameTaken(name string) bool {
	var result string
	_ = db.QueryRow("SELECT name FROM mud WHERE name=?", name).Scan(&result)
	if result == "" {
		return false
	}
	return true
}

// Update player into database
func (player *Player) registerPlayer() {

	log := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

	status := ""

	for _, value := range player.Status {
		status = status + "," + value
	}

	// Removing leading and trailing commas
	status = strings.Trim(status, ",")

	stmt, err := db.Prepare("INSERT mud SET username=?, name=?, class=?, race=?, gender=?, status=?, level=?, experience=?")
	defer stmt.Close()

	_, err = stmt.Exec(player.Username, player.Name, player.Class, player.Race, player.Gender, status, player.Level, player.Experience)
	if err != nil {
		log.Println(err)
	}

	stmt, err = db.Prepare("INSERT location SET name=?, area=?, x=?, y=?, z=?")
	if err != nil {
		log.Println(err)
	}

	stmt, err = db.Prepare("INSERT equipment SET name=?, weapon=?, sidearm=?, shield=?, helmet=?, torso=?, belt=?, arms=?, legs=?, shoes=?, ring=?, floating=?")
	if err != nil {
		log.Println(err)
	}

	// All new players get some newbie gear
	_, err = stmt.Exec(player.Name, "None", "None", "None", "Newbie Helm", "Newbie Mail", "Newbie Belt", "None", "Newbie Pants", "Newbie Boots", "None", "None")
	if err != nil {
		log.Println(err)
	}

	_, err = stmt.Exec(player.Name, player.Area.Name, player.Location.Row, player.Location.Col, player.Location.Level)
	if err != nil {
		log.Println(err)
	}
}
