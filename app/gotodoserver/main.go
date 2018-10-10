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

	sv := &server{gotodo.NewService(db), httprouter.New()}

	sv.router.GET("/", sv.getAll)
	sv.router.GET("/:id", sv.get)
	sv.router.POST("/", sv.add)
	sv.router.PATCH("/:id", sv.edit)
	sv.router.PUT("/:id/done", sv.markAsDone)
	sv.router.DELETE("/:id", sv.delete)
	sv.router.DELETE("/", sv.deleteFinished)

	log.Fatal(http.ListenAndServe(":"+os.Getenv("GOTODO_API_PORT"), sv))
}
