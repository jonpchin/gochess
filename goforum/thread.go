package goforum

type Thread struct {
	ID       int    // Unique ID of the thread
	ForumID  int    // Used to find all threads in a forum section
	Username string // The one who created the thread
	Title    string // Title of thread
	Views    int    // Number of views the thread has
	Replies  int    // Number of replies the thread has
	LastPost string // The user who last made a post
	Date     string // Date when the thread was created
}

// Gets threads from forumId
func GetThreads(forumId string) (threads []Thread) {
	return threads
}
