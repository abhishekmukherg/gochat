package router

import (
	"github.com/linkinpark342/gchat/users"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"encoding/json"
	"log"
	"strconv"
)

type GchatRouter struct {
	user *users.UserManager
}

func Create(userMgr *users.UserManager) http.Handler {
	r := mux.NewRouter()
	gcr := GchatRouter{userMgr}
	s := r.PathPrefix("/users").Subrouter()
	s.HandleFunc("/", gcr.userAddHandler).Methods("POST")
	//s.HandleFunc("/", UsersHandler)
	s.HandleFunc("/{id:[0-9]+}/", gcr.userGetHandler)
	return r
}

func (gc *GchatRouter) userAddHandler(w http.ResponseWriter, r *http.Request) {
	var u users.User
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	err = json.Unmarshal(buf, &u)
	if err != nil {
		log.Printf("Failed to deserialize: %q\n", err)
		http.Error(w, http.StatusText(400), 400)
		return
	}

	newUser, err := gc.user.Create(u.Name)
	if err != nil {
		log.Printf("Failed to create user: %q\n", err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	b, err := json.Marshal(newUser)
	if err != nil {
		log.Printf("Failed to serialize user %q\n", err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
	w.Header().Add("Content-Type", "text/plain")
	w.WriteHeader(200)
	w.Write(b)
}

func (gc *GchatRouter) userGetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	user, err := gc.user.GetById(id)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	if user == nil {
		http.NotFound(w, r)
		return
	}

	b, err := json.Marshal(user)
	if err != nil {
		log.Printf("Failed to serialize user %q\n", err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
	w.Header().Add("Content-Type", "text/plain")
	w.WriteHeader(200)
	w.Write(b)
}
