package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getMyConversations(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")
	userId, err := strconv.Atoi(ps.ByName("Id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	count, err := rt.db.CountUsers()
	if err != nil {
		fmt.Println("error in function CountUsers")
	}
	if userId < 0 || userId > count {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	var conversationIds []database.ConversationId
	conversationIds, err = rt.db.GetConversations(database.UserId{userId})
	var conversations []database.Conversation
	for i := 0; i < len(conversationIds); i++ {
		var toAdd database.Conversation
		toAdd, err = rt.db.GetConversation(conversationIds[i])
		if err != nil {
			fmt.Println("error in function GetConversation")
		}
		if toAdd.Groupchat == false {
			var members []string
			members, err = rt.db.GetAllMembers(conversationIds[i])
			var memberid database.UserId
			memberid, err = rt.db.GetUserId(members[0])

			if memberid.Id == userId {
				var receiverId database.UserId
				receiverId, err = rt.db.GetUserId(members[1])
				var receiver database.User
				receiver, err = rt.db.GetUser(receiverId)
				toAdd.Picture = receiver.Picture
				toAdd.Name = receiver.Name
			} else {

				var receiverId database.UserId
				receiverId, err = rt.db.GetUserId(members[0])
				var receiver database.User
				receiver, err = rt.db.GetUser(receiverId)
				toAdd.Picture = receiver.Picture
				toAdd.Name = receiver.Name
			}
		}
		conversations = append(conversations, toAdd)

	}

	if err != nil {
		fmt.Println("error in GetConversations function")
	}
	json.NewEncoder(w).Encode(conversations)
}

func (rt *_router) startNewConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")
	userId, err := strconv.Atoi(ps.ByName("Id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	count, err := rt.db.CountUsers()
	if err != nil {
		fmt.Println("error in function CountUsers")
	}
	if userId < 0 || userId > count {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var requestData struct {
		Receivers string `json:"receivers"`
	}
	err = json.NewDecoder(r.Body).Decode(&requestData)
	var receiversst string = requestData.Receivers
	receivers := strings.Split(receiversst, ",")

	var conversationId database.ConversationId
	for _, receiver := range receivers {
		exists, err := rt.db.UserLookup(receiver)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if exists == false {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Println("user doesn't exist")
			return
		}
	}

	if len(receivers) == 1 {
		conversationId, err = rt.db.CheckIfInConversation(database.UserId{userId}, receivers[0])
		if err != nil {
			fmt.Println("error in function CheckIfInConversation")
		}
		fmt.Println(conversationId)
		if conversationId.Id != 0 {
			json.NewEncoder(w).Encode(conversationId)
		} else {
			conversationId, err = rt.db.StartConversation(database.UserId{userId}, receivers)
			json.NewEncoder(w).Encode(conversationId)
		}

	} else {
		conversationId, err = rt.db.StartConversation(database.UserId{userId}, receivers)
		json.NewEncoder(w).Encode(conversationId)
	}

}

func (rt *_router) getConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")
	conversationIdint, err := strconv.Atoi(ps.ByName("conversationId"))
	var conversationId = database.ConversationId{conversationIdint}
	userIdint, err := strconv.Atoi(ps.ByName("Id"))
	var userId = database.UserId{userIdint}

	var userconversations []database.ConversationId
	userconversations, err = rt.db.GetConversations(userId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	count, err := rt.db.CountConversations()
	if err != nil {
		fmt.Println("error in function CountConversations")
	}
	if conversationIdint < 0 || conversationIdint > count {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	for i := 0; i < len(userconversations); i++ {
		if userconversations[i] == conversationId {
			var conversation database.Conversation
			conversation, err = rt.db.GetConversation(conversationId)
			if err != nil {
				fmt.Println("error in function GetConversation")
			}
			for i := 0; i < len(conversation.Content); i++ {
				var message database.Message
				message, err = rt.db.GetMessage(conversation.Content[i])
				if message.Sender != userId && message.Status == "sent" {
					err = rt.db.UpdateStatus(conversation.Content[i])
				}
			}
			if conversation.Groupchat == false {
				var members []string
				members, err = rt.db.GetAllMembers(conversation.Id)
				var memberid database.UserId
				memberid, err = rt.db.GetUserId(members[0])

				if memberid == userId {
					var receiverId database.UserId
					receiverId, err = rt.db.GetUserId(members[1])
					var receiver database.User
					receiver, err = rt.db.GetUser(receiverId)
					conversation.Picture = receiver.Picture
					conversation.Name = receiver.Name
				} else {

					var receiverId database.UserId
					receiverId, err = rt.db.GetUserId(members[0])
					var receiver database.User
					receiver, err = rt.db.GetUser(receiverId)
					conversation.Picture = receiver.Picture
					conversation.Name = receiver.Name
				}
			}
			json.NewEncoder(w).Encode(conversation)
		}
	}

}

func (rt *_router) sendMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")
	conversationIdint, err := strconv.Atoi(ps.ByName("conversationId"))
	var conversationId = database.ConversationId{conversationIdint}
	userIdint, err := strconv.Atoi(ps.ByName("Id"))
	var userId = database.UserId{userIdint}
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	count, err := rt.db.CountConversations()
	if err != nil {
		fmt.Println("error in function CountConversations")
	}
	if conversationIdint < 0 || conversationIdint > count {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var requestData struct {
		Content string `json:"content"`
	}
	err = json.NewDecoder(r.Body).Decode(&requestData)
	var content string = requestData.Content

	var messageId database.MessageId
	messageId, err = rt.db.CreateMessage(userId, conversationId, content)

	json.NewEncoder(w).Encode(messageId)
}
