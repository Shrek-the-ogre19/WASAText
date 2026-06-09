package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/Shrek-the-ogre19/WASAText/service/database"
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
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	if userId < 0 || userId > count {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	conversationIds, err := rt.db.GetConversations(database.UserId{userId})
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	var conversations []database.Conversation
	for i := 0; i < len(conversationIds); i++ {
		toAdd, err := rt.db.GetConversation(conversationIds[i])
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			return
		}
		if toAdd.Groupchat == false {
			members, err := rt.db.GetAllMembers(conversationIds[i])
			if err != nil {
				w.WriteHeader(http.StatusBadGateway)
				return
			}
			memberid, err := rt.db.GetUserId(members[0])
			if err != nil {
				w.WriteHeader(http.StatusBadGateway)
				return
			}

			if memberid.Id == userId {
				receiverId, err := rt.db.GetUserId(members[1])
				if err != nil {
					w.WriteHeader(http.StatusBadGateway)
					return
				}
				receiver, err := rt.db.GetUser(receiverId)
				if err != nil {
					w.WriteHeader(http.StatusBadGateway)
					return
				}
				toAdd.Picture = receiver.Picture
				toAdd.Name = receiver.Name
			} else {
				receiverId, err := rt.db.GetUserId(members[0])
				if err != nil {
					w.WriteHeader(http.StatusBadGateway)
					return
				}
				receiver, err := rt.db.GetUser(receiverId)
				if err != nil {
					w.WriteHeader(http.StatusBadGateway)
					return
				}
				toAdd.Picture = receiver.Picture
				toAdd.Name = receiver.Name
			}
		}
		conversations = append(conversations, toAdd)

	}

	err = json.NewEncoder(w).Encode(conversations)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
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
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	if userId < 0 || userId > count {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var requestData struct {
		Receivers string `json:"receivers"`
	}
	err = json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	receivers := strings.Split(requestData.Receivers, ",")

	var conversationId database.ConversationId
	for _, receiver := range receivers {
		exists, err := rt.db.UserLookup(receiver)
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			return
		}
		if exists == false {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	if len(receivers) == 1 {
		conversationId, err = rt.db.CheckIfInConversation(database.UserId{userId}, receivers[0])
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			return
		}
		if conversationId.Id != 0 {
			err = json.NewEncoder(w).Encode(conversationId)
			if err != nil {
				w.WriteHeader(http.StatusBadGateway)
				return
			}
		} else {
			conversationId, err = rt.db.StartConversation(database.UserId{userId}, receivers)
			if err != nil {
				w.WriteHeader(http.StatusBadGateway)
				return
			}
			err = json.NewEncoder(w).Encode(conversationId)
			if err != nil {
				w.WriteHeader(http.StatusBadGateway)
				return
			}
		}

	} else {
		conversationId, err = rt.db.StartConversation(database.UserId{userId}, receivers)
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			return
		}
		err = json.NewEncoder(w).Encode(conversationId)
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			return
		}
	}

}

func (rt *_router) getConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")
	conversationIdint, err := strconv.Atoi(ps.ByName("conversationId"))
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	conversationId := database.ConversationId{conversationIdint}
	userIdint, err := strconv.Atoi(ps.ByName("Id"))
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	userId := database.UserId{userIdint}

	userconversations, err := rt.db.GetConversations(userId)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	count, err := rt.db.CountConversations()
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	if conversationIdint < 0 || conversationIdint > count {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	for i := 0; i < len(userconversations); i++ {
		if userconversations[i] == conversationId {
			conversation, err := rt.db.GetConversation(conversationId)
			if err != nil {
				w.WriteHeader(http.StatusBadGateway)
				return
			}
			for i := 0; i < len(conversation.Content); i++ {
				message, err := rt.db.GetMessage(conversation.Content[i])
				if err != nil {
					w.WriteHeader(http.StatusBadGateway)
					return
				}
				if message.Sender != userId && message.Status == "sent" {
					err = rt.db.UpdateStatus(conversation.Content[i])
					if err != nil {
						w.WriteHeader(http.StatusBadGateway)
						return
					}
				}
			}
			if conversation.Groupchat == false {
				members, err := rt.db.GetAllMembers(conversation.Id)
				if err != nil {
					w.WriteHeader(http.StatusBadGateway)
					return
				}
				memberid, err := rt.db.GetUserId(members[0])
				if err != nil {
					w.WriteHeader(http.StatusBadGateway)
					return
				}

				if memberid == userId {
					receiverId, err := rt.db.GetUserId(members[1])
					if err != nil {
						w.WriteHeader(http.StatusBadGateway)
						return
					}
					receiver, err := rt.db.GetUser(receiverId)
					if err != nil {
						w.WriteHeader(http.StatusBadGateway)
						return
					}
					conversation.Picture = receiver.Picture
					conversation.Name = receiver.Name
				} else {
					receiverId, err := rt.db.GetUserId(members[0])
					if err != nil {
						w.WriteHeader(http.StatusBadGateway)
						return
					}
					receiver, err := rt.db.GetUser(receiverId)
					if err != nil {
						w.WriteHeader(http.StatusBadGateway)
						return
					}
					conversation.Picture = receiver.Picture
					conversation.Name = receiver.Name
				}
			}
			err = json.NewEncoder(w).Encode(conversation)
			if err != nil {
				w.WriteHeader(http.StatusBadGateway)
				return
			}

		}
	}

}

func (rt *_router) sendMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")
	conversationIdint, err := strconv.Atoi(ps.ByName("conversationId"))
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	conversationId := database.ConversationId{conversationIdint}
	userIdint, err := strconv.Atoi(ps.ByName("Id"))
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	userId := database.UserId{userIdint}
	count, err := rt.db.CountConversations()
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	if conversationIdint < 0 || conversationIdint > count {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var requestData struct {
		Content string `json:"content"`
	}
	err = json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	messageId, err := rt.db.CreateMessage(userId, conversationId, requestData.Content)
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
