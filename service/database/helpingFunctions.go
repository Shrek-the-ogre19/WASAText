package database

import (
	"strconv"
	"strings"
)


func ConvertConversations(conversationsR string) ([]ConversationId){
	var conversations []ConversationId
	if conversationsR == ""{
		return conversations
	}
	conversations1 := strings.Split(conversationsR, ",")
	var size int
	size = len(conversations1)
	for i := 0; i < size; i++ {
		a := conversations1[i]
		number,err := strconv.Atoi(a)
		var id ConversationId
		if err == nil{
			id.Id = number
		}
			conversations = append(conversations, id)
	}
	return conversations
}


func ConvertUsers(usersR string) ([]UserId){
	var users []UserId
	if usersR == ""{
		return users
	}
	if usersR[:1] == ","{
		usersR = usersR[1:]
	}
	users1 := strings.Split(usersR, ",")
	var size int
	size = len(users1)
	for i := 0; i < size; i++ {
		a := users1[i]
		number,err := strconv.Atoi(a)
		var id UserId
		if err == nil{
			id.Id = number
		}
		users = append(users, id)
	}
	return users
}

func ConvertMessages(messagesR string) ([]MessageId){
	var messages []MessageId
	if messagesR == ""{
		return messages
	}
	messagesR = messagesR[1:]
	messages1 := strings.Split(messagesR, ",")
	var size int
	size = len(messages1)
	for i := 0; i < size; i++ {
		a := messages1[i]
		number,err := strconv.Atoi(a)
		var id MessageId
		if err == nil{
			id.Id = number
		}
		messages = append(messages, id)
	}
	return messages
}


func ConvertComments(commentsR string) ([]CommentId){
	var comments []CommentId
	if commentsR == ""{
		return comments
	}
	commentsR = commentsR[1:]
	comments1 := strings.Split(commentsR, ",")
	var size int
	size = len(comments1)
	for i := 0; i < size; i++ {
		a := comments1[i]
		number,err := strconv.Atoi(a)
		var id CommentId
		if err == nil{
			id.Id = number
		}
		comments = append(comments, id)
	}
	return comments
}
