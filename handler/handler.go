package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"

	"github.com/saifulwebid/gotodo"
)

func respondWithErrorInJSON(w http.ResponseWriter, code int, err error) {
	respondInJSON(w, code, map[string]string{"error": err.Error()})
}

func respondInJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	response, _ := json.Marshal(payload)
	w.Write(response)
}

type Server struct {
	Service gotodo.Service
	Router  *httprouter.Router
}

func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	s.Router.GET("/", s.getAll)
	s.Router.GET("/:id", s.get)
	s.Router.POST("/", s.add)
	s.Router.PATCH("/:id", s.edit)
	s.Router.PUT("/:id/done", s.markAsDone)
	s.Router.DELETE("/:id", s.delete)
	s.Router.DELETE("/", s.deleteFinished)

	s.Router.ServeHTTP(w, req)
}

// Get is a handler for GET "/:id" route. It will return a Todo with specified
// ID if exists. It will return an error otherwise.
func (s *Server) get(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		respondWithErrorInJSON(w, http.StatusBadRequest, errors.New("cannot parse id"))
		return
	}

	todo, err := s.Service.Get(id)
	if err != nil {
		respondWithErrorInJSON(w, http.StatusNotFound, errors.New("Todo not found"))
		return
	}

	respondInJSON(w, http.StatusOK, todo)
}

// GetAll is a handler for GET "/" route. It will return an array of Todos.
//
// A query string called "done" can also exist on the request. This query string
// should be either "true" or "false". "true" means that user wants to get all
// finished Todos; "false" otherwise.
func (s *Server) getAll(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var todos []*gotodo.Todo

	done, ok := r.URL.Query()["done"]
	if ok && len(done[0]) > 0 {
		if done[0] == "true" {
			todos = s.Service.GetFinished()
		} else {
			todos = s.Service.GetPending()
		}
	} else {
		todos = s.Service.GetAll()
	}

	respondInJSON(w, http.StatusOK, todos)
}

// Add is a handler for POST "/" route. It receives a JSON which corresponds to
// a Todo structure, adds the Todo using gotodo.Service, and returns the JSON
// from the gotodo.Service. It returns an error if such error occurs.
//
// Add will only respect .title and .description from the JSON request body.
func (s *Server) add(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	defer r.Body.Close()

	todo := &gotodo.Todo{}
	if err := json.NewDecoder(r.Body).Decode(todo); err != nil {
		respondWithErrorInJSON(w, http.StatusBadRequest, errors.New("Invalid request payload"))
		return
	}

	todo, err := s.Service.Add(todo.Title, todo.Description)
	if err != nil {
		respondWithErrorInJSON(w, http.StatusInternalServerError, err)
		return
	}

	respondInJSON(w, http.StatusCreated, todo)
}

// Edit is a handler for PATCH "/:id" route to edit a Todo. It receives a JSON
// which corresponds a Todo structure, applies the new values to the old values
// using gotodo.Service, and returns back the Todo from the service. It returns
// an error if such error occurs.
//
// It will only respect .title and .description attribute, as .done is modified
// only through MarkAsDone, as the gotodo package requests.
func (s *Server) edit(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		respondWithErrorInJSON(w, http.StatusBadRequest, errors.New("cannot parse id"))
		return
	}

	todo, err := s.Service.Get(id)
	if err != nil {
		respondWithErrorInJSON(w, http.StatusNotFound, errors.New("Todo not found"))
		return
	}

	type InputJSON struct {
		Title       string  `json:"title"`
		Description *string `json:"description,omitempty"`
	}

	todoEdit := &InputJSON{}
	if err := json.NewDecoder(r.Body).Decode(todoEdit); err != nil {
		respondWithErrorInJSON(w, http.StatusBadRequest, errors.New("Invalid request payload"))
		return
	}

	todo.Title = todoEdit.Title
	if todoEdit.Description != nil {
		todo.Description = *todoEdit.Description
	}

	err = s.Service.Edit(todo)
	if err != nil {
		respondWithErrorInJSON(w, http.StatusInternalServerError, err)
		return
	}

	respondInJSON(w, http.StatusOK, todo)
}

// MarkAsDone is a handler for PUT "/:id/done" route to mark a Todo as done.
// It receives an empty request and returns the marked Todo from the service,
// or an error if such error exists.
func (s *Server) markAsDone(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		respondWithErrorInJSON(w, http.StatusBadRequest, errors.New("cannot parse id"))
		return
	}

	todo, err := s.Service.Get(id)
	if err != nil {
		respondWithErrorInJSON(w, http.StatusNotFound, errors.New("Todo not found"))
		return
	}

	err = s.Service.MarkAsDone(todo)
	if err != nil {
		respondWithErrorInJSON(w, http.StatusInternalServerError, err)
		return
	}

	respondInJSON(w, http.StatusOK, todo)
}

// Delete is a handler for DELETE "/:id" route to delete a Todo. It will return
// an error if such error exists.
func (s *Server) delete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		respondWithErrorInJSON(w, http.StatusBadRequest, errors.New("cannot parse id"))
		return
	}

	todo, err := s.Service.Get(id)
	if err != nil {
		respondWithErrorInJSON(w, http.StatusNotFound, errors.New("Todo not found"))
		return
	}

	err = s.Service.Delete(todo)
	if err != nil {
		respondWithErrorInJSON(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(200)
}

// DeleteFinished is a handler for DELETE "/" route to delete all finished
// Todos. It must receive a "done" query string with "true" value; otherwise,
// it will return an error message.
func (s *Server) deleteFinished(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	done, ok := r.URL.Query()["done"]
	if !ok || done[0] != "true" {
		respondWithErrorInJSON(w, http.StatusBadRequest, errors.New("?done=true should be set"))
		return
	}

	s.Service.DeleteFinished()

	w.WriteHeader(200)
}
