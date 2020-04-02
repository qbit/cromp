package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"suah.dev/cromp/db"
)

var (
	pd, err     = sql.Open("postgres", "host=localhost dbname=qbit sslmode=disable password=''")
	ctx, cancel = context.WithCancel(context.Background())
	base        = db.New(pd)
)

func logger(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s\n", r.URL.Path)
		f(w, r)
	}
}

func getUser(r *http.Request) (*db.User, error) {
	token := r.Header.Get("X-Access-Token")
	if token == "" {
		log.Printf("checkAuth: %s received empty token\n", r.URL.Path)
		return nil, fmt.Errorf("Unauthorized")
	}

	u, err := uuid.Parse(token)
	if err != nil {
		return nil, err
	}

	user, err := base.GetUserByToken(ctx, u)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func checkAuth(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := getUser(r)
		if err != nil {
			http.Error(w, "Invalid User", http.StatusUnauthorized)
			return
		}

		if time.Now().Before(user.TokenExpires) {
			log.Printf("checkAuth: %s received valid token\n", r.URL.Path)
			f(w, r)
		} else {
			log.Printf("checkAuth: %s received expired token\n", r.URL.Path)
			http.Error(w, "Token Expired", http.StatusUnauthorized)
			return
		}
	}
}

func main() {

	defer cancel() // cancel when we are finished consuming integers

	if err != nil {
		panic(err)
	}
	defer pd.Close()

	http.HandleFunc("/user/new", logger(NewUser))
	http.HandleFunc("/user/auth", logger(Auth))

	http.HandleFunc("/entries/add", checkAuth(logger(AddEntry)))
	http.HandleFunc("/entries/list", checkAuth(logger(ListEntries)))
	http.HandleFunc("/entries/delete", checkAuth(logger(Entries)))
	http.HandleFunc("/entries/get", checkAuth(logger(Entries)))
	http.HandleFunc("/entries/update", checkAuth(logger(Entries)))
	http.HandleFunc("/entries/similar", checkAuth(logger(SimilarEntries)))

	log.Fatalln(http.ListenAndServe(":8080", nil))
}
