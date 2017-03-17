package goforum

type Thread struct {
	ID       int
	ForumID  int
	Username string // The one who created the thread
	Title    string // Title of thread
	Views    int
	Replies  int
	LastPost string // The user who last made a post
	Date     string
}

// Gets threads from forumId
func GetThreads(forumId string) (threads []Thread) {
	return threads
}
