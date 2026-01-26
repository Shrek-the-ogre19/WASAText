package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

var AllUsers = []User{}

func (rt *_router) getUsers(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(AllUsers)
}
