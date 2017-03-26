package goforum

import (
	"log"
	"os"
)

type Post struct {
	ID       int    // Unique ID of the post
	ThreadID int64  // Used to find other posts in the thread
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

// Returns true if post was created with no errors
func (post *Post) createPost() bool {

	log := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

	stmt, err := db.Prepare(`INSERT posts SET threadId=?, orderId=?, username=?, title=?
		, body=?, date=?`)
	if err != nil {
		log.Println(err)
		return false
	}

	_, err = stmt.Exec(post.ThreadID, post.OrderID, post.Username, post.Title,
		post.Body, post.Date)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}
