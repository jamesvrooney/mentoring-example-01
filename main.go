package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

type Person struct {
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
	Email     string `json:"email"`
}

var people = map[string]Person{
	"joe": {
		FirstName: "joe",
		LastName:  "strong",
		Email:     "joe.strong@test.com",
	}, "mike": {
		FirstName: "mike",
		LastName:  "strong",
		Email:     "mike.strong@test.com",
	}, "pete": {
		FirstName: "pete",
		LastName:  "strong",
		Email:     "pete.strong@test.com",
	},
}

func main() {
	http.HandleFunc("/foo", fooHandler)
	http.HandleFunc("/bar", barHandler)
	http.HandleFunc("/person", personHandler)

	go http.ListenAndServe(":8080", nil)

	time.Sleep(300 * time.Millisecond)

	response, err := http.Get("http://localhost:8080/person?user=joe")

	if err != nil {
		panic(err)
	}

	var person Person

	if err := json.NewDecoder(response.Body).Decode(&person); err != nil {
		panic(err)
	}

	fmt.Printf("%#v\n", person)
}

func personHandler(w http.ResponseWriter, r *http.Request) {
	user := r.URL.Query().Get("user")

	person := people[user]

	if err := json.NewEncoder(w).Encode(person); err != nil {
		http.Error(w, "Error encoding", http.StatusInternalServerError)
	}
}

func barHandler(w http.ResponseWriter, r *http.Request) {
	_, err := io.Copy(w, r.Body)

	if err != nil {
		http.Error(w, "Failed reading body", http.StatusInternalServerError)

		return
	}
}

func fooHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		http.Error(w, "Failed reading body", http.StatusInternalServerError)

		return
	}

	if _, err := w.Write(body); err != nil {
		http.Error(w, "Failed writing body", http.StatusInternalServerError)
	}
}
