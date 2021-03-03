package main

import (
	"net/http"
	"encoding/json"
    "github.com/gorilla/mux"
	"fmt"
	"strconv"
	"restApi/myErrors"
	"restApi/pgDB"
	"github.com/joho/godotenv"
	"os"
)




var todoController *pgDB.TodoModel

type HomePage struct{}
func (h *HomePage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "this is cool")
}

type GetAllTodos struct{}
func (g *GetAllTodos) ServeHTTP(w http.ResponseWriter, r *http.Request){
	todoList, err := todoController.GetTodos()
	if err != nil {
		myErrors.Check(err)
	} else{
		json.NewEncoder(w).Encode(todoList)
	}
}

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

type DeleteSingleTodo struct{}
func (d *DeleteSingleTodo) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	stringId := vars["id"]

	intId, err := strconv.Atoi(stringId)
	if (err != nil) {
		myErrors.Check(err)
	}

	_, err = todoController.Delete(intId)
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
	db, err := pgDB.OpenDB(connectionString)

	if err!= nil {
		myErrors.Check(err)
	}

	todoController = &pgDB.TodoModel{db}

	http.ListenAndServe(":8080", websiteRouter)

}
