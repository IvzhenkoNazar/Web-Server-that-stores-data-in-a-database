package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/mysql"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

//User is a structur for users
type User struct {
	Login               string `json:"login"`
	ID                  int    `json:"id"`
	Node_ID             string `json:"node_id"`
	Avatar_URL          string `json:"avatar_url"`
	Gravatar_ID         string `json:"gravatar_id"`
	URL                 string `json:"url"`
	HTML_URL            string `json:"html_url"`
	Followers_URL       string `json:"followers_url"`
	Following_URL       string `json:"following_url"`
	Gists_URL           string `json:"gists_url"`
	Starred_URL         string `json:"starred_url"`
	Subscriptions_URL   string `json:"subscriptions_url"`
	Organizations_URL   string `json:"organizations_url"`
	Repos_URL           string `json:"repos_url"`
	Events_URL          string `json:"events_url"`
	Received_Events_URL string `json:"received_events_url"`
	Type                string `json:"type"`
	Site_Admin          bool   `json:"site_admin"`
	Name                string `json:"name"`
	Company             string `json:"company"`
	Blog                string `json:"blog"`
	Location            string `json:"location"`
	Email               string `json:"email"`
	Hireable            string `json:"hireable"`
	Bio                 string `json:"bio"`
	Public_Repos        int    `json:"public_repos"`
	Public_Gists        int    `json:"public_gists"`
	Followers           int    `json:"followers"`
	Following           int    `json:"following"`
	Created_At          string `json:"created_at"`
	Update_At           string `json:"update_at"`
}

var db *sqlx.DB
var err error

func check(e error) {
	if e != nil {
		log.Fatal(e.Error())
	}
}

func main() {
	db, err = sqlx.Connect("mysql", "root:userAccess1@tcp(127.0.0.1:3306)/testdb")
	check(err)

	defer func() {
		err := db.Close()
		if err != nil {
			log.Println(err.Error())
		}
	}()

	driver, err := mysql.WithInstance(db.DB, &mysql.Config{})
	check(err)

	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations",
		"mysql",
		driver,
	)
	check(err)

	err = m.Up()
	check(err)

	router := mux.NewRouter()
	router.HandleFunc("/users/{id}", getUser).Methods("GET")
	router.HandleFunc("/users/add", addUser).Methods("POST")

	port := ":333"
	err = http.ListenAndServe(port, router)
	if err != nil && err != http.ErrServerClosed {
		log.Panicln(err.Error())
	}
}

func getUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	getID := mux.Vars(r)

	result, err := db.Queryx("SELECT * FROM users WHERE id=?", getID["id"])
	check(err)

	var user User

	for result.Next() {
		err := result.Scan(&user.Login, &user.ID, &user.Node_ID, &user.Avatar_URL, &user.Gravatar_ID,
			&user.URL, &user.HTML_URL, &user.Followers_URL, &user.Following_URL, &user.Gists_URL,
			&user.Starred_URL, &user.Subscriptions_URL, &user.Organizations_URL, &user.Repos_URL,
			&user.Events_URL, &user.Received_Events_URL, &user.Type, &user.Site_Admin, &user.Name,
			&user.Company, &user.Blog, &user.Location, &user.Email, &user.Hireable, &user.Bio,
			&user.Public_Repos, &user.Public_Gists, &user.Followers, &user.Following, &user.Created_At,
			&user.Update_At)
		check(err)
	}

	json.NewEncoder(w).Encode(user)
}

func addUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user User

	dat, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(dat, &user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = db.NamedExec(`INSERT INTO users VALUES (:login, :id, :node_id, :avatar_url, :gravatar_id, 
		:url, :html_url, :followers_url, :following_url, :gists_url, :starred_url, 
		:subscriptions_url, :organizations_url, :repos_url, :events_url, :received_events_url, 
		:type, :site_admin, :name, :company, :blog, :location, :email, :hireable, :bio, :public_repos, 
		:public_gists, :followers, :following, :created_at, :update_at)`, user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, _  = fmt.Fprintf(w, "New data was append")
}
