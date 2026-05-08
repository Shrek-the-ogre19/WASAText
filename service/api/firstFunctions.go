package api

import (
	"encoding/json"
	"net/http"

	"strconv"

	"github.com/Shrek-the-ogre19/WASAText/service/database"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")
	var id database.UserId
	var requestData struct {
		Name string `json:"name"`
	}
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	var name string = requestData.Name

	exists, err := rt.db.UserLookup(name)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	if exists == false {
		err = rt.db.AddUser(name)
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			return
		}
	}
	id, err = rt.db.GetUserId(name)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	err = json.NewEncoder(w).Encode(id.Id)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}

}

func (rt *_router) getSelf(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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
	user, err := rt.db.GetUser(database.UserId{userId})
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
}

func (rt *_router) getUsers(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")
	users, err := rt.db.ListAllUsers()
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
}

func (rt *_router) getSpecificUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")
	id, err := strconv.Atoi(ps.ByName("specificId"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var user database.User
	user, err = rt.db.GetUser(database.UserId{id})
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
}

func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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
		Name string `json:"name"`
	}
	err = json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var name string = requestData.Name

	users, err := rt.db.ListAllUsers()
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	for _, user := range users {
		// Reject duplicate usernames owned by other users.
		if user == name {
			http.Error(w, `{"error":"username already exists"}`, http.StatusBadRequest)
			return
		}
	}

	var id database.UserId = database.UserId{userId}
	err = rt.db.ChangeUserName(id, name)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
}

func (rt *_router) setMyPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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
		Picture string `json:"photo"`
	}
	err = json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	var picture string = requestData.Picture

	var id database.UserId = database.UserId{userId}
	err = rt.db.ChangeUserPicture(id, picture)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
}
