package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Shrek-the-ogre19/WASAText/service/database"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")
	messageIdint, err := strconv.Atoi(ps.ByName("messageId"))
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	messageId := database.MessageId{messageIdint}
	conversationIdint, err := strconv.Atoi(ps.ByName("conversationId"))
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	var conversationId = database.ConversationId{conversationIdint}
	count, err := rt.db.CountMessages()
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	if messageIdint < 0 || messageIdint > count {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	var conversationMessages []database.MessageId
	conversationMessages, err = rt.db.GetConversationMessages(conversationId)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	for i := 0; i < len(conversationMessages); i++ {
		if conversationMessages[i] == messageId {
			var message database.Message
			message, err = rt.db.GetMessage(messageId)
			if err != nil {
				w.WriteHeader(http.StatusBadGateway)
				return
			}
			err = json.NewEncoder(w).Encode(message)
			if err != nil {
				w.WriteHeader(http.StatusBadGateway)
				return
			}
		}
	}

}

func (rt *_router) forwardMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")
	var requestData struct {
		Conversation string `json:"conversation"`
	}
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	var newconversation string = requestData.Conversation

	messageIdint, err := strconv.Atoi(ps.ByName("messageId"))
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	messageId := database.MessageId{messageIdint}
	userIdint, err := strconv.Atoi(ps.ByName("Id"))
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	var userId = database.UserId{userIdint}
	var conversations []database.ConversationId
	conversations, err = rt.db.GetConversations(userId)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	var newuserId database.UserId
	newuserId, err = rt.db.GetUserId(newconversation)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}

	for i := 0; i < len(conversations); i++ {
		var conversation database.Conversation
		conversation, err = rt.db.GetConversation(conversations[i])
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			return
		}
		if conversation.Name == newconversation || conversation.Members[0] == newuserId || conversation.Members[1] == newuserId {

			err = rt.db.ForwardMessage(userId, messageId, conversations[i])
			if err != nil {
				w.WriteHeader(http.StatusBadGateway)
				return
			}
			break
		}
	}

	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	err = json.NewEncoder(w).Encode(messageId)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}

}

func (rt *_router) deleteMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	messageIdint, err := strconv.Atoi(ps.ByName("messageId"))
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	messageId := database.MessageId{messageIdint}
	userIdint, err := strconv.Atoi(ps.ByName("Id"))
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	var userId = database.UserId{userIdint}
	var message database.Message
	conversationIdint, err := strconv.Atoi(ps.ByName("conversationId"))
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	var conversationId = database.ConversationId{conversationIdint}
	message, err = rt.db.GetMessage(messageId)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}

	if message.Sender == userId {
		err = rt.db.DeleteMessage(messageId, conversationId)
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			return
		}
	}

}

func (rt *_router) commentMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")
	var requestData struct {
		Comment string `json:"comment"`
	}
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var newComment string = requestData.Comment
	messageIdint, err := strconv.Atoi(ps.ByName("messageId"))
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	messageId := database.MessageId{messageIdint}
	userIdint, err := strconv.Atoi(ps.ByName("Id"))
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	var userId = database.UserId{userIdint}
	var commentId database.CommentId
	var comments []database.CommentId
	var message database.Message
	message, err = rt.db.GetMessage(messageId)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	comments = message.Comments
	for i := 0; i < len(comments); i++ {
		var comment database.Comment
		comment, err = rt.db.GetComment(comments[i])
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			return
		}
		if userId == comment.User {
			return
		}
	}

	commentId, err = rt.db.CreateComment(userId, messageId, newComment)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	err = json.NewEncoder(w).Encode(commentId)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}

}

func (rt *_router) getComments(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")
	messageIdint, err := strconv.Atoi(ps.ByName("messageId"))
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	messageId := database.MessageId{messageIdint}
	var comments []database.CommentId
	comments, err = rt.db.GetComments(messageId)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	err = json.NewEncoder(w).Encode(comments)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
}

func (rt *_router) getSpecificComment(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")
	commentIdint, err := strconv.Atoi(ps.ByName("commentId"))
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	var commentId = database.CommentId{commentIdint}
	messageIdint, err := strconv.Atoi(ps.ByName("messageId"))
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	messageId := database.MessageId{messageIdint}
	count, err := rt.db.CountComments()
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	if commentIdint < 0 || commentIdint > count {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var message database.Message
	message, err = rt.db.GetMessage(messageId)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	for i := 0; i < len(message.Comments); i++ {
		if message.Comments[i] == commentId {
			var comment database.Comment
			comment, err = rt.db.GetComment(commentId)
			if err != nil {
				w.WriteHeader(http.StatusBadGateway)
				return
			}
			err = json.NewEncoder(w).Encode(comment)
			if err != nil {
				w.WriteHeader(http.StatusBadGateway)
				return
			}
		}
	}
}

func (rt *_router) uncommentMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	messageIdint, err := strconv.Atoi(ps.ByName("messageId"))
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	messageId := database.MessageId{messageIdint}
	userIdint, err := strconv.Atoi(ps.ByName("Id"))
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	var userId = database.UserId{userIdint}
	commentIdint, err := strconv.Atoi(ps.ByName("commentId"))
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	var commentId = database.CommentId{commentIdint}
	var comment database.Comment
	comment, err = rt.db.GetComment(commentId)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}

	if comment.User == userId {
		err = rt.db.DeleteComment(commentId, messageId)
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			return
		}
	}
}
