package pgdb

import (
	"restApi/myErrors"

	"database/sql"
	"fmt"

	//
	_ "github.com/lib/pq"
)

func OpenDB(connStr string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		myErrors.Check(err)
	}

	if err = db.Ping(); err != nil {
		myErrors.Check(err)
	}

	return db, nil
}

type Todo struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
}

type TodoModel struct {
	DB *sql.DB
}

func (t *TodoModel) Insert(description string) error {
	statement := `INSERT INTO todo (todo_description) VALUES ($1);`

	_, err := t.DB.Exec(statement, description)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (t *TodoModel) Delete(id int) (int, error) {
	statement := `DELETE FROM todo WHERE todo_id = $1;`

	result, err := t.DB.Exec(statement, id)
	if err != nil {
		myErrors.Check(err)
		return -1, err
	}

	row, err := result.RowsAffected()

	if err != nil {
		myErrors.Check(err)
		return -1, err
	}

	return int(row), nil
}

func (t *TodoModel) GetTodos() ([]Todo, error) {
	var todoList []Todo
	var todoItem Todo

	var id int
	var description string
	var returnErr error = nil

	rows, err := t.DB.Query(`SELECT * FROM todo;`)

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&id, &description)
		todoItem = Todo{id, description}

		if err != nil {
			myErrors.Check(err)
			returnErr = err
		}
		todoList = append(todoList, todoItem)
	}

	return todoList, returnErr
}
