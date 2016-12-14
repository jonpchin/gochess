package gostuff

import (
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
	Rating   string
}

// fetches rating history, unmarshals it, adds a new rating history, then marshals data and then
// store it back in the database returns true if sucessfully updates rating history with no errors
func updateRatingHistory(name string, gameType string, rating string) bool {

	problems, _ := os.OpenFile("logs/errors.txt", os.O_APPEND|os.O_WRONLY, 0666)
	defer problems.Close()
	log := log.New(problems, "", log.LstdFlags|log.Lshortfile)

	//check if database connection is open
	if db.Ping() != nil {
		log.Println("DATABASE DOWN!")
		return false
	}

	var ratingHistory string

	// getting player's rating history
	err := db.QueryRow("SELECT ? FROM ratinghistory WHERE username=?", name).Scan(&ratingHistory)

	if err != nil {
		log.Println(err)
		return false
	}

	// used to append current game rating history into rating history memory
	var ratingHistoryMemory []RatingDate

	//unmarshall JSON string into ratingHistoryMemory which is a memory model
	if err := json.Unmarshal([]byte(ratingHistory), &ratingHistoryMemory); err != nil {
		fmt.Println("Just receieved a message I couldn't decode:", ratingHistory, err)
		return false
	}

	var ratingInfo RatingDate
	ratingInfo.DateTime = time.Now().String()
	ratingInfo.Rating = rating
	ratingHistoryMemory = append(ratingHistoryMemory, ratingInfo)

	// need to marshall memory model before storing in database
	updatedRatingHistory, err := json.Marshal(ratingHistoryMemory)
	if err != nil {
		fmt.Println("updateRatingHistory problem marshalling ", err)
		return false
	}

	//store in database
	stmt, err := db.Prepare("INSERT ratinghistory SET ?=? WHERE username=?")
	if err != nil {
		log.Println(err)
		return false
	}

	_, err = stmt.Exec(gameType, updatedRatingHistory, name)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}
