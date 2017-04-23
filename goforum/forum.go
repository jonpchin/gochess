package goforum

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
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

func GetForums() (forums []Forum) {

	log := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

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

func GetForumTitle(forumId string) string {

	var forumTitle string
	err := db.QueryRow("SELECT title from forums where id=?", forumId).Scan(&forumTitle)
	if err != nil {
		fmt.Println("Could not fetch forum title", err)
	}
	return forumTitle
}

func getForumId(forumName string) int {

	var id int

	err := db.QueryRow("SELECT id from forums where title=?", forumName).Scan(&id)
	if err != nil {
		log.Println(err)
		return 0 // Invalid id
	}
	return id
}

// Returns true if succesfully updated forum count, also returns forumId of the forumName
func updateForumCount(forumName string, name string) (bool, int) {

	log := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
	id := getForumId(forumName)

	stmt, err := db.Prepare(`UPDATE forums SET totalthreads=totalthreads+1, totalposts=totalposts+1
		, recentuser=?, date=? WHERE id=?`)
	if err != nil {
		log.Println(err)
		return false, 0
	}

	_, err = stmt.Exec(name, time.Now(), id)
	if err != nil {
		log.Println(err)
		return false, 0
	}
	return true, id
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
// if the user is allowed to post, also returns number of seconds
// user has to wait before posting again
func canUserPost(username string) (bool, string) {

	log := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

	var lastpost sql.NullString

	err := db.QueryRow("SELECT lastpost from userinfo WHERE username=?", username).Scan(&lastpost)
	if err != nil {
		log.Println(err)
		return false, "30"
	}

	timeFormat := "2006-01-02 15:04:05"

	// If not valid that means there is no existing timestamp in the database
	if lastpost.Valid {
		timeInSeconds := 30
		isElapse, timeDiff := gostuff.HasTimeElapsed(lastpost.String, timeInSeconds, timeFormat, false)
		if isElapse == false {
			diff := strconv.Itoa(int(timeInSeconds - timeDiff))
			log.Println("Please wait "+diff+" seconds before posting another post user:", username)
			return false, diff
		}
	}
	updateLastPostTime(time.Now().Format(timeFormat), username)
	return true, "0"
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

func updateForumPostCount(forumId int) {

	log := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

	stmt, err := db.Prepare("UPDATE forums SET totalposts=totalposts+1 WHERE id=?")
	if err != nil {
		log.Println(err)
		return
	}

	_, err = stmt.Exec(forumId)
	if err != nil {
		log.Println(err)
	}
}
