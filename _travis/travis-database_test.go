package gostuff

import (
	"database/sql"
	"testing"

	"github.com/icrowley/fake"
	"github.com/jonpchin/gochess/gostuff"
)

var db *sql.DB

// Travis CI default MySQL username and pass is public information
func TestTravisConnect(t *testing.T) {

	// make sure MySQL connection is alive before proceeding
	if gostuff.CheckDBConnection("data/dbtravis.txt") == false {
		t.Fatal("Failed to connect to MySQL in Travis CI")
	}

	var err error
	dbString, _ := gostuff.ReadFile("data/dbtravis.txt")
	db, err = sql.Open("mysql", dbString)
	//defer db.Close()

	if err != nil {
		t.Fatal("Can't open MySQL")
	}

	gostuff.SetDb(db)

	//if database ping fails here that means connection is alive but database is missing
	if db.Ping() != nil {
		t.Fatal("Can't ping MySQL")
	}

	// registers a random person to the database
	var userInfo gostuff.UserInfo
	userInfo.Username = fake.UserName()

	// Ensure username is between 3 and 12 characters
	if len(userInfo.Username) < 3 {
		userInfo.Username += "tes"
	} else if len(username) > 12 {
		userInfo.Username = userInfo.Username[:12]
	}

	userInfo.Password = fake.Password(5, 32, true, true, false)

	userInfo.Email = fake.EmailAddress()
	userInfo.IpAddress = fake.IPv4()

	// doesnt matter what parameter as its only for handling corner case in localhost
	_, err = userInfo.Register("")
	if err != nil {
		t.Fatal(err)
	}

	found := gostuff.CheckUserNameInDb("test1234")
	if found {
		t.Fatal("Username was found in the database when it was not suppose to be")
	}
	found = gostuff.CheckUserNameInDb(userInfo.Username)
	if found == false {
		t.Fatal("Username was not found  in the database when it was suppose to be")
	}
}
