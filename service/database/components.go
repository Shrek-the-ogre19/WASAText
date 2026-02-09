package database

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
	Id   UserId
	Name string
	Picture string
	Conversations []ConversationId
}

type Conversation struct {
	Id      ConversationId
	Snippet string
	Name    string
	Picture string
	Date    string
	Content   []MessageId
	Groupchat bool
	Members   []UserId
}

type Message struct {
	Id        MessageId
	Status    string
	Content   string
	Comments   []CommentId
	Timestamp string
	Sender      UserId
	ConversationId ConversationId
}

type Comment struct {
	Id      CommentId
	Content string
	User    UserId
	Message MessageId
}
