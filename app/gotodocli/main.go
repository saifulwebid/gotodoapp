package main

import (
	"fmt"
	"log"
	"os"

	"github.com/subosito/gotenv"
	"gopkg.in/urfave/cli.v1"

	"github.com/saifulwebid/gotodo"
	"github.com/saifulwebid/gotodo/database"
)

func init() {
	gotenv.Load()
}

var db gotodo.Repository
var service gotodo.Service

func getAll(c *cli.Context) error {
	var todos []*gotodo.Todo

	if c.NumFlags() > 0 {
		if c.Bool("done") {
			todos = service.GetFinished()
		} else {
			todos = service.GetPending()
		}
	} else {
		todos = service.GetAll()
	}

	fmt.Println(todosToString(todos))

	return nil
}

func main() {
	db, err := database.NewRepository()
	if err != nil {
		log.Fatal(err)
	}

	service = gotodo.NewService(db)

	doneFlags := []cli.Flag{
		cli.BoolFlag{
			Name: "done, d",
		},
	}

	app := cli.NewApp()

	app.Name = "gotodocli"
	app.Usage = "manage your todos"
	app.Commands = []cli.Command{
		{
			Name:   "getall",
			Usage:  "get all todos on the database",
			Flags:  doneFlags,
			Action: getAll,
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

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
