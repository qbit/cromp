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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if user.Authed {
		// TODO respond with token
		json.NewEncoder(w).Encode(user)
	}
}

// ListEntries handles requests to /entries/list
func ListEntries(w http.ResponseWriter, r *http.Request) {
	user, err := getUser(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println(user.UserID)
	entries, err := base.GetEntries(ctx, user.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(entries)
}

// AddEntry handles requests to /entries/add
func AddEntry(w http.ResponseWriter, r *http.Request) {
	var e db.CreateEntryParams

	err := json.NewDecoder(r.Body).Decode(&e)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user, err := getUser(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	e.UserID = user.UserID

	entry, err := base.CreateEntry(ctx, e)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(entry)
}

// UpdateEntries handles requests to /etries/update
func UpdateEntries(w http.ResponseWriter, r *http.Request) {
	var params db.UpdateEntryParams

	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user, err := getUser(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	params.UserID = user.UserID

	entry, err := base.UpdateEntry(ctx, params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(entry)
}

// Entries handles requests to /entries
func Entries(w http.ResponseWriter, r *http.Request) {
	var params db.GetEntryParams

	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user, err := getUser(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	params.UserID = user.UserID

	entry, err := base.GetEntry(ctx, params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(entry)
}

// SimilarEntries are entries that match some text
func SimilarEntries(w http.ResponseWriter, r *http.Request) {
	var e db.SimilarEntriesParams
	err := json.NewDecoder(r.Body).Decode(&e)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user, err := getUser(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	e.UserID = user.UserID
	entries, err := base.SimilarEntries(ctx, e)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(entries)
}
