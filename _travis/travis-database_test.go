package travis

import (
	"database/sql"
	"os/exec"
	"testing"

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

	err, result := importDbIntoTravis()
	if err != nil {
		t.Fatal("Can't import database", err, "result is ", result)
	}
	/*
		err = importTablesIntoTravis()
		if err != nil {
			t.Fatal("Error importing tables into travis", err)
		}
	*/
	/*
		var userInfo UserInfo
		userInfo.username = "jon"
		userInfo.password = "test"
		userInfo.email = "fake@email.com"
		userinfo.ipAddress = "1.1.1.1"
		userInfo.random

		success := userInfo.register(w, r)
		if success {

		}
	*/
}

// imports template database into travis, returns error if there was one
func importDbIntoTravis() (error, string) {

	result, err := exec.Command("/bin/bash", "-c", "cd data && bash importTravisTemplate.sh").Output()
	if err != nil {
		return err, string(result)
	}
	return nil, ""
}

// imports fake data into tables on travis, returns error if there was one
func importTablesIntoTravis() error {
	_, err := exec.Command("/bin/bash", "-c", "cd data && bash importTravisTables.sh").Output()
	if err != nil {
		return err
	}
	return nil
}
