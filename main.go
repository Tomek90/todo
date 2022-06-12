package main

import (
	"net/http"
	"todo/application"
	ht "todo/handlers"

	"github.com/gorilla/mux"
)

func main() {
	app := &application.Application{}

	err := app.StartApp()
	if err != nil {
		panic(err)
	}

	baseHandler := ht.NewAppHandler(app)
	r := mux.NewRouter()
	r.HandleFunc("/", baseHandler.MainPage).Methods("GET")
	r.HandleFunc("/", baseHandler.CreateToDo).Methods("POST")
	r.HandleFunc("/{id}/delete", baseHandler.DeleteToDo).Methods("GET")
	r.HandleFunc("/{id}/update", baseHandler.UpdateToDo).Methods("POST")
	r.HandleFunc("/api/alltodos", baseHandler.GetAllAPI).Methods("GET")
	r.HandleFunc("/api/{id}", baseHandler.GetByIDAPI).Methods("GET")
	r.HandleFunc("/api/create", baseHandler.CreateToDoAPI).Methods("POST")
	r.HandleFunc("/api/{id}/delete", baseHandler.DeleteToDoAPI).Methods("POST")
	r.HandleFunc("/api/{id}/update", baseHandler.UpdateToDoAPI).Methods("POST")
	http.ListenAndServe(":8080", r)
}
