package handlers

import (
	"log"
	"net/http"
	"strconv"
	"text/template"
	"time"
	"todo/application"
	tm "todo/models"

	"github.com/gorilla/mux"
)

type BaseHandler struct {
	App *application.Application
}

type ValidationError struct {
	Message string
}

type ToDoView struct {
	ToDos            []tm.ToDo
	ValidationErrors []ValidationError
}

func NewAppHandler(app *application.Application) *BaseHandler {
	return &BaseHandler{
		App: app,
	}
}

// MainPage renders main page of the ToDo application
func (bh *BaseHandler) MainPage(w http.ResponseWriter, r *http.Request) {
	tpl, err := template.ParseFiles("templates/index.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		log.Fatalln(err)
	}

	toDos, err := tm.GetAllToDos(bh.App.MyDB)
	if err != nil {
		log.Fatalln(err)
	}

	toDoView := ToDoView{
		ToDos: toDos,
	}

	tpl.Execute(w, toDoView)
}

// func CreateToDo creates a new ToDo from form
func (bh *BaseHandler) CreateToDo(w http.ResponseWriter, r *http.Request) {
	var (
		validationErrors []ValidationError
		dateOfExpiry     *time.Time
		percentage       = 0
	)

	tpl, err := template.ParseFiles("templates/index.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		log.Fatalln(err)
	}

	toDos, err := tm.GetAllToDos(bh.App.MyDB)
	if err != nil {
		log.Fatalln(err)
	}

	err = r.ParseForm()
	if r.FormValue("title") == "" {
		validationErrors = append(validationErrors, ValidationError{"Please fill in title of ToDo"})
	}

	if r.FormValue("description") == "" {
		validationErrors = append(validationErrors, ValidationError{"Please fill in description of ToDo"})
	}

	if r.FormValue("dateOfExpiry") != "" {
		dateOfExpiryTime, err := time.Parse("2006-01-02", r.FormValue("dateOfExpiry"))
		if err != nil {
			log.Fatalln(err)
		}

		dateOfExpiry = &dateOfExpiryTime
	}

	if r.FormValue("percentage") != "" {
		percentage, err = strconv.Atoi(r.FormValue("percentage"))
		if err != nil {
			log.Fatalln(err)
		}
	}

	if len(validationErrors) == 0 {
		newToDo := &tm.ToDo{
			Title:               r.FormValue("title"),
			Description:         r.FormValue("description"),
			DateAndTimeOfExpiry: dateOfExpiry,
			CompletePercentage:  percentage,
		}

		err = newToDo.Create(bh.App.MyDB)
		if err != nil {
			log.Fatalln(err)
		}

		http.Redirect(w, r, "/", http.StatusFound)
	}

	toDoView := ToDoView{
		ToDos:            toDos,
		ValidationErrors: validationErrors,
	}

	tpl.Execute(w, toDoView)
}

// function UpdateToDo updates particular ToDo
func (bh *BaseHandler) UpdateToDo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	toDoIDStr := params["id"]

	err := r.ParseForm()
	toDoID, err := strconv.ParseUint(toDoIDStr, 10, 32)
	if err != nil {
		log.Fatalln(err)
	}

	toDo, err := tm.GetToDoByID(bh.App.MyDB, uint(toDoID))
	if err != nil {
		log.Fatalln(err)
	}

	toDo.CompletePercentage, err = strconv.Atoi(r.FormValue("percentage"))
	if err != nil {
		log.Fatalln(err)
	}

	err = toDo.Update(bh.App.MyDB)
	if err != nil {
		log.Fatalln(err)
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

// function DeleteToDo deletes particular ToDo
func (bh *BaseHandler) DeleteToDo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	toDoIDStr := params["id"]

	toDoID, err := strconv.ParseUint(toDoIDStr, 10, 32)
	if err != nil {
		log.Fatalln(err)
	}

	toDo, err := tm.GetToDoByID(bh.App.MyDB, uint(toDoID))
	if err != nil {
		log.Fatalln(err)
	}

	err = toDo.Delete(bh.App.MyDB)
	if err != nil {
		log.Fatalln(err)
	}

	http.Redirect(w, r, "/", http.StatusFound)
}
