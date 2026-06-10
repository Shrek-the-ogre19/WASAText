package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Shrek-the-ogre19/WASAText/service/database"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) setGroupName(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var requestData struct {
		Newname string `json:"name"`
	}
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	conversationIdint, err := strconv.Atoi(ps.ByName("conversationId"))
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	conversationId := database.ConversationId{conversationIdint}

	groupchat, err := rt.db.CheckGroupchat(conversationId)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	if groupchat {
		err = rt.db.ChangeConversationName(conversationId, requestData.Newname)
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			return
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (rt *_router) setGroupPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var requestData struct {
		Newpicture string `json:"picture"`
	}
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	conversationIdint, err := strconv.Atoi(ps.ByName("conversationId"))
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	conversationId := database.ConversationId{conversationIdint}

	groupchat, err := rt.db.CheckGroupchat(conversationId)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	if groupchat {
		err = rt.db.ChangeConversationPhoto(conversationId, requestData.Newpicture)
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			return
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (rt *_router) listGroupMembers(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")
	conversationIdint, err := strconv.Atoi(ps.ByName("conversationId"))
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	conversationId := database.ConversationId{conversationIdint}

	groupchat, err := rt.db.CheckGroupchat(conversationId)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	if groupchat {
		members, err := rt.db.GetAllMembers(conversationId)
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			return
		}
		err = json.NewEncoder(w).Encode(members)
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			return
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (rt *_router) addToGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")
	conversationIdint, err := strconv.Atoi(ps.ByName("conversationId"))
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	conversationId := database.ConversationId{conversationIdint}
	groupchat, err := rt.db.CheckGroupchat(conversationId)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	if !groupchat {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	members, err := rt.db.GetAllMembers(conversationId)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	var requestData struct {
		Name string `json:"name"`
	}
	err = json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	memberToAdd := requestData.Name

	for i := 0; i < len(members); i++ {
		if memberToAdd == members[i] {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	exists, err := rt.db.UserLookup(memberToAdd)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	if !exists {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err = rt.db.AddMembers(conversationId, memberToAdd)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(memberToAdd)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
}

func (rt *_router) leaveGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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
	groupchat, err := rt.db.CheckGroupchat(conversationId)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	if groupchat {
		err = rt.db.LeaveGroup(conversationId, userId)
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			return
		}
	}

}
