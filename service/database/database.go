/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.

To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.

For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/
package database

import (
	"database/sql"
	"errors"
	"fmt"
)

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	GetName() (string, error)
	SetName(name string) error

	Ping() error

	CountUsers() (int, error)
	AddUser(name string) error
	ChangeUserName(id UserId, name string) error
	ChangeUserPicture(id UserId, picture string) error
	GetUserId(name string) (UserId, error)
	GetUser(id UserId) (User, error)
	ListAllUsers() ([]string, error)
	UserLookup(name string) (bool, error)
	AddConversations(id UserId, newConvId int) error

	GetConversations(id UserId) ([]ConversationId, error)
	CountConversations() (int, error)
	StartConversation(sender UserId, receivers []string) (ConversationId, error)
	GetConversation(id ConversationId) (Conversation, error)
	ChangeConversationName(id ConversationId, NewName string) error
	ChangeConversationPhoto(id ConversationId, NewPicture string) error
	GetAllMembers(id ConversationId) ([]string, error)
	AddMembers(id ConversationId, Name string) error
	LeaveGroup(convid ConversationId, userid UserId) error
	CheckGroupchat(id ConversationId) (bool, error)
	UpdateConversationHead(id ConversationId, snippetNew string, dateNew string) error
	UpdateContent(id ConversationId, newMessageId int) error
	RemoveFromContent(id ConversationId, oldMessageId int) error
	GetConversationMessages(id ConversationId) ([]MessageId, error)
	CheckIfInConversation(userId UserId, receiver string) (ConversationId, error)
	IsConversationMember(conversationId ConversationId, userId UserId) (bool, error)

	CountMessages() (int, error)
	CreateMessage(senderId UserId, conversationId ConversationId, text string) (MessageId, error)
	GetMessage(id MessageId) (Message, error)
	MarkMessageRead(messageId MessageId, readerId UserId) error
	ForwardMessage(userId UserId, messageId MessageId, conversationId ConversationId) (MessageId, error)
	DeleteMessage(messageId MessageId, conversationId ConversationId) error
	GetComments(id MessageId) ([]CommentId, error)
	AddCommentToMessage(messageId MessageId, newCommentId CommentId) error
	RemoveCommentFromMessage(messageId MessageId, oldCommentId CommentId) error

	CountComments() (int, error)
	CreateComment(userId UserId, messageId MessageId, content string) (CommentId, error)
	GetComment(id CommentId) (Comment, error)
	DeleteComment(id CommentId, messageId MessageId) error
	CheckIfUserCommented(id CommentId, userId UserId) (bool, error)
}

type appdbimpl struct {
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	// Check if table exists. If not, the database is empty, and we need to create the structure
	/*
		var tableName string
		err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='example_table';`).Scan(&tableName)
		if errors.Is(err, sql.ErrNoRows) {
			sqlStmt := `CREATE TABLE example_table (id INTEGER NOT NULL PRIMARY KEY, name TEXT);`
			_, err = db.Exec(sqlStmt)
			if err != nil {
				return nil, fmt.Errorf("error creating database structure: %w", err)
			}
		}
	*/

	CreateTableUsers := `CREATE TABLE IF NOT EXISTS users (
	id INTEGER NOT NULL UNIQUE,
	name TEXT NOT NULL UNIQUE,
	picture TEXT,
	conversations TEXT);`
	_, err := db.Exec(CreateTableUsers)
	if err != nil {
		return nil, fmt.Errorf("error creating database structure: %w", err)
	}

	CreateTableConversation := `CREATE TABLE IF NOT EXISTS conversations (
	id INTEGER NOT NULL UNIQUE,
	snippet TEXT,
	name TEXT,
	picture TEXT,
	date DATETIME,
	content TEXT,
	groupchat BOOL,
	members TEXT);`
	_, err = db.Exec(CreateTableConversation)
	if err != nil {
		return nil, fmt.Errorf("error creating database structure: %w", err)
	}

	CreateTableMessages := `CREATE TABLE IF NOT EXISTS messages (
	id INTEGER NOT NULL UNIQUE,
	status    TEXT,
	content   TEXT,
	comments   TEXT,
	timestamp DATETIME,
	senderId INTEGER,
	conversationId INTEGER);`
	_, err = db.Exec(CreateTableMessages)
	if err != nil {
		return nil, fmt.Errorf("error creating database structure: %w", err)
	}

	CreateTableComments := `CREATE TABLE IF NOT EXISTS comments (
	id INTEGER NOT NULL UNIQUE,
	content   TEXT,
	user int,
	message int);`
	_, err = db.Exec(CreateTableComments)
	if err != nil {
		return nil, fmt.Errorf("error creating database structure: %w", err)
	}

	_, _ = db.Exec(`ALTER TABLE messages ADD COLUMN readBy TEXT DEFAULT ''`)

	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
