package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func main() {
	dbName := "reunion"
	user := ""
	pass := ""
	protocol := "tcp"

	DB, err := sql.Open("mysql", fmt.Sprintf("%s:%s@%s(127.0.0.1:3306)/%s", user, pass, protocol, dbName))
	if err != nil {
		panic(err)
	}
	defer DB.Close()

	http.HandleFunc("/submit", submitForm)
	http.HandleFunc("/rsvp", rsvp)
	http.HandleFunc("/submit/rsvp", submitRSVP)
	err = http.ListenAndServe(":3333", nil)
	fmt.Println(err)
}

func submitForm(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")

	if len(r.URL.Query()["fname"]) < 1 ||
		len(r.URL.Query()["lname"]) < 1 ||
		len(r.URL.Query()["email"]) < 1 ||
		len(r.URL.Query()["phone"]) < 1 ||
		len(r.URL.Query()["method"]) < 1 {
		w.WriteHeader(400)
		return
	}

	firstName := r.URL.Query()["fname"][0]
	lastName := r.URL.Query()["lname"][0]
	email := r.URL.Query()["email"][0]
	phone := r.URL.Query()["phone"][0]
	method := r.URL.Query()["method"][0]

	/* TODO
	sql := "INSERT INTO cities(name, population) VALUES ('Moscow', 12506000)"
	res, err := db.Exec(sql)
	*/

	f, err := os.OpenFile("records.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	if _, err = f.WriteString(fmt.Sprintf("%s, %s, %s, %s, %s\n", firstName, lastName, email, phone, method)); err != nil {
		panic(err)
	}

	w.WriteHeader(200)
}

func rsvp(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")

	// check that request has token
	token := r.URL.Query()["token"]
	if len(token) < 1 {
		w.WriteHeader(403)
		return
	}

	/* TODO
	sql := "INSERT INTO cities(name, population) VALUES ('Moscow', 12506000)"
	res, err := db.Exec(sql)
	*/
	// check that token exists in tokens db and is not expire
	if token[0] != "1234" {
		w.WriteHeader(403)
		return
	}

	// route to rsvp form with token information added to url
	f, err := os.Open("./rsvp/index.html")
	if err != nil {
		w.WriteHeader(500)
		return
	}
	defer f.Close()

	bs, err := ioutil.ReadAll(f)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	fmt.Fprint(w, string(bs))
}

func submitRSVP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")

	// extract params from url
	if len(r.URL.Query()["name"]) < 1 ||
		len(r.URL.Query()["attending"]) < 1 ||
		len(r.URL.Query()["plusone"]) < 1 {
		w.WriteHeader(400)
		return
	}

	// insert into rsvp's db
	name := r.URL.Query()["name"][0]
	attending := r.URL.Query()["attending"][0]
	plusone := r.URL.Query()["plusone"][0]

	/* TODO
	sql := "INSERT INTO cities(name, population) VALUES ('Moscow', 12506000)"
	res, err := db.Exec(sql)
	*/

	f, err := os.OpenFile("rsvps.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	defer f.Close()

	if _, err = f.WriteString(fmt.Sprintf("%s, %s, %s\n", name, attending, plusone)); err != nil {
		w.WriteHeader(500)
		return
	}

	// redirect to home page
	w.WriteHeader(200)
}
