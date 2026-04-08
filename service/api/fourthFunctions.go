package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Shrek-the-ogre19/WASAText/service/database"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")
	messageIdint, err := strconv.Atoi(ps.ByName("messageId"))
	messageId := database.MessageId{messageIdint}
	conversationIdint, err := strconv.Atoi(ps.ByName("conversationId"))
	var conversationId = database.ConversationId{conversationIdint}
	count, err := rt.db.CountMessages()
	if err != nil {
		fmt.Println("error in function CountConversations")
	}
	if messageIdint < 0 || messageIdint > count {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	var conversationMessages []database.MessageId
	conversationMessages, err = rt.db.GetConversationMessages(conversationId)
	for i := 0; i < len(conversationMessages); i++ {
		if conversationMessages[i] == messageId {
			var message database.Message
			message, err = rt.db.GetMessage(messageId)
			json.NewEncoder(w).Encode(message)
		}
	}

}

func (rt *_router) forwardMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")
	var requestData struct {
		Conversation string `json:"conversation"`
	}
	err := json.NewDecoder(r.Body).Decode(&requestData)
	var newconversation string = requestData.Conversation

	messageIdint, err := strconv.Atoi(ps.ByName("messageId"))
	messageId := database.MessageId{messageIdint}
	userIdint, err := strconv.Atoi(ps.ByName("Id"))
	var userId = database.UserId{userIdint}
	var conversations []database.ConversationId
	conversations, err = rt.db.GetConversations(userId)

	for i := 0; i < len(conversations); i++ {
		var conversation database.Conversation
		conversation, err = rt.db.GetConversation(conversations[i])
		if conversation.Name == newconversation {
			err = rt.db.ForwardMessage(userId, messageId, conversations[i])
			break
		}
	}

	if err != nil {
		fmt.Println("error in function")
	}
	json.NewEncoder(w).Encode(messageId)

}

func (rt *_router) deleteMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	messageIdint, err := strconv.Atoi(ps.ByName("messageId"))
	messageId := database.MessageId{messageIdint}
	userIdint, err := strconv.Atoi(ps.ByName("Id"))
	var userId = database.UserId{userIdint}
	var message database.Message
	conversationIdint, err := strconv.Atoi(ps.ByName("conversationId"))
	var conversationId = database.ConversationId{conversationIdint}
	message, err = rt.db.GetMessage(messageId)

	if message.Sender == userId {
		err = rt.db.DeleteMessage(messageId, conversationId)
	}

	if err != nil {
		fmt.Println("error in function")
	}
}

func (rt *_router) commentMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")
	var requestData struct {
		Comment string `json:"comment"`
	}
	err := json.NewDecoder(r.Body).Decode(&requestData)
	var newComment string = requestData.Comment
	messageIdint, err := strconv.Atoi(ps.ByName("messageId"))
	messageId := database.MessageId{messageIdint}
	userIdint, err := strconv.Atoi(ps.ByName("Id"))
	var userId = database.UserId{userIdint}
	var commentId database.CommentId
	commentId, err = rt.db.CreateComment(userId, messageId, newComment)
	json.NewEncoder(w).Encode(commentId)
	if err != nil {
		fmt.Println("error in function")
	}

}

func (rt *_router) getComments(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")
	messageIdint, err := strconv.Atoi(ps.ByName("messageId"))
	messageId := database.MessageId{messageIdint}
	var comments []database.CommentId
	comments, err = rt.db.GetComments(messageId)
	json.NewEncoder(w).Encode(comments)

	if err != nil {
		fmt.Println("error in function")
	}
}

func (rt *_router) getSpecificComment(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")
	commentIdint, err := strconv.Atoi(ps.ByName("commentId"))
	var commentId = database.CommentId{commentIdint}
	messageIdint, err := strconv.Atoi(ps.ByName("messageId"))
	messageId := database.MessageId{messageIdint}
	count, err := rt.db.CountComments()
	if commentIdint < 0 || commentIdint > count {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var message database.Message
	message, err = rt.db.GetMessage(messageId)
	for i := 0; i < len(message.Comments); i++ {
		if message.Comments[i] == commentId {
			var comment database.Comment
			comment, err = rt.db.GetComment(commentId)
			json.NewEncoder(w).Encode(comment)
		}
	}

	if err != nil {
		fmt.Println("error in function")
	}
}

func (rt *_router) uncommentMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	messageIdint, err := strconv.Atoi(ps.ByName("messageId"))
	messageId := database.MessageId{messageIdint}
	userIdint, err := strconv.Atoi(ps.ByName("Id"))
	var userId = database.UserId{userIdint}
	commentIdint, err := strconv.Atoi(ps.ByName("commentId"))
	var commentId = database.CommentId{commentIdint}
	var comment database.Comment
	comment, err = rt.db.GetComment(commentId)

	if comment.User == userId {
		err = rt.db.DeleteComment(commentId, messageId)
	}

	if err != nil {
		fmt.Println("error in function")
	}
}
