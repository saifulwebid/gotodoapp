package main

import (
	"log"
	"net/http"
	"os"

	"github.com/subosito/gotenv"

	"github.com/saifulwebid/gotodo"
	"github.com/saifulwebid/gotodo/database"
	"github.com/saifulwebid/gotodoapp/handler"
)

func init() {
	gotenv.Load()
}

func main() {
	db, err := database.NewRepository()
	if err != nil {
		log.Fatal(err)
	}

	sv := handler.NewServer(gotodo.NewService(db))

	log.Fatal(http.ListenAndServe(":"+os.Getenv("GOTODO_API_PORT"), sv))
}
