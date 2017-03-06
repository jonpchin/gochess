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

	// Setting up users for chrominum driver web test
	userInfo.Username = "ben"
	userInfo.Password = "test123"
	userInfo.Email = "fake@email.com"
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

	result, err := userInfo.Login("", "", "", "localhost")
	if err != nil {
		t.Fatal(err)
	}

	const (
		needToActivateMessage = "<img src='img/ajax/not-available.png' /> " +
			"You must activate your account by entering the activation token in " +
			"your email at the activation page. An email has been sent again " +
			"containing your activation code."
	)

	if result != needToActivateMessage {
		t.Fatal("Incorrect login results", result)
	}

	_, isActivate := userInfo.Activate("", "", "")
	if isActivate == false {
		t.Fatal("Could not activate account on travis unit test", isActivate)
	}

	// Retrying login in after account activation
	result, err = userInfo.Login("", "", "", "localhost")
	if err != nil {
		t.Fatal(err)
	}

	if result != "" {
		t.Fatal("Incorrect login results 2 results should be blank string")
	}

	errMessage, bullet, blitz, standard, correspondence := gostuff.GetRating(userInfo.Username)
	if errMessage != "" {
		t.Fatal("Error in fetching rating for ", userInfo.Username, errMessage)
	}
	if blitz != 1500 || bullet != 1500 || standard != 1500 || correspondence != 1500 {
		t.Fatal("Ratings are not set to 1500", blitz, bullet, standard, correspondence)
	}

	errMessage, bulletF, blitzF, standardF, correspondenceF, bulletRD,
		blitzRD, standardRD, correspondenceRD := gostuff.GetRatingAndRD(userInfo.Username)

	if errMessage != "" {
		t.Fatal("Error in fetching rating and ratngRD for ", userInfo.Username, errMessage)
	}
	if blitzF != 1500 || bulletF != 1500 || standardF != 1500 || correspondenceF != 1500 {
		t.Fatal("Ratings part 2 are not set to 1500", blitz, bullet, standard, correspondence)
	}

	if bulletRD != 350 || blitzRD != 350 ||
		standardRD != 350 || correspondenceRD != 350 {
		t.Fatal("Rating RD float values are incorrect", bulletRD, blitzRD, standardRD, correspondenceRD)
	}

	userInfo.Username = "can"
	userInfo.Email = fake.EmailAddress()
	userInfo.IpAddress = fake.IPv4()

	// doesnt matter what parameter as its only for handling corner case in localhost
	_, err = userInfo.Register("")
	if err != nil {
		t.Fatal(err)
	}

	_, isActivate = userInfo.Activate("", "", "")
	if isActivate == false {
		t.Fatal("Could not activate account on travis unit test", isActivate)
	}
}
