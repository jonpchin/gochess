package goforum

type Post struct {
	ID       int
	ThreadID int
	OrderID  int // The order of this post in the thread
	Username string
	Title    string // Title of the thread the post is in
	Body     string
	Date     string
}

// Gets posts from thread ID
func GetPosts(threadId string) (posts []Post) {
	return posts
}
