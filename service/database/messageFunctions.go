package database

import (
	"strconv"
	"time"
)

func (db *appdbimpl) CountMessages() (int, error) {
	var count int
	err := db.c.QueryRow("SELECT COUNT(*) FROM messages").Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, err
}

func (db *appdbimpl) CreateMessage(senderId UserId, conversationId ConversationId, text string) (MessageId, error) {
	id, err := db.CountMessages()
	if err != nil {
		return MessageId{}, err
	}
	status := "sent"
	timestamp := time.Now().String()

	_, err = db.c.Exec("INSERT INTO messages (id, status, content, comments, timestamp, senderId, conversationId) VALUES (?, ?, ?, ?, ?, ?, ?)", id+1, status, text, "", timestamp, senderId.Id, conversationId.Id)
	if err != nil {
		return MessageId{}, err
	}
	err = db.UpdateContent(conversationId, id+1)
	if err == nil {
		user, userErr := db.GetUser(senderId)
		if userErr == nil {
			senderName := user.Name
			//For snippet we cut username to 6 characters and text to 30
			if len(senderName) > 6 {
				senderName = senderName[:6] + "..."
			}
			snippet := senderName + ": " + text
			db.UpdateConversationHead(conversationId, snippet, timestamp)
		}
	}

	return MessageId{id + 1}, err
}

func (db *appdbimpl) GetMessage(id MessageId) (Message, error) {
	var status string
	var content string
	var comments []CommentId
	var timestamp string
	var senderId UserId
	var conversationId ConversationId
	var senderIdR int
	var conversationIdR int
	var commentsR string
	var message Message

	err := db.c.QueryRow("SELECT status, content, comments, timestamp, senderId, conversationId FROM messages WHERE id = ?", id.Id).Scan(&status, &content, &commentsR, &timestamp, &senderIdR, &conversationIdR)

	comments = ConvertComments(commentsR)
	senderId = UserId{senderIdR}
	conversationId = ConversationId{conversationIdR}

	if err == nil {
		message = Message{id, status, content, comments, timestamp, senderId, conversationId}
	}

	return message, err
}

func (db *appdbimpl) UpdateStatus(id MessageId) error {
	var newstatus = "read"
	_, err := db.c.Exec(`UPDATE messages
SET status = ?
WHERE id = ?;`, newstatus, id.Id)
	return err
}

func (db *appdbimpl) ForwardMessage(userId UserId, messageId MessageId, conversationId ConversationId) error {
	id, err := db.CountMessages()
	if err != nil {
		return err
	}
	status := "sent"
	timestamp := time.Now().String()

	message, err := db.GetMessage(messageId)
	if err != nil {
		return err
	}
	text := "FORWARDED: " + message.Content

	_, err = db.c.Exec("INSERT INTO messages (id, status, content, comments, timestamp, senderId, conversationId) VALUES (?, ?, ?, ?, ?, ?, ?)", id+1, status, text, "", timestamp, userId.Id, conversationId.Id)
	if err != nil {
		return err
	}
	err = db.UpdateContent(conversationId, id+1)
	if err == nil {
		user, userErr := db.GetUser(userId)
		if userErr == nil {
			senderName := user.Name
			//For snippet we cut username to 6 characters and text to 30
			if len(senderName) > 6 {
				senderName = senderName[:6] + "..."
			}
			snippetText := text
			if len(snippetText) > 30 {
				snippetText = snippetText[:30] + "..."
			}
			snippet := senderName + " :  " + snippetText
			db.UpdateConversationHead(conversationId, snippet, timestamp)
		}
	}
	return err
}

func (db *appdbimpl) DeleteMessage(messageId MessageId, conversationId ConversationId) error {
	_, err := db.c.Exec(`UPDATE messages
SET status = ?, content = ?, comments = ? , timestamp = ?
WHERE id = ?;`, "", "", "", "", messageId.Id)

	if err == nil {
		db.RemoveFromContent(conversationId, messageId.Id)

	}
	return err
}

func (db *appdbimpl) GetComments(id MessageId) ([]CommentId, error) {
	var comments []CommentId
	var commentsR string

	err := db.c.QueryRow("SELECT  comments FROM messages WHERE id = ?", id.Id).Scan(&commentsR)
	comments = ConvertComments(commentsR)
	return comments, err
}

func (db *appdbimpl) AddCommentToMessage(messageId MessageId, newCommentId CommentId) error {
	var commentsR string
	err := db.c.QueryRow("SELECT  comments FROM messages WHERE id = ?", messageId.Id).Scan(&commentsR)
	if err != nil {
		return err
	}
	commentsR = commentsR + "," + strconv.Itoa(newCommentId.Id)

	_, err = db.c.Exec(`UPDATE messages
SET comments = ?
WHERE id = ?;`, commentsR, messageId.Id)
	return err
}

func (db *appdbimpl) RemoveCommentFromMessage(messageId MessageId, oldCommentId CommentId) error {
	var commentsR string
	err := db.c.QueryRow("SELECT  comments FROM messages WHERE id = ?", messageId.Id).Scan(&commentsR)
	if err != nil {
		return err
	}
	comments := ConvertComments(commentsR)
	var newcomments string
	for i := 0; i < len(comments); i++ {
		if comments[i].Id != oldCommentId.Id {
			newcomments = newcomments + "," + strconv.Itoa(comments[i].Id)
		}
	}

	_, err = db.c.Exec(`UPDATE messages
SET comments = ?
WHERE id = ?;`, newcomments, messageId.Id)
	return err
}
