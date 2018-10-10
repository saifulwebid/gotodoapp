package main

import (
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"

	"github.com/saifulwebid/gotodo"
)

type server struct {
	service gotodo.Service
}

// Get is a handler for GET "/:id" route. It will return a Todo with specified
// ID if exists. It will return an error otherwise.
func (s *server) get(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		// TODO: return error JSON (invalid id)
	}

	todo, err := s.service.Get(id)
	if err != nil {
		// TODO: return error JSON (id not found)
	}

	// TODO: return todo as JSON
	todo = todo
}

// GetAll is a handler for GET "/" route. It will return an array of Todos.
//
// A query string called "done" can also exist on the request. This query string
// should be either "true" or "false". "true" means that user wants to get all
// finished Todos; "false" otherwise.
func (s *server) getAll(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var res []*gotodo.Todo

	done, ok := r.URL.Query()["done"]
	if ok && len(done[0]) > 0 {
		if done[0] == "true" {
			res = s.service.GetFinished()
		} else {
			res = s.service.GetPending()
		}
	} else {
		res = s.service.GetAll()
	}

	// TODO: return res as JSON
	res = res
}

// Add is a handler for POST "/" route. It receives a JSON which corresponds to
// a Todo structure, adds the Todo using gotodo.Service, and returns the JSON
// from the gotodo.Service. It returns an error if such error occurs.
//
// Add will only respect .title and .description from the JSON request body.
func (s *server) add(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// TODO: read JSON from request
	var todo *gotodo.Todo

	todo, err := s.service.Add(todo.Title, todo.Description)
	if err != nil {
		// TODO: return error JSON
	}

	// TODO: return todo as JSON
	todo = todo
}

// Edit is a handler for PATCH "/:id" route to edit a Todo. It receives a JSON
// which corresponds a Todo structure, applies the new values to the old values
// using gotodo.Service, and returns back the Todo from the service. It returns
// an error if such error occurs.
//
// It will only respect .title and .description attribute, as .done is modified
// only through MarkAsDone, as the gotodo package requests.
func (s *server) edit(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		// TODO: return error JSON (invalid id)
	}

	todo, err := s.service.Get(id)
	if err != nil {
		// TODO: return error JSON (id not found)
	}

	// TODO: read JSON from request
	var todoEdit *gotodo.Todo

	todo.Title = todoEdit.Title
	todo.Description = todoEdit.Description

	err = s.service.Edit(todo)
	if err != nil {
		// TODO: return error JSON
	}

	// TODO: return todo as JSON
	todo = todo
}

// MarkAsDone is a handler for PUT "/:id/done" route to mark a Todo as done.
// It receives an empty request and returns the marked Todo from the service,
// or an error if such error exists.
func (s *server) markAsDone(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		// TODO: return error JSON (invalid id)
	}

	todo, err := s.service.Get(id)
	if err != nil {
		// TODO: return error JSON (id not found)
	}

	err = s.service.MarkAsDone(todo)
	if err != nil {
		// TODO: return error JSON
	}

	// TODO: return todo as JSON
	todo = todo
}

// Delete is a handler for DELETE "/:id" route to delete a Todo. It will return
// an error if such error exists.
func (s *server) delete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		// TODO: return error JSON (invalid id)
	}

	todo, err := s.service.Get(id)
	if err != nil {
		// TODO: return error JSON (id not found)
	}

	err = s.service.Delete(todo)
	if err != nil {
		// TODO: return error JSON
	}

	// TODO: determine proper response
}

// DeleteFinished is a handler for DELETE "/" route to delete all finished
// Todos. It must receive a "done" query string with "true" value; otherwise,
// it will return an error message.
func (s *server) deleteFinished(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	done, ok := r.URL.Query()["done"]
	if !ok || done[0] != "true" {
		// TODO: return error JSON
	}

	s.service.DeleteFinished()
}
