package gostuff

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
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

	flag := true
	ratingHistory, pass, err := GetRatingHistory(name, gametype)

	if pass == false {
		log.Println(err)
		return false
	} else if ratingHistory == "" {
		flag = false
	}

	// used to append current game rating history into rating history memory
	var ratingHistoryMemory []RatingDate

	if flag { // when there is no rating history do not need to unmarshal data from database as there is none

		//unmarshall JSON string into ratingHistoryMemory which is a memory model
		if err := json.Unmarshal([]byte(ratingHistory), &ratingHistoryMemory); err != nil {
			log.Println("Just receieved a message I couldn't decode:", ratingHistory, err)
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
		log.Println("updateRatingHistory problem marshalling ", err)
		return false
	}

	return modifyRatingHistory(name, gametype, updatedRatingHistory)
}

// Updates rating history based on game type for user, returns false if there was an error
func modifyRatingHistory(name string, gametype string, updatedRatingHistory []byte) bool {
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

// Removes game history after a certain number of days
func RemoveGameHistory(days string) {

	log := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

	rows, err := db.Query("SELECT username FROM ratinghistory")
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()

	var name string
	daysConverted, err := strconv.Atoi(days)
	if err != nil {
		log.Println(err)
	}

	for rows.Next() {

		err = rows.Scan(&name)
		if err != nil {
			log.Println(err)
		}
		removeHistoryFromPlayer(name, daysConverted, "bullet")
		removeHistoryFromPlayer(name, daysConverted, "blitz")
		removeHistoryFromPlayer(name, daysConverted, "standard")
		removeHistoryFromPlayer(name, daysConverted, "correspondence")
	}
}

// Removes game history older then a certain amount of days for a player
// Returns true if no errors
func removeHistoryFromPlayer(name string, days int, gametype string) bool {

	log := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
	ratingHistory, pass, err := GetRatingHistory(name, gametype)

	if pass {
		// If there is no ratingHistory then there is nothing to remove
		if ratingHistory != "" {
			var ratingHistoryMemory []RatingDate

			if err := json.Unmarshal([]byte(ratingHistory), &ratingHistoryMemory); err != nil {
				log.Println("Just receieved a message I couldn't decode:", ratingHistory, "test", err)
				return false
			}

			hours := days * 24
			timeFormat := "20060102150405"

			for i, game := range ratingHistoryMemory {
				isElpase, _ := HasTimeElapsed(game.DateTime, hours, timeFormat, true)
				if isElpase == false { // Use function to get difference of today and game.Datetime
					ratingHistoryMemory = ratingHistoryMemory[i:]
					break
				}
			}
			updatedRatingHistory, err := json.Marshal(ratingHistoryMemory)
			if err != nil {
				log.Println("removeHistoryFromPlayer problem marshalling ", err)
				return false
			}
			modifyRatingHistory(name, gametype, updatedRatingHistory)
		}

		return true
	} else {
		log.Println(err)
	}

	return false
}
