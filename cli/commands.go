package cli

import (
	"fmt"
	"log"
	"strconv"

	"gopkg.in/urfave/cli.v1"

	"github.com/saifulwebid/gotodo"
)

func parseIDFromCli(c *cli.Context) int {
	idStr := c.Args().Get(0)
	if idStr == "" {
		log.Fatal("ID argument must not be blank")
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Fatal(err)
	}

	return id
}

func (a *Application) getAll(c *cli.Context) error {
	var todos []*gotodo.Todo

	if c.IsSet("done") {
		if c.Bool("done") {
			todos = a.Service.GetFinished()
		} else {
			todos = a.Service.GetPending()
		}
	} else {
		todos = a.Service.GetAll()
	}

	fmt.Println(todosToString(todos))

	return nil
}

func (a *Application) get(c *cli.Context) error {
	id := parseIDFromCli(c)

	todo, err := a.Service.Get(id)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(todoToString(todo))

	return nil
}

func (a *Application) create(c *cli.Context) error {
	todo, err := a.Service.Add(c.String("title"), c.String("description"))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Created todo:")
	fmt.Println(todoToString(todo))

	return nil
}

func (a *Application) edit(c *cli.Context) error {
	if c.NumFlags() == 0 {
		log.Fatal("No --title or --description set; exiting")
	}

	id := parseIDFromCli(c)

	todo, err := a.Service.Get(id)
	if err != nil {
		log.Fatal(err)
	}

	if c.IsSet("title") {
		todo.Title = c.String("title")
	}

	if c.IsSet("description") {
		todo.Title = c.String("description")
	}

	err = a.Service.Edit(todo)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Edited todo:")
	fmt.Println(todoToString(todo))

	return nil
}

func (a *Application) markAsDone(c *cli.Context) error {
	id := parseIDFromCli(c)

	todo, err := a.Service.Get(id)
	if err != nil {
		log.Fatal(err)
	}

	err = a.Service.MarkAsDone(todo)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Todo marked as done:")
	fmt.Println(todoToString(todo))

	return nil
}

func (a *Application) delete(c *cli.Context) error {
	id := parseIDFromCli(c)

	todo, err := a.Service.Get(id)
	if err != nil {
		log.Fatal(err)
	}

	err = a.Service.Delete(todo)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Todo deleted:")
	fmt.Println(todoToString(todo))

	return nil
}

func (a *Application) deleteFinished(c *cli.Context) error {
	a.Service.DeleteFinished()

	fmt.Println("all finished todo are deleted")

	return nil
}
