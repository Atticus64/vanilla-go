package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Note struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

var notes []Note

func main() {
	r := chi.NewRouter()
	r.Get("/", homeHandler)
	r.Get("/notes", getNotes)
	r.Get("/notes/{id}", noteById)
	r.Put("/notes/{id}", UpdateNote)
	r.Post("/notes", addNote)
	r.Delete("/notes/{id}", deleteNote)

	log.Fatal(http.ListenAndServe(":3000", r))
}

func handleErr(err error ) {
	if err != nil {
		panic(err)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome to the API!")
}

func getNotesDB() []Note {
	notesJson, err := os.OpenFile("notes.json", os.O_RDWR|os.O_CREATE, 0644)
	handleErr(err)

	defer notesJson.Close()

	data, err := notesJson.Stat()

	handleErr(err)

	var notes []Note
	if data.Size() != 0 {
		bytes, err := io.ReadAll(notesJson)
		handleErr(err)
		err = json.Unmarshal(bytes, &notes)
		handleErr(err)
	}

	return notes
}

func getNotes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	notesJson, err := os.OpenFile("notes.json", os.O_RDWR|os.O_CREATE, 0644)
	handleErr(err)

	defer notesJson.Close()

	data, err := notesJson.Stat()
	handleErr(err)

	var notes []Note
	if data.Size() != 0 {
		bytes, err := io.ReadAll(notesJson)
		handleErr(err)
		err = json.Unmarshal(bytes, &notes)
		handleErr(err)
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

	note.Id = uuid.NewString()

	notes = append(notes, note)
	file, err := json.MarshalIndent(notes, "", " ")
	handleErr(err)

	err = ioutil.WriteFile("notes.json", file, 0644)
	handleErr(err)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Note added"))
	json.NewEncoder(w).Encode(note)
}

func searchNote(id string, notes []Note) Note {
	for _, n := range notes {
		if n.Id == id {
			return n
		}
	}
	return Note{
		Id: "not found",
	}
}

func removeNoteById(notes []Note, id string) []Note {
	// Creamos un slice de personas vacío
	result := []Note{}

	// Iteramos sobre el slice original y agregamos todas las personas
	// excepto la que tenga el nombre indicado
	for _, note := range notes {
		if note.Id != id {
			result = append(result, note)
		}
	}

	return result
}

func deleteNote(w http.ResponseWriter, r *http.Request) {
	noteId := chi.URLParam(r, "id")

	notes := getNotesDB()

	noteFinded := searchNote(noteId, notes)

	if noteFinded.Id == "not found" {
		w.Write([]byte(fmt.Sprintf("Not found Note with id: %d", &noteId)))
		w.WriteHeader(http.StatusNotFound)
		return
	}

	newNotes := removeNoteById(notes, noteId)

	file, err := json.MarshalIndent(newNotes, "", " ")
	handleErr(err)

	err = ioutil.WriteFile("notes.json", file, 0644)
	handleErr(err)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(newNotes)
}

func noteById(w http.ResponseWriter, r *http.Request) {
	noteId := chi.URLParam(r, "id")

	notes := getNotesDB()

	noteFinded := searchNote(noteId, notes)

	if noteFinded.Id == "not found" {
		w.Write([]byte(fmt.Sprintf("Not found Note with id: %d", &noteId)))
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(notes)
}

func UpdateNote(w http.ResponseWriter, r *http.Request) {
	noteId := chi.URLParam(r, "id")

	notes := getNotesDB()

	noteFinded := searchNote(noteId, notes)

	if noteFinded.Id == "not found" {
		w.Write([]byte(fmt.Sprintf("Not found Note with id: %d", &noteId)))
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var note Note
	err := json.NewDecoder(r.Body).Decode(&note)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	note.Id = noteFinded.Id

	notes = removeNoteById(notes, noteId)

	notes = append(notes, note)

	file, err := json.MarshalIndent(notes, "", " ")
	handleErr(err)
	err = ioutil.WriteFile("notes.json", file, 0644)
	handleErr(err)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(notes)
}
