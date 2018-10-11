package cli

import (
	"fmt"

	"github.com/saifulwebid/gotodo"
)

func todoToString(todo *gotodo.Todo) string {
	template := `ID: %d
Title: %s
Description: %s
Done: %s
`

	done := "Pending"
	if todo.Done {
		done = "Finished"
	}

	return fmt.Sprintf(template, todo.ID, todo.Title, todo.Description, done)
}

func todosToString(todos []*gotodo.Todo) string {
	ret := ""

	for i, todo := range todos {
		if i > 0 {
			ret += "----------------------------\n"
		}
		ret += todoToString(todo)
	}

	return ret
}
