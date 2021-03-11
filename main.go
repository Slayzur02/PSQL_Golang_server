package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"restApi/myErrors"
	pgdb "restApi/pgDB"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

var todoController *pgdb.TodoModel

// HomePage - handler for /
type HomePage struct{}

func (h *HomePage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "this is cool")
}

// GetAllTodos - handler for /todos
type GetAllTodos struct{}

func (g *GetAllTodos) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	todoList, err := todoController.GetTodos()
	if err != nil {
		myErrors.Check(err)
	} else {
		json.NewEncoder(w).Encode(todoList)
	}
}

// AddSingleTodo - handler for /todo/id
type AddSingleTodo struct{}

func (a *AddSingleTodo) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	desc := vars["string"]

	err := todoController.Insert(desc)
	if err != nil {
		myErrors.Check(err)
	}

	http.Redirect(w, r, "/getall", http.StatusSeeOther)
}

// DeleteSingleTodo ndler for /delete
type DeleteSingleTodo struct{}

func (d *DeleteSingleTodo) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	stringID := vars["id"]

	intID, err := strconv.Atoi(stringID)
	if err != nil {
		myErrors.Check(err)
	}

	_, err = todoController.Delete(intID)
	if err != nil {
		myErrors.Check(err)
	}
}

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		myErrors.Check(err)
	}

	username := os.Getenv("PSQL_USERNAME")
	password := os.Getenv("PSQL_PASSWORD")

	websiteRouter := mux.NewRouter().StrictSlash(true)
	websiteRouter.Handle("/", new(HomePage))
	websiteRouter.Handle("/getall", new(GetAllTodos))
	websiteRouter.Handle("/insert/{string}", new(AddSingleTodo))
	websiteRouter.Handle("/delete/{id}", new(DeleteSingleTodo))

	connectionString := "user=" + username + " password='" + password + "' dbname=todo sslmode=disable"
	db, err := pgdb.OpenDB(connectionString)

	if err != nil {
		myErrors.Check(err)
	}

	todoController = &pgdb.TodoModel{DB: db}

	http.ListenAndServe(":8080", websiteRouter)

}
