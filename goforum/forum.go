package goforum

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jonpchin/gochess/gostuff"
)

type Forum struct {
	ID           int
	Title        string
	Description  string
	TotalThreads int
	TotalPosts   int
	RecentUser   string   // Most recent user that made a post
	RecentDate   string   // Most recent date the post was made
	Threads      []Thread // List of threads in forum
}

var db *sql.DB

func ConnectForumDb() {
	db = gostuff.GetDb()
}

type NullTime struct {
	Time  time.Time
	Valid bool // Valid is true if Time is not NULL
}

func GetForums() (forums []Forum) {

	problems, err := os.OpenFile("logs/errors.txt", os.O_APPEND|os.O_WRONLY, 0666)
	defer problems.Close()
	log := log.New(problems, "", log.LstdFlags|log.Lshortfile)

	rows, err := db.Query("SELECT * FROM forums")
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()

	var forum Forum

	for rows.Next() {

		err = rows.Scan(&forum.ID, &forum.Title, &forum.Description, &forum.TotalThreads,
			&forum.TotalPosts, &forum.RecentUser, &forum.RecentDate)

		if err != nil {
			log.Println(err)
		}
		forums = append(forums, forum)
	}
	return forums
}

// Returns forumId from forumId, if none is found it returns "0" as forumId
func GetForumIdFromName(forumName string) string {
	forumId := "0"
	err := db.QueryRow("SELECT id from forums where title=?", forumName).Scan(&forumId)
	if err != nil {
		fmt.Println("Could not fetch forumId from forumName", forumName, err)
	}
	return forumId
}

// Checks if 30 seconds has passed since a user has last post, returns true
// if the user is allowed to post
func CanUserPost(username string) bool {

	log := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

	var lastpost sql.NullString

	err := db.QueryRow("SELECT lastpost from userinfo WHERE username=?", username).Scan(&lastpost)
	if err != nil {
		log.Println(err)
		return false
	}

	timeFormat := "2006-01-02 15:04:05"

	if lastpost.Valid {

		then, err := time.Parse(timeFormat, lastpost.String)
		if err != nil {
			fmt.Println(err)
			return false
		}

		duration := time.Now().Sub(then)

		var timeZoneDiff float64
		timeZoneDiff = 14400.0

		timeDiff := duration.Seconds() - timeZoneDiff
		log.Println(timeDiff)
		if timeDiff < 30 {
			log.Println("Please wait 30 seconds before posting another post user:", username)
			return false
		}
	}

	updateLastPostTime(time.Now().Format(timeFormat), username)
	return true
}

func updateLastPostTime(dateTime string, username string) {

	log := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

	stmt, err := db.Prepare("UPDATE userinfo SET lastpost=? WHERE username=?")
	if err != nil {
		log.Println(err)
		return
	}

	_, err = stmt.Exec(dateTime, username)
	if err != nil {
		log.Println(err)
		return
	}
}
