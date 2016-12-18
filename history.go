package gostuff

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// contains the date-time of the rating
type RatingDate struct {
	DateTime string
	Rating   float64
}

// fetches rating history, unmarshals it, adds a new rating history, then marshals data and then
// store it back in the database returns true if sucessfully updates rating history with no errors
func updateRatingHistory(name string, gametype string, rating float64) bool {

	problems, _ := os.OpenFile("logs/errors.txt", os.O_APPEND|os.O_WRONLY, 0666)
	defer problems.Close()
	log := log.New(problems, "", log.LstdFlags|log.Lshortfile)

	//check if database connection is open
	if db.Ping() != nil {
		log.Println("DATABASE DOWN!")
		return false
	}

	var ratingHistory string
	flag := true

	// TODO: Replace SQL string concatention, gametype is already filtered to only bullet
	// blitz or standard but still should attempt to avoid SQL injection,
	// parameterized SQL statements put quotations and gametype can't have quotations to work properly
	err := db.QueryRow("SELECT "+gametype+" FROM ratinghistory WHERE username=?", name).Scan(&ratingHistory)
	if err == sql.ErrNoRows { // this will occur if there is no name exist
		log.Println("No name found in ratinghistory for ", name)
		return false
	} else if ratingHistory == "" { // Then there is no history so insert instead of update
		flag = false
	} else if err != nil {
		log.Println(err)
		return false
	}

	// used to append current game rating history into rating history memory
	var ratingHistoryMemory []RatingDate

	if flag { // when there is no rating history do not need to unmarshal data from database as there is none

		//unmarshall JSON string into ratingHistoryMemory which is a memory model
		if err := json.Unmarshal([]byte(ratingHistory), &ratingHistoryMemory); err != nil {
			fmt.Println("Just receieved a message I couldn't decode:", ratingHistory, "test", err)
			return false
		}
	}
	// TODO: Need to add a corner case if user is not in history table
	var ratingInfo RatingDate
	ratingInfo.DateTime = time.Now().Format("20060102150405")
	ratingInfo.Rating = rating
	ratingHistoryMemory = append(ratingHistoryMemory, ratingInfo)

	// need to marshall memory model before storing in database
	updatedRatingHistory, err := json.Marshal(ratingHistoryMemory)
	if err != nil {
		fmt.Println("updateRatingHistory problem marshalling ", err)
		return false
	}

	stmt, err := db.Prepare("UPDATE ratinghistory SET " + gametype + "=? WHERE username=?")
	if err != nil {
		log.Println(err)
		return false
	}

	_, err = stmt.Exec(string(updatedRatingHistory), name)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}
