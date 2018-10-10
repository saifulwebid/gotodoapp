package main

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/urfave/cli.v1"
)

func main() {
	app := cli.NewApp()

	app.Name = "gotodocli"
	app.Usage = "manage your todos"

	app.Action = func(c *cli.Context) error {
		fmt.Println("Hello world!")

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
