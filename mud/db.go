package mud

import (
	"database/sql"
	"strings"

	"github.com/jonpchin/gochess/gostuff"
)

var db *sql.DB

func ConnectMudDb() {
	db = gostuff.GetDb()
}

// Looks up the name MUD based on goplaychess username
// First checks username to see if it exist, if not it prompts for name
func lookupName(username string) {

}

// Returns true if MUD name exists
func checkNameExist(username string) bool {
	var name string
	//checking if name exists
	_ = db.QueryRow("SELECT username FROM mud WHERE username=?", username).Scan(&name)
	if strings.EqualFold(username, name) { // already exists, case insensitive comparison
		return true
	} else {
		return false
	}
}

func promptNewname() {
}
