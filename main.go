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
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	stan "github.com/nats-io/stan.go"
)

//User is a structur for users
type User struct {
	UserID              int    `json:"userid"`
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
	Updated_At          string `json:"updated_at"`
}

var db *sqlx.DB
var err error

func main() {
	db, err = sqlx.Connect("mysql", "root:userAccess1@tcp(127.0.0.1:3306)/testdb?multiStatements=true")
	defer func() {
		err := db.Close()
		if err != nil {
			log.Println(err.Error())
		}
	}()

	driver, err := mysql.WithInstance(db.DB, &mysql.Config{})
	if err != nil {
		log.Fatal(err.Error())
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations",
		"mysql",
		driver,
	)
	if err != nil {
		log.Fatal(err.Error())
	}

	err = m.Up()
	if err != nil {
		log.Println(err.Error())
	}

	router := mux.NewRouter()
	router.HandleFunc("/users/{id}", getUser).Methods("GET")
	router.HandleFunc("/users/add", addUser).Methods("POST")

	port := ":333"
	http.ListenAndServe(port, router)
	err = http.ListenAndServe(port, router)
	if err != nil && err != http.ErrServerClosed {
		log.Panicln(err.Error())
	}
}

func getUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	getID := mux.Vars(r)

	result, err := db.Queryx("SELECT * FROM users WHERE userid=?", getID["id"])
	if err != nil {
		log.Fatal(err.Error())
	}

	var user User

	for result.Next() {
		err := result.Scan(&user.UserID, &user.Login, &user.ID, &user.Node_ID, &user.Avatar_URL,
			&user.Gravatar_ID, &user.URL, &user.HTML_URL, &user.Followers_URL, &user.Following_URL,
			&user.Gists_URL, &user.Starred_URL, &user.Subscriptions_URL, &user.Organizations_URL,
			&user.Repos_URL, &user.Events_URL, &user.Received_Events_URL, &user.Type, &user.Site_Admin,
			&user.Name, &user.Company, &user.Blog, &user.Location, &user.Email, &user.Hireable, &user.Bio,
			&user.Public_Repos, &user.Public_Gists, &user.Followers, &user.Following, &user.Created_At,
			&user.Updated_At)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
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

	_, err = db.NamedExec(`INSERT INTO users VALUES (:userid, :login, :id, :node_id, :avatar_url, 
		:gravatar_id, :url, :html_url, :followers_url, :following_url, :gists_url, :starred_url, 
		:subscriptions_url, :organizations_url, :repos_url, :events_url, :received_events_url, 
		:type, :site_admin, :name, :company, :blog, :location, :email, :hireable, :bio, :public_repos, 
		:public_gists, :followers, :following, :created_at, :updated_at)`, user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, _ = fmt.Fprintf(w, "New data was append")

	sc, err := stan.Connect("test-cluster", "testID")
	if err != nil {
		log.Fatal(err)
	}

	defer sc.Close()

	m := &user

	me, err := json.Marshal(m)
	if err != nil {
		log.Fatal(err)
	}

	if err := sc.Publish("foo", me); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Publish message")
}
