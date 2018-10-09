package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
	"github.com/subosito/gotenv"
)

func main() {
	gotenv.Load()

	router := httprouter.New()

	router.GET("/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		fmt.Fprintln(w, "Hello world!")
	})

	log.Fatal(http.ListenAndServe(":"+os.Getenv("GOTODO_API_PORT"), router))
}
