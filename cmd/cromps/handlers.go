package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"suah.dev/cromp/db"
)

// NewUser is the handler for /user/new
func NewUser(w http.ResponseWriter, r *http.Request) {
	var u db.CreateUserParams
	json.NewDecoder(r.Body).Decode(&u)

	user, err := base.CreateUser(ctx, u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(user)
}

// Auth is the handler for /usr/auth
func Auth(w http.ResponseWriter, r *http.Request) {
	var p db.AuthUserParams
	json.NewDecoder(r.Body).Decode(&p)

	user, err := base.AuthUser(ctx, p)
	if err != nil {
		panic(err)
	}

	if user.Authed {
		authedUsers[user.Token.String()] = user
		// TODO respond with token
		json.NewEncoder(w).Encode(user)
	}
}

// AddEntry handles requests to /entries/add
func AddEntry(w http.ResponseWriter, r *http.Request) {
	var e db.CreateEntryParams

	json.NewDecoder(r.Body).Decode(&e)

	entry, err := base.CreateEntry(ctx, e)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(entry)
}

// Entries handles requests to /entries
func Entries(w http.ResponseWriter, r *http.Request) {
	// Print secret message
	fmt.Fprintln(w, "The cake is a lie!")
}

// SimilarEntries
func SimilarEntries(w http.ResponseWriter, r *http.Request) {
}
