package database

import (
	"strconv"
	"time"
)

func (db *appdbimpl) CountMessages()(int, error){
	var count int
	err := db.c.QueryRow("SELECT COUNT(*) FROM messages").Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, err
}


func (db *appdbimpl) CreateMessage(senderId UserId, conversationId ConversationId, text string)(error){
	id, err := db.CountMessages()
	var status = "sent"
	var content = text
	var comments = ""
	var timestamp = time.Now().String()
	var snippet string
	var senderName string
	var user User


	_, err = db.c.Exec("INSERT INTO messages (id, status, content, comments, timestamp, senderId, conversationId) VALUES (?, ?, ?, ?, ?, ?, ?)", id+1, status, content, comments, timestamp, senderId.Id, conversationId.Id )
	err = db.UpdateContent(conversationId, id+1)
	if err == nil{
		user, err = db.GetUser(senderId)
		senderName = user.Name
		//For snippet we cut username to 6 characters and text to 30
		if len(senderName)> 6{
			senderName = senderName[:6] + "..."
		}
		if len(text)> 30{
			text = text[:30] + "..."
		}
		snippet = senderName + ": " + text

		db.UpdateConversationHead(conversationId, snippet, timestamp)
	}


	return err
}


func (db *appdbimpl) GetMessage(id MessageId)(Message, error){
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

	if err ==nil{
		message = Message{id, status, content, comments, timestamp, senderId, conversationId}
	}

	return message, err
}

func (db *appdbimpl) UpdateStatus(id MessageId)(error){
	var newstatus = "read"
	_, err := db.c.Exec(`UPDATE messages
SET status = ?
WHERE id = ?;`, newstatus, id.Id)
	return err
}

func (db *appdbimpl) ForwardMessage(userId UserId, messageId MessageId, conversationId ConversationId)(error){
	id, err := db.CountMessages()
	var status = "sent"
	var comments = ""
	var timestamp = time.Now().String()

	var message Message
	message, err = db.GetMessage(messageId)
	var snippet string
	var senderName string
	var user User
	text := message.Content
	senderId := userId


	_, err = db.c.Exec("INSERT INTO messages (id, status, content, comments, timestamp, senderId, conversationId) VALUES (?, ?, ?, ?, ?, ?, ?)", id+1, status, text, comments, timestamp, senderId.Id, conversationId.Id )
	err = db.UpdateContent(conversationId, id+1)
	if err == nil{
		user, err = db.GetUser(senderId)
		senderName = user.Name
		//For snippet we cut username to 6 characters and text to 30
		if len(senderName)> 6{
			senderName = senderName[:6] + "..."
		}
		if len(text)> 30{
			text = text[:30] + "..."
		}
		snippet = senderName + " :  " + text



		db.UpdateConversationHead(conversationId, snippet, timestamp)
	}
	return err
}

func (db *appdbimpl) DeleteMessage(messageId MessageId, conversationId ConversationId)(error) {
	_, err := db.c.Exec(`UPDATE messages
SET status = ?, content = ?, comments = ? , timestamp = ?
WHERE id = ?;`, "", "", "", "", messageId.Id)

	if err == nil {
		db.RemoveFromContent(conversationId, messageId.Id)

	}
	return err
}

func (db *appdbimpl) GetComments(id MessageId)([]CommentId, error){
	var comments []CommentId
	var commentsR string

	err := db.c.QueryRow("SELECT  comments FROM messages WHERE id = ?", id.Id).Scan(&commentsR)
	comments = ConvertComments(commentsR)
	return comments, err
}

func (db *appdbimpl) AddCommentToMessage(messageId MessageId, newCommentId CommentId)(error){
	var commentsR string
	err := db.c.QueryRow("SELECT  comments FROM messages WHERE id = ?", messageId.Id).Scan(&commentsR)
	commentsR = commentsR + "," + strconv.Itoa(newCommentId.Id)


	_, err = db.c.Exec(`UPDATE messages
SET comments = ?
WHERE id = ?;`, commentsR, messageId.Id)
	return err
}

func (db *appdbimpl) RemoveCommentFromMessage(messageId MessageId, oldCommentId CommentId)(error){
	var commentsR string
	err := db.c.QueryRow("SELECT  comments FROM messages WHERE id = ?", messageId.Id).Scan(&commentsR)
	comments := ConvertComments(commentsR)
	var newcomments string
	for i := 0; i < len(comments); i++{
		if comments[i].Id != oldCommentId.Id{
			newcomments= newcomments + "," + strconv.Itoa(comments[i].Id)
		}
	}

	_, err = db.c.Exec(`UPDATE messages
SET comments = ?
WHERE id = ?;`, newcomments, messageId.Id)
	return err
}
