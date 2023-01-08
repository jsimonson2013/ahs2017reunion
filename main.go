package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/submit", submitForm)
	err := http.ListenAndServe(":3333", nil)
	fmt.Println(err)
}

func submitForm(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Query())

	firstName := r.URL.Query()["fname"][0]
	lastName := r.URL.Query()["lname"][0]
	email := r.URL.Query()["email"][0]
	phone := r.URL.Query()["phone"][0]

	f, err := os.OpenFile("records.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	if _, err = f.WriteString(fmt.Sprintf("%s, %s, %s, %s\n", firstName, lastName, email, phone)); err != nil {
		panic(err)
	}

	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.WriteHeader(200)
}
