package handlers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	tm "todo/models"

	"github.com/gobuffalo/envy"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// function GetByIDAPI provides ToDo by ID
func (bh *BaseHandler) GetByIDAPI(w http.ResponseWriter, r *http.Request) {
	appApiKey, err := envy.MustGet("API_KEY")
	if err != nil {
		log.Fatalln(err)
	}

	headerApiKey := r.Header.Get("Authorization")

	if headerApiKey != appApiKey {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.WriteHeader(401)

		_ = json.NewEncoder(w).Encode("Unauthorized")

		return
	}

	params := mux.Vars(r)
	toDoIDStr := params["id"]

	toDoID, err := strconv.ParseUint(toDoIDStr, 10, 32)
	if err != nil {
		log.Fatalln(err)
	}

	toDo, err := tm.GetToDoByID(bh.App.MyDB, uint(toDoID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.Header().Set("X-Content-Type-Options", "nosniff")
			w.WriteHeader(400)

			_ = json.NewEncoder(w).Encode("No ToDo found under provided ID")

			return
		} else {
			log.Fatalln(err)
		}
	}

	respBytes, err := json.Marshal(toDo)
	if err != nil {
		log.Fatalln("Could not marshal a ToDo")
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(200)
	w.Write(respBytes)
}

// function GetAllAPI provides all ToDos
func (bh *BaseHandler) GetAllAPI(w http.ResponseWriter, r *http.Request) {
	appApiKey, err := envy.MustGet("API_KEY")
	if err != nil {
		log.Fatalln(err)
	}

	headerApiKey := r.Header.Get("Authorization")

	if headerApiKey != appApiKey {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.WriteHeader(401)

		_ = json.NewEncoder(w).Encode("Unauthorized")

		return
	}

	toDos, err := tm.GetAllToDos(bh.App.MyDB)
	if err != nil {
		log.Fatalln(err)
	}

	respBytes, err := json.Marshal(toDos)
	if err != nil {
		log.Fatalln("Could not marshal a ToDo")
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(200)
	w.Write(respBytes)
}

// function DeleteToDoAPI deletes particular ToDo from REST request
func (bh *BaseHandler) DeleteToDoAPI(w http.ResponseWriter, r *http.Request) {
	appApiKey, err := envy.MustGet("API_KEY")
	if err != nil {
		log.Fatalln(err)
	}

	headerApiKey := r.Header.Get("Authorization")

	if headerApiKey != appApiKey {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.WriteHeader(401)

		_ = json.NewEncoder(w).Encode("Unauthorized")

		return
	}

	params := mux.Vars(r)
	toDoIDStr := params["id"]

	toDoID, err := strconv.ParseUint(toDoIDStr, 10, 32)
	if err != nil {
		log.Fatalln(err)
	}

	toDo, err := tm.GetToDoByID(bh.App.MyDB, uint(toDoID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.Header().Set("X-Content-Type-Options", "nosniff")
			w.WriteHeader(400)

			_ = json.NewEncoder(w).Encode("No ToDo found under provided ID")

			return
		} else {
			log.Fatalln(err)
		}
	}

	err = toDo.Delete(bh.App.MyDB)
	if err != nil {
		log.Fatalln(err)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(200)
}

// function UpdateToDoAPI updates particular ToDo from REST request
func (bh *BaseHandler) UpdateToDoAPI(w http.ResponseWriter, r *http.Request) {
	appApiKey, err := envy.MustGet("API_KEY")
	if err != nil {
		log.Fatalln(err)
	}

	headerApiKey := r.Header.Get("Authorization")

	if headerApiKey != appApiKey {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.WriteHeader(401)

		_ = json.NewEncoder(w).Encode("Unauthorized")

		return
	}

	params := mux.Vars(r)
	toDoIDStr := params["id"]

	toDoID, err := strconv.ParseUint(toDoIDStr, 10, 32)
	if err != nil {
		log.Fatalln(err)
	}

	toDo, err := tm.GetToDoByID(bh.App.MyDB, uint(toDoID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.Header().Set("X-Content-Type-Options", "nosniff")
			w.WriteHeader(400)

			_ = json.NewEncoder(w).Encode("No ToDo found under provided ID")

			return
		} else {
			log.Fatalln(err)
		}
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.WriteHeader(400)

		_ = json.NewEncoder(w).Encode("Could not read body from request")

		return
	}

	err = json.Unmarshal(body, &toDo)
	if err != nil {
		log.Fatalln("Could not unmarshall JSON into ToDo")
	}

	if toDo.Title == "" ||
		toDo.Description == "" {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.WriteHeader(400)

		_ = json.NewEncoder(w).Encode("Make sure to provide title and description")

		return
	}

	err = toDo.Update(bh.App.MyDB)
	if err != nil {
		log.Fatalln(err)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(200)
}

// function CreateToDoAPI creates ToDo from REST request
func (bh *BaseHandler) CreateToDoAPI(w http.ResponseWriter, r *http.Request) {
	appApiKey, err := envy.MustGet("API_KEY")
	if err != nil {
		log.Fatalln(err)
	}

	headerApiKey := r.Header.Get("Authorization")

	if headerApiKey != appApiKey {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.WriteHeader(401)

		_ = json.NewEncoder(w).Encode("Unauthorized")

		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.WriteHeader(400)

		_ = json.NewEncoder(w).Encode("Could not read body from request")

		return
	}

	toDo := tm.ToDo{}

	err = json.Unmarshal(body, &toDo)
	if err != nil {
		log.Fatalln("Could not unmarshall JSON into ToDo")
	}

	if toDo.Title == "" ||
		toDo.Description == "" {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.WriteHeader(400)

		_ = json.NewEncoder(w).Encode("Make sure to provide title and description")

		return
	}

	err = toDo.Create(bh.App.MyDB)
	if err != nil {
		log.Fatalln(err)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(200)
}
