package travis

import (
	"database/sql"
	"net/http"
	"testing"

	"github.com/icrowley/fake"
	"github.com/jonpchin/gochess/gostuff"
)

var db *sql.DB

// Travis CI default MySQL username and pass is public information
func TestDbConnect(t *testing.T) {

	// make sure MySQL connection is alive before proceeding
	if gostuff.CheckDBConnection("data/dbtravis.txt") == false {
		t.Fatal("Failed to connect to MySQL in Travis CI")
	}
	dbString, _ := gostuff.ReadFile("data/dbtravis.txt")
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

	// registers a random person to the database
	var userInfo gostuff.UserInfo
	userInfo.Username = fake.UserName()
	userInfo.Password = fake.Password(5, 32, true, true, false)

	userInfo.Email = fake.EmailAddress()
	userInfo.IpAddress = fake.IPv4()

	var w http.ResponseWriter
	var r *http.Request
	err = userInfo.Register(w, r)
	if err != nil {
		t.Fatal(err)
	}

	found := gostuff.CheckUserNameInDb(userInfo.Username)
	if found == false {
		t.Fatal("Username was not found in the database")
	}
}
