package main

import (
	"log"
	"os"

	"github.com/saifulwebid/gotodo"
	"github.com/saifulwebid/gotodo/database"
	"github.com/subosito/gotenv"

	"github.com/saifulwebid/gotodoapp/cli"
)

func init() {
	gotenv.Load()
}

func main() {
	db, err := database.NewRepository()
	if err != nil {
		log.Fatal(err)
	}

	service := gotodo.NewService(db)

	app := &cli.Application{
		Service: service,
	}

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
