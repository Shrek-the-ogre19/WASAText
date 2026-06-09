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
	conversationId := database.ConversationId{conversationIdint}
	count, err := rt.db.CountMessages()
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	if messageIdint < 0 || messageIdint > count {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	conversationMessages, err := rt.db.GetConversationMessages(conversationId)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	for i := 0; i < len(conversationMessages); i++ {
		if conversationMessages[i] == messageId {
			message, err := rt.db.GetMessage(messageId)
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
	userId := database.UserId{userIdint}
	conversations, err := rt.db.GetConversations(userId)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	newuserId, err := rt.db.GetUserId(requestData.Conversation)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}

	for i := 0; i < len(conversations); i++ {
		conversation, err := rt.db.GetConversation(conversations[i])
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			return
		}
		if conversation.Name == requestData.Conversation || conversation.Members[0] == newuserId || conversation.Members[1] == newuserId {

			err = rt.db.ForwardMessage(userId, messageId, conversations[i])
			if err != nil {
				w.WriteHeader(http.StatusBadGateway)
				return
			}
			break
		}
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
	userId := database.UserId{userIdint}
	conversationIdint, err := strconv.Atoi(ps.ByName("conversationId"))
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	conversationId := database.ConversationId{conversationIdint}
	message, err := rt.db.GetMessage(messageId)
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
	userId := database.UserId{userIdint}
	message, err := rt.db.GetMessage(messageId)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	for i := 0; i < len(message.Comments); i++ {
		comment, err := rt.db.GetComment(message.Comments[i])
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			return
		}
		if userId == comment.User {
			return
		}
	}

	commentId, err := rt.db.CreateComment(userId, messageId, requestData.Comment)
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
	comments, err := rt.db.GetComments(messageId)
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
	commentId := database.CommentId{commentIdint}
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

	message, err := rt.db.GetMessage(messageId)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	for i := 0; i < len(message.Comments); i++ {
		if message.Comments[i] == commentId {
			comment, err := rt.db.GetComment(commentId)
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
	userId := database.UserId{userIdint}
	commentIdint, err := strconv.Atoi(ps.ByName("commentId"))
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	commentId := database.CommentId{commentIdint}
	comment, err := rt.db.GetComment(commentId)
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
