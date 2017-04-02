package goforum

import (
	"database/sql"
	"fmt"
	"log"
	"os"

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
