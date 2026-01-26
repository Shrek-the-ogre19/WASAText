package api

type UserId struct {
	Id int
}

type ConversationId struct {
	Id int
}

type MessageId struct {
	Id int
}

type CommentId struct {
	Id int
}

type User struct {
	Id      UserId
	Name    string
	Picture string
}

type Conversation struct {
	Id      ConversationId
	Snippet string
	Name    string
	Picture string
	Date    string
	Content []Message
}

type Message struct {
	Id        MessageId
	Status    string
	Content   string
	Comment   []Comment
	Timestamp string
	User      User
}

type Comment struct {
	Id      CommentId
	Content string
	User    User
}
