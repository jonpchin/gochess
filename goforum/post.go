package goforum

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
