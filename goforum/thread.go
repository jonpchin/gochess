package goforum

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jonpchin/gochess/gostuff"
)

type ThreadSection struct {
	Title   string
	Threads []Thread
}

type Thread struct {
	ID         int64  // Unique ID of the thread
	ForumID    int    // Used to find all threads in a forum section
	ForumTitle string // Title of the forum
	Username   string // The one who created the thread
	Title      string // Title of thread
	Views      int    // Number of views the thread has
	Replies    int    // Number of replies the thread has
	LastPost   string // The user who last made a post
	Date       string // Date when the thread was created
	Posts      []Post // List of posts in the Thread
}

// Gets threads from forumId
func GetThreads(forumId string) (threadSection ThreadSection) {

	log := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

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

// Creates the first post in a thread, must be logged in to do this
func SendFirstForumPost(w http.ResponseWriter, r *http.Request) {
	username, err := r.Cookie("username")
	if err == nil {
		sessionID, err := r.Cookie("sessionID")
		if err == nil {
			if gostuff.SessionManager[username.Value] == sessionID.Value {

				var thread Thread
				date := time.Now().Format("20060102150405")
				thread.ForumTitle = template.HTMLEscapeString(r.FormValue("forumname"))
				thread.Username = username.Value
				threadTitle := template.HTMLEscapeString(r.FormValue("title"))
				thread.Title = threadTitle
				thread.Views = 0
				thread.Replies = 0
				thread.LastPost = username.Value
				thread.Date = date

				var post Post
				// First post of thread always starts with ID zero
				post.OrderID = 0
				post.Body = template.HTMLEscapeString(r.FormValue("message"))
				post.Username = username.Value
				post.Title = threadTitle
				post.Date = date

				thread.Posts = append(thread.Posts, post)

				updated, forumId := updateForumCount(thread.ForumTitle, post.Username)
				if updated {
					thread.ForumID = forumId
					thread.createThread()
					w.Write([]byte(""))
					return
				}
			}
		}
	}
	w.Write([]byte("<img src='img/ajax/not-available.png' /> Invalid credentials"))
}

// Returns true if succesfully updated forum count, also returns forumId of the forumName
func updateForumCount(forumName string, name string) (bool, int) {

	log := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

	var id int
	var totalthreads int
	var totalposts int

	err := db.QueryRow("SELECT id, totalthreads, totalposts from forums where title=?", forumName).Scan(
		&id, &totalthreads, &totalposts)
	if err != nil {
		log.Println(err)
		return false, 0
	}

	stmt, err := db.Prepare("UPDATE forums SET totalthreads=?, totalposts=?, recentuser=?, date=? WHERE id=?")
	if err != nil {
		log.Println(err)
		return false, 0
	}

	totalthreads += 1
	totalposts += 1

	_, err = stmt.Exec(totalthreads, totalposts, name, time.Now(), id)
	if err != nil {
		log.Println(err)
		return false, 0
	}
	return true, id
}

// Creates new thread with message and title
// Returns false if failed to create a new thread
func (thread *Thread) createThread() bool {

	log := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

	// Do not set auto incrmement id, it will be automatically set
	stmt, err := db.Prepare(`INSERT threads SET forumId=?, username=?, title=?, views=?
		, replies=?, lastpost=?, date=?`)
	if err != nil {
		log.Println(err)
		return false
	}

	res, err := stmt.Exec(thread.ForumID, thread.Username, thread.Title, thread.Views, thread.Replies,
		thread.LastPost, thread.Date)
	if err != nil {
		log.Println(err)
		return false
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Println(err)
	} else {
		thread.ID = id
		thread.Posts[0].ThreadID = id
	}

	// A newly created thread only has 1 post
	return thread.Posts[0].createPost()
}
