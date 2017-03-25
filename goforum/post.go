package goforum

import (
	"net/http"

	"github.com/jonpchin/gochess/gostuff"
)

type Post struct {
	ID       int    // Unique ID of the post
	ThreadID int    // Used to find other posts in the thread
	OrderID  int    // The order of this post in the thread
	Username string // User who made the post
	Title    string // Title of the thread the post is in
	Body     string // The actual message of the post
	Date     string // Date when the post was made
}

// Gets posts from thread ID
func GetPosts(threadId string) (posts []Post) {
	return posts
}

// Creates the first post in a thread, must be logged in to do this
func SendFirstForumPost(w http.ResponseWriter, r *http.Request) {
	username, err := r.Cookie("username")
	if err == nil {
		sessionID, err := r.Cookie("sessionID")
		if err == nil {
			if gostuff.SessionManager[username.Value] == sessionID.Value {

			}
		}
	}
	w.Write([]byte("<img src='img/ajax/not-available.png' /> Invalid credentials"))
}
