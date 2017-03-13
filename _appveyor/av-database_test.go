package gostuff

import (
	"database/sql"
	"testing"

	"github.com/jonpchin/gochess/gostuff"
)

// App Veyor default MySQL username and pass is public information
func TestAppVeyorConnect(t *testing.T) {

	db := gostuff.GetDb()

	// make sure MySQL connection is alive before proceeding
	if gostuff.CheckDBConnection("data/dbapp-veyor.txt") == false {
		t.Fatal("Failed to connect to MySQL in App Veyor")
	}
	dbString, _ := gostuff.ReadFile("data/dbapp-veyor.txt")
	var err error
	db, err = sql.Open("mysql", dbString)
	defer db.Close()

	if err != nil {
		t.Fatal("Can't open MySQL")
	}

	//if database ping fails here that means connection is alive but database is missing
	if db.Ping() != nil {
		t.Fatal("Can't ping MySQL")
	}
}
