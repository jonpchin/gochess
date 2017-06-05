package mud

import (
	"database/sql"
	"fmt"

	"github.com/jonpchin/gochess/gostuff"
)

var db *sql.DB

func ConnectMudDb() {
	db = gostuff.GetDb()
}

// Checks if Go Play Chess username is present in MUD table if not then add it
func lookupName(username string) {
	var name string
	//checking if name exists
	_ = db.QueryRow("SELECT username FROM mud WHERE username=?", username).Scan(&name)
	if name == "" { // already exists, case insensitive comparison
		registerUsername(username)
	} else {
		fmt.Println("Username is already present for", username)
	}
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

func registerName(name string, username string) {

	stmt, err := db.Prepare("UPDATE mud SET name=? WHERE username=?")
	defer stmt.Close()

	_, err = stmt.Exec(name, username)
	if err != nil {
		fmt.Println(err)
	}
}

// If a player username does not exist then register it
func registerUsername(username string) {

	stmt, err := db.Prepare(`INSERT INTO mud (username, name, class, race, gender, status, level, experience) 
		VALUES (?, "", "", "", "", "", 0, 0)`)
	defer stmt.Close()

	_, err = stmt.Exec(username)
	if err != nil {
		fmt.Println(err)
	}
}

func registerClass(class string, username string) bool {
	valid, fullClass := isValidClass(class)
	if valid {
		stmt, err := db.Prepare("UPDATE mud SET class=? WHERE username=?")
		defer stmt.Close()

		_, err = stmt.Exec(fullClass, username)
		if err != nil {
			fmt.Println(err)
		}
		return true
	}
	return false
}
