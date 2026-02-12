package api

import (
	"encoding/json"
	"net/http"

	"fmt"
	"strconv"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
	"github.com/julienschmidt/httprouter"
)



func ( rt *_router) doLogin(w http. ResponseWriter, r *http.Request, ps httprouter.Params){
	w.Header().Set("content-type", "application/json")
	var id database.UserId
	var name string = r.URL.Query().Get("")
	exists, err := rt.db.UserLookup(name)
	if err != nil {
		fmt.Println("error in function UserLookup")
	}
	if exists == false{
		err = rt.db.AddUser(name)
		if err != nil {
			fmt.Println("error in function AddUser")
		}
	}
	id, err = rt.db.GetUserId(name)
	if err != nil {
		fmt.Println("error in function GetUserId")
	}
	json.NewEncoder(w).Encode(id.Id)

}

func (rt *_router) getUsers(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")
	userId,err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	count, err := rt.db.CountUsers()
	if err!= nil{
		fmt.Println("error in function CountUsers")
	}
	if userId <0 || userId > count{
		w.WriteHeader(http.StatusNotFound)
		return
	}
	users, err := rt.db.ListAllUsers()
	if err!= nil{
		fmt.Println("error in function ListAllUsers")
	}
	json.NewEncoder(w).Encode(users)
}

func (rt *_router) getSpecificUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params){
	w.Header().Set("content-type", "application/json")
	userId,err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	count, err := rt.db.CountUsers()
	if err!= nil{
		fmt.Println("error in function CountUsers")
	}
	if userId <0 || userId > count{
		w.WriteHeader(http.StatusNotFound)
		return
	}
	var name string = r.URL.Query().Get("name")
	id, err := rt.db.GetUserId(name)
	if err != nil{
		w.WriteHeader(http.StatusNotFound)
	}
	var user database.User
	user, err = rt.db.GetUser(id)
	if err != nil{
		fmt.Println("function getUser failed")
	}
	json.NewEncoder(w).Encode(user)
}

func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params){
	w.Header().Set("content-type", "application/json")
	userId,err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	count, err := rt.db.CountUsers()
	if err!= nil{
		fmt.Println("error in function CountUsers")
	}
	if userId <0 || userId > count{
		w.WriteHeader(http.StatusNotFound)
		return
	}
	var name string = r.URL.Query().Get("")
	var id database.UserId = database.UserId{userId}
	rt.db.ChangeUserName(id,name)
}


func (rt *_router) setMyPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params){
	w.Header().Set("content-type", "application/json")
	userId,err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	count, err := rt.db.CountUsers()
	if err!= nil{
		fmt.Println("error in function CountUsers")
	}
	if userId <0 || userId > count{
		w.WriteHeader(http.StatusNotFound)
		return
	}
	var picture string = r.URL.Query().Get("")
	var id database.UserId = database.UserId{userId}
	rt.db.ChangeUserName(id,picture)
}
