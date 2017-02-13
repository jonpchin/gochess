package gostuff

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/icrowley/fake"
)

// generates a text file with fake data to be used to import into travis
// this should be used as a one time call
func FakeDataForTravis() {

	const (
		// number of rows in each table to populate data with
		rows = 20
	)
	populateUserInfo(rows)
}

// populates userinfo table in travis
// rows is the number of entries in the table
func populateUserInfo(rows int) {

	const path = "_travis/data/database/userinfo.txt"
	file, err := os.Create(path)
	if err != nil {
		fmt.Println("Cannot create file fakeTravis.go populateUserInfo 0", err)
		return
	}

	// defer is evaluated LIFO so reverse order is used
	defer os.Rename(path, "_travis/data/database/userinfo.csv")
	defer file.Close()

	var row string

	for i := 0; i < rows; i++ {

		//reset row to blank string
		row = ""

		// makes sure username is between 3 to 12 characters
		username := fake.UserName()
		if len(username) < 3 {
			username += "tes"
		}
		if len(username) < 12 {
			row = row + username + ","
		} else {
			row = row + username[:12] + ","
		}

		// atLeast, atMost, allowUpper, allowNumeric, allowSpecial
		row = row + fake.Password(5, 32, true, true, false) + ","
		row = row + fake.EmailAddress() + ","
		// Date
		row = row + time.Now().String() + ","
		// IP address
		row = row + fake.IPv4() + ","

		//default to yes they are activated, will still test
		//account activiation by setting one to no
		row = row + "No,"
		// set to zero to signify no incorect captcha entered yet
		row = row + "0,"
		// default to us country code
		row = row + "us,"

		// remove all leading and trailing whitespace
		row = strings.TrimSpace(row)
		// removes the extra comma at the end if there is one
		row = strings.TrimSuffix(row, ",")

		// write out csv into file
		fmt.Fprintln(file, row)
	}
}
