package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

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

	if c.IsSet("done") {
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

func get(c *cli.Context) error {
	idStr := c.Args().Get(0)
	if idStr == "" {
		log.Fatal("ID argument must not be blank")
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Fatal(err)
	}

	todo, err := service.Get(id)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(todoToString(todo))

	return nil
}

func create(c *cli.Context) error {
	todo, err := service.Add(c.String("title"), c.String("description"))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Created todo:")
	fmt.Println(todoToString(todo))

	return nil
}

func edit(c *cli.Context) error {
	if c.NumFlags() == 0 {
		log.Fatal("No --title or --description set; exiting")
	}

	idStr := c.Args().Get(0)
	if idStr == "" {
		log.Fatal("ID argument must not be blank")
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Fatal(err)
	}

	todo, err := service.Get(id)
	if err != nil {
		log.Fatal(err)
	}

	if c.IsSet("title") {
		todo.Title = c.String("title")
	}

	if c.IsSet("description") {
		todo.Title = c.String("description")
	}

	err = service.Edit(todo)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Edited todo:")
	fmt.Println(todoToString(todo))

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
	todoFlags := []cli.Flag{
		cli.StringFlag{
			Name: "title, t",
		},
		cli.StringFlag{
			Name: "description, desc, d",
		},
	}

	app := cli.NewApp()

	app.Name = "gotodocli"
	app.Usage = "manage your todos"
	app.Commands = []cli.Command{
		{
			Name:   "getall",
			Usage:  "get all todos from the database",
			Flags:  doneFlags,
			Action: getAll,
		},
		{
			Name:   "get",
			Usage:  "get a todo from the database",
			Action: get,
		},
		{
			Name:   "create",
			Usage:  "create a todo in the database",
			Flags:  todoFlags,
			Action: create,
		},
		{
			Name:   "edit",
			Usage:  "edit a todo in the database",
			Flags:  todoFlags,
			Action: edit,
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
