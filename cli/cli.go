package cli

import (
	"github.com/saifulwebid/gotodo"
	"gopkg.in/urfave/cli.v1"
)

// Application is a wrapper to urfave/cli package. It also contains an instance
// to gotodo.Service to be used by all CLI commands.
type Application struct {
	Service gotodo.Service
}

// Run will set up an urfave/cli.App instance and run it.
func (a *Application) Run(arguments []string) error {
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
			Action: a.getAll,
		},
		{
			Name:   "get",
			Usage:  "get a todo from the database",
			Action: a.get,
		},
		{
			Name:   "create",
			Usage:  "create a todo in the database",
			Flags:  todoFlags,
			Action: a.create,
		},
		{
			Name:   "edit",
			Usage:  "edit a todo in the database",
			Flags:  todoFlags,
			Action: a.edit,
		},
		{
			Name:   "done",
			Usage:  "mark a todo as done",
			Action: a.markAsDone,
		},
		{
			Name:   "delete",
			Usage:  "delete a todo from the database",
			Action: a.delete,
		},
		{
			Name:   "delete-finished",
			Usage:  "delete all finished todos from the database",
			Action: a.deleteFinished,
		},
	}

	return app.Run(arguments)
}
