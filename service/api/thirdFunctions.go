package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) setGroupName(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var requestData struct {
		Newname string `json:"name"`
	}
	err := json.NewDecoder(r.Body).Decode(&requestData)
	var newname string = requestData.Newname

	conversationIdint, err := strconv.Atoi(ps.ByName("conversationId"))
	var conversationId = database.ConversationId{conversationIdint}

	var groupchat bool
	groupchat, err = rt.db.CheckGroupchat(conversationId)
	if err != nil {
		fmt.Println("error in function CheckGroupchat")
	}
	if groupchat {
		err = rt.db.ChangeConversationName(conversationId, newname)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (rt *_router) setGroupPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var requestData struct {
		Newpicture string `json:"picture"`
	}
	err := json.NewDecoder(r.Body).Decode(&requestData)
	var newpicture string = requestData.Newpicture

	conversationIdint, err := strconv.Atoi(ps.ByName("conversationId"))
	var conversationId = database.ConversationId{conversationIdint}

	var groupchat bool
	groupchat, err = rt.db.CheckGroupchat(conversationId)
	if err != nil {
		fmt.Println("error in function CheckGroupchat")
	}
	if groupchat {
		err = rt.db.ChangeConversationPhoto(conversationId, newpicture)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (rt *_router) listGroupMembers(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")
	conversationIdint, err := strconv.Atoi(ps.ByName("conversationId"))
	var conversationId = database.ConversationId{conversationIdint}

	var groupchat bool
	groupchat, err = rt.db.CheckGroupchat(conversationId)
	if err != nil {
		fmt.Println("error in function CheckGroupchat")
	}
	if groupchat {
		var members []string
		members, err = rt.db.GetAllMembers(conversationId)
		json.NewEncoder(w).Encode(members)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (rt *_router) addToGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	conversationIdint, err := strconv.Atoi(ps.ByName("conversationId"))
	var conversationId = database.ConversationId{conversationIdint}
	var groupchat bool
	groupchat, err = rt.db.CheckGroupchat(conversationId)
	if err != nil {
		fmt.Println("error in function CheckGroupchat")
	}
	if groupchat {
		var members []string
		members, err = rt.db.GetAllMembers(conversationId)
		var requestData struct {
			Name string `json:"name"`
		}
		err = json.NewDecoder(r.Body).Decode(&requestData)
		var memberToAdd string = requestData.Name

		for i := 0; i < len(members); i++ {
			if memberToAdd == members[i] {
				w.WriteHeader(http.StatusBadRequest)
				break
			}
		}
		var users []string
		users, err = rt.db.ListAllUsers()
		for i := 0; i < len(users); i++ {
			if memberToAdd == users[i] {
				err = rt.db.AddMembers(conversationId, memberToAdd)
				if err != nil {
					fmt.Println("error in function AddMembers")
					break
				}
			}

		}
		w.WriteHeader(http.StatusBadRequest)

	} else {
		w.WriteHeader(http.StatusBadRequest)
	}

}

func (rt *_router) leaveGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	conversationIdint, err := strconv.Atoi(ps.ByName("conversationId"))
	var conversationId = database.ConversationId{conversationIdint}
	userIdint, err := strconv.Atoi(ps.ByName("Id"))
	var userId = database.UserId{userIdint}
	var groupchat bool
	groupchat, err = rt.db.CheckGroupchat(conversationId)
	if err != nil {
		fmt.Println("error in function CheckGroupchat")
	}
	if groupchat {
		err = rt.db.LeaveGroup(conversationId, userId)
	}

}
