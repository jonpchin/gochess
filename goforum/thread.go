package goforum

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
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
	Locked     string // No posts can be made on a locked thread, options: Yes or No
	Date       string // Date when the thread was created
	Posts      []Post // List of posts in the Thread
}

// Gets threads from forumId
func GetThreads(forumId string) (threadSection ThreadSection) {

	log := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

	rows, err := db.Query("SELECT * FROM threads WHERE forumId=?", forumId)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()

	var thread Thread

	for rows.Next() {

		err = rows.Scan(&thread.ID, &thread.ForumID, &thread.Username, &thread.Title,
			&thread.Views, &thread.Replies, &thread.LastPost, &thread.Locked, &thread.Date)

		if err != nil {
			log.Println(err)
		}
		threadSection.Threads = append(threadSection.Threads, thread)
	}
	threadSection.Title = GetForumTitle(forumId)
	return threadSection
}

// Creates the first post in a thread, must be logged in to do this
func SendForumPost(w http.ResponseWriter, r *http.Request) {
	username, err := r.Cookie("username")
	if err == nil {
		sessionID, err := r.Cookie("sessionID")
		if err == nil {
			if gostuff.SessionManager[username.Value] == sessionID.Value {

				log := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
				canPost, seconds := canUserPost(username.Value)

				if canPost == false {
					w.Write([]byte("Please wait " + seconds + " seconds before posting again."))
					return
				}

				firstPost := template.HTMLEscapeString(r.FormValue("firstPost"))
				tempThreadId := template.HTMLEscapeString(r.FormValue("threadId"))
				var threadId int64
				threadId = 0
				if tempThreadId != "" {
					threadId, err = strconv.ParseInt(tempThreadId, 10, 64)

					if err != nil {
						log.Println(err)
					}
				}

				forumTitle := template.HTMLEscapeString(r.FormValue("forumname"))
				threadTitle := template.HTMLEscapeString(r.FormValue("title"))
				message := template.HTMLEscapeString(r.FormValue("message"))
				totalPosts := template.HTMLEscapeString(r.FormValue("totalPosts"))
				date := time.Now().Format("20060102150405")

				if firstPost == "Yes" {

					var thread Thread

					thread.ForumTitle = forumTitle
					thread.Username = username.Value
					thread.Title = threadTitle
					thread.Views = 0
					thread.Replies = 0
					thread.LastPost = username.Value
					thread.Date = date

					var post Post
					// First post of thread always starts with ID zero
					post.OrderID = 0
					post.Body = message
					post.Username = username.Value
					post.Title = threadTitle
					post.Date = date

					thread.Posts = append(thread.Posts, post)
					updated, forumId := updateForumCount(thread.ForumTitle, post.Username)

					if updated {
						thread.ForumID = forumId
						if thread.createThread() {
							w.Write([]byte("createPost"))
						} else {
							w.Write([]byte("<img src='img/ajax/not-available.png' /> Failed to create new thread"))
						}
						return
					} else {
						w.Write([]byte("<img src='img/ajax/not-available.png' /> Failed to update forum count"))
						return
					}
				} else {

					var post Post

					newTotalPosts, err := strconv.Atoi(totalPosts)
					if err != nil {
						log.Println(err)
					}

					post.ThreadID = threadId
					post.OrderID = newTotalPosts
					post.Body = message
					post.Username = username.Value
					post.Title = threadTitle
					post.Date = date

					if post.createPost() {
						updateThreadReplies(post.ThreadID)
						updateForumPostCount(getForumId(forumTitle))
						w.Write([]byte("createPost"))
						return
					} else {
						w.Write([]byte("<img src='img/ajax/not-available.png' /> Failed to create post"))
						return
					}
				}
			}
		}
	}
	w.Write([]byte("<img src='img/ajax/not-available.png' /> Invalid credentials"))
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
		return false
	} else {
		thread.ID = id
		thread.Posts[0].ThreadID = id
	}

	// A newly created thread only has 1 post
	return thread.Posts[0].createPost()
}

// Returns true if thread is locked
func IsLocked(threadId string) bool {

	var locked string
	log := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

	err := db.QueryRow("SELECT locked from threads where id=?", threadId).Scan(&locked)
	if err != nil {
		log.Println(err)
		return false
	}
	if locked == "Yes" {
		return true
	}

	return false
}

func LockThread(w http.ResponseWriter, r *http.Request) {

	if gostuff.ValidateCredentials(w, r) == false {
		return
	}
	id := template.HTMLEscapeString(r.FormValue("id"))
	updateThreadLock("Yes", id)
}

func UnlockThread(w http.ResponseWriter, r *http.Request) {

	if gostuff.ValidateCredentials(w, r) == false {
		return
	}
	id := template.HTMLEscapeString(r.FormValue("id"))
	updateThreadLock("No", id)
}

// Updates the lock thread based on the lock string
// id is the id of the thread to update
func updateThreadLock(lock string, id string) {

	log := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

	stmt, err := db.Prepare("UPDATE threads SET locked=? WHERE id=?")
	if err != nil {
		log.Println(err)
		return
	}

	_, err = stmt.Exec(lock, id)
	if err != nil {
		log.Println(err)
	}
}

func updateThreadViewCount(threadId int64) {
	log := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

	stmt, err := db.Prepare("UPDATE threads SET views=views+1 WHERE id=?")
	if err != nil {
		log.Println(err)
		return
	}

	_, err = stmt.Exec(threadId)
	if err != nil {
		log.Println(err)
	}
}
