package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

type DBHandler struct {
	DB *sql.DB
}

func main() {
	dbName := "reunion"
	user := os.Args[1]
	pass := os.Args[2]
	protocol := "tcp"
	port := os.Args[3]

	DB, err := sql.Open("mysql", fmt.Sprintf("%s:%s@%s(127.0.0.1:3306)/%s", user, pass, protocol, dbName))
	if err != nil {
		panic(err)
	}
	defer DB.Close()

	if err := DB.Ping(); err != nil {
		panic(err)
	}

	dbh := DBHandler{DB}

	fmt.Println("Connected to the DB...")

	http.HandleFunc("/submit", submitForm)
	http.HandleFunc("/rsvpWithToken", dbh.rsvp)
	http.HandleFunc("/submitRSVP", dbh.submitRSVP)
	err = http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
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

type Token struct {
	ID         *int64
	ContactID  *int64
	Token      *string
	Expiration *int64
}

func (h DBHandler) rsvp(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")

	// check that request has token
	token := r.URL.Query()["token"]
	if len(token) < 1 {
		w.WriteHeader(403)
		return
	}

	// A token to hold data from the returned row.
	var tok Token

	row := h.DB.QueryRow("SELECT ID, ContactID, Token, expiration FROM reunion.tokens WHERE Token = ?", token[0])
	if err := row.Scan(&tok.ID, &tok.ContactID, &tok.Token, &tok.Expiration); err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(404)
			return
		}
		if err != nil {
			w.WriteHeader(500)
			fmt.Println(err)
			return
		}
	}

	if time.Now().Unix() > *tok.Expiration {
		w.WriteHeader(401)
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

func (h DBHandler) submitRSVP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")

	// extract params from url
	if len(r.URL.Query()["name"]) < 1 ||
		len(r.URL.Query()["attending"]) < 1 ||
		len(r.URL.Query()["plusone"]) < 1 ||
		len(r.URL.Query()["token"]) < 1 {
		w.WriteHeader(400)
		return
	}

	// insert into rsvp's db
	name := r.URL.Query()["name"][0]
	attending := r.URL.Query()["attending"][0] == "y"
	plusone := r.URL.Query()["plusone"][0] == "y"
	token := r.URL.Query()["token"][0]

	// A token to hold data from the returned row.
	var tok Token

	row := h.DB.QueryRow("SELECT ID, ContactID, Token, expiration FROM reunion.tokens WHERE Token = ?", token)
	if err := row.Scan(&tok.ID, &tok.ContactID, &tok.Token, &tok.Expiration); err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(404)
			return
		}
		if err != nil {
			w.WriteHeader(500)
			fmt.Println(err)
			return
		}
	}

	_, err := h.DB.Exec("INSERT INTO reunion.rsvps (ContactID, Name, Attending, PlusOne) VALUES (?, ?, ?, ?)", tok.ContactID, name, attending, plusone)
	if err != nil {
		w.WriteHeader(500)
		fmt.Println(err)
		return
	}

	// route to rsvp form with token information added to url
	f, err := os.Open("./index.html")
	if err != nil {
		w.WriteHeader(500)
		fmt.Println(err)
		return
	}
	defer f.Close()

	bs, err := ioutil.ReadAll(f)
	if err != nil {
		w.WriteHeader(500)
		fmt.Println(err)
		return
	}

	fmt.Fprint(w, string(bs))

	if _, err := h.DB.Exec("UPDATE reunion.tokens SET expiration=?", 0); err != nil {
		fmt.Printf("Error while updating token expiration: %v\n", err)
	}

	fmt.Println("HERE")
}
