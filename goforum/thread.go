package goforum

import (
	"fmt"
	"log"
	"os"
)

type ThreadSection struct {
	Title   string
	Threads []Thread
}

type Thread struct {
	ID         int    // Unique ID of the thread
	ForumID    int    // Used to find all threads in a forum section
	ForumTitle string // Title of the forum
	Username   string // The one who created the thread
	Title      string // Title of thread
	Views      int    // Number of views the thread has
	Replies    int    // Number of replies the thread has
	LastPost   string // The user who last made a post
	Date       string // Date when the thread was created
}

// Gets threads from forumId
func GetThreads(forumId string) (threadSection ThreadSection) {

	problems, err := os.OpenFile("logs/errors.txt", os.O_APPEND|os.O_WRONLY, 0666)
	defer problems.Close()
	log := log.New(problems, "", log.LstdFlags|log.Lshortfile)

	rows, err := db.Query("SELECT * FROM threads")
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()

	var thread Thread

	for rows.Next() {

		err = rows.Scan(&thread.ID, &thread.ForumID, &thread.Username, &thread.Title,
			&thread.Views, &thread.Replies, &thread.LastPost, &thread.Date)

		if err != nil {
			log.Println(err)
		}
		threadSection.Threads = append(threadSection.Threads, thread)
	}
	threadSection.Title = getForumTitle(forumId)
	return threadSection
}

func getForumTitle(forumId string) string {

	var forumTitle string
	err := db.QueryRow("SELECT title from forums where id=?", forumId).Scan(&forumTitle)
	if err != nil {
		fmt.Println("Could not fetch forum title", err)
	}
	return forumTitle
}
