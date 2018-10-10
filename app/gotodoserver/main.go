package main

import (
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
	"github.com/subosito/gotenv"

	"github.com/saifulwebid/gotodo"
	"github.com/saifulwebid/gotodo/database"
)

func init() {
	gotenv.Load()
}

func main() {
	db, err := database.NewRepository()
	if err != nil {
		log.Fatal(err)
	}

	sv := &server{gotodo.NewService(db)}

	router := httprouter.New()

	router.GET("/", sv.getAll)
	router.GET("/:id", sv.get)
	router.POST("/", sv.add)
	router.PATCH("/:id", sv.edit)
	router.PUT("/:id/done", sv.markAsDone)
	router.DELETE("/:id", sv.delete)
	router.DELETE("/", sv.deleteFinished)

	log.Fatal(http.ListenAndServe(":"+os.Getenv("GOTODO_API_PORT"), router))
}
