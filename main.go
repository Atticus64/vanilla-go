package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Note struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

var notes []Note

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/notes", notesHandler)

	log.Fatal(http.ListenAndServe(":3000", nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome to the API!")
}

func notesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getNotes(w, r)
	case http.MethodPost:
		addNote(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getNotes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	notesJson, err := ioutil.ReadFile("notes.json")
	if err != nil {
		ioutil.WriteFile("notes.json", []byte("[]"), 0644)
	}

	notes := []Note{}
	if notesJson != nil {
		json.Unmarshal(notesJson, &notes)
	}

	json.NewEncoder(w).Encode(notes)
}

func addNote(w http.ResponseWriter, r *http.Request) {
	var note Note
	err := json.NewDecoder(r.Body).Decode(&note)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	notes = append(notes, note)
	fmt.Println(notes)
	file, _ := json.MarshalIndent(notes, "", " ")
	_ = ioutil.WriteFile("notes.json", file, 0644)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Note added"))
}
