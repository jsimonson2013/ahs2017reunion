package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/submit", submitForm)
	http.HandleFunc("/rsvp", rsvp)
	http.HandleFunc("/submit/rsvp", submitRSVP)
	err := http.ListenAndServe(":3333", nil)
	fmt.Println(err)
}

func submitForm(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Query())
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
	// check that request has token

	// check that token exists in tokens db and is not expired

	// route to rsvp form with token information added to url
}

func submitRSVP(w http.ResponseWriter, r *http.Request) {
	// extract params from url

	// insert into rsvp's db

	// redirect to home page
}
