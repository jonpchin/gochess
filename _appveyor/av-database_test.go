package travis

import (
	"database/sql"
	"testing"

	"github.com/jonpchin/gochess/gostuff"
)

var db *sql.DB

// App Veyor default MySQL username and pass is public information
func TestDbConnect(t *testing.T) {
	// make sure MySQL connection is alive before proceeding
	if gostuff.CheckDBConnection("data/dbapp-veyor.txt") == false {
		t.Fatal("Failed to connect to MySQL in App Veyor")
	}
	dbString, _ := gostuff.ReadFile("data/dbapp-veyor.txt")
	db, err := sql.Open("mysql", dbString)
	defer db.Close()
	//	db.SetMaxIdleConns(20)
	if err != nil {
		t.Fatal("Can't open MySQL")
	}

	//if database ping fails here that means connection is alive but database is missing
	if db.Ping() != nil {
		t.Fatal("Can't ping MySQL")
	}
}
