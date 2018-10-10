package handler_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/saifulwebid/gotodo"
	"github.com/saifulwebid/gotodoapp/handler"
)

func execute(h *handler.Server, req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	return rr
}

func TestGet(t *testing.T) {
	svc := &mockService{
		GetFn: func(id int) (*gotodo.Todo, error) {
			return nil, nil
		},
	}
	h := handler.NewServer(svc)

	t.Run("bad request", func(t *testing.T) {
		svc.GetInvoked = 0

		req := httptest.NewRequest("GET", "/1a2b", nil)
		rr := execute(h, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, 0, svc.GetInvoked)
	})

	t.Run("not found", func(t *testing.T) {
		svc.GetInvoked = 0
		svc.GetFn = func(id int) (*gotodo.Todo, error) {
			return nil, errors.New("not found")
		}

		req := httptest.NewRequest("GET", "/1", nil)
		rr := execute(h, req)

		assert.Equal(t, http.StatusNotFound, rr.Code)
		assert.Equal(t, 1, svc.GetInvoked)
	})

	t.Run("found", func(t *testing.T) {
		svc.GetInvoked = 0
		svc.GetFn = func(id int) (*gotodo.Todo, error) {
			return &gotodo.Todo{}, nil
		}

		req := httptest.NewRequest("GET", "/1", nil)
		rr := execute(h, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, 1, svc.GetInvoked)
	})
}

func TestGetTodos(t *testing.T) {
	svc := &mockService{
		GetAllFn: func() []*gotodo.Todo {
			return []*gotodo.Todo{}
		},
		GetPendingFn: func() []*gotodo.Todo {
			return []*gotodo.Todo{}
		},
		GetFinishedFn: func() []*gotodo.Todo {
			return []*gotodo.Todo{}
		},
	}
	h := handler.NewServer(svc)

	t.Run("get all", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		rr := execute(h, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, 1, svc.GetAllInvoked)
	})

	t.Run("get pending", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/?done=false", nil)
		rr := execute(h, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, 1, svc.GetPendingInvoked)
	})

	t.Run("get finished", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/?done=true", nil)
		rr := execute(h, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, 1, svc.GetFinishedInvoked)
	})
}

func TestAdd(t *testing.T) {
	svc := &mockService{
		AddFn: func(title string, description string) (*gotodo.Todo, error) {
			if title == "" {
				return nil, errors.New("invalid todo")
			}

			return &gotodo.Todo{1, title, description, false}, nil
		},
	}
	h := handler.NewServer(svc)

	t.Run("invalid todo", func(t *testing.T) {
		svc.AddInvoked = 0
		payload := []byte(`{"title": "", "description": "desc"}`)

		req := httptest.NewRequest("POST", "/", bytes.NewBuffer(payload))
		rr := execute(h, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Equal(t, 1, svc.AddInvoked)
	})

	t.Run("valid todo", func(t *testing.T) {
		svc.AddInvoked = 0
		payload := []byte(`{"title": "Title", "description": "desc"}`)

		req := httptest.NewRequest("POST", "/", bytes.NewBuffer(payload))
		rr := execute(h, req)

		assert.Equal(t, http.StatusCreated, rr.Code)
		assert.Equal(t, 1, svc.AddInvoked)
	})
}

func TestEdit(t *testing.T) {
	svc := &mockService{
		GetFn: func(id int) (*gotodo.Todo, error) {
			if id != 1 {
				return nil, errors.New("not found")
			}

			return &gotodo.Todo{1, "title", "description", false}, nil
		},
		EditFn: func(todo *gotodo.Todo) error {
			if todo.Title == "" {
				return errors.New("invalid todo")
			}

			return nil
		},
	}
	h := handler.NewServer(svc)

	t.Run("invalid id", func(t *testing.T) {
		svc.GetInvoked = 0
		svc.EditInvoked = 0

		req := httptest.NewRequest("PATCH", "/abc", nil)
		rr := execute(h, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, 0, svc.GetInvoked)
		assert.Equal(t, 0, svc.EditInvoked)
	})

	t.Run("non-existent todo", func(t *testing.T) {
		svc.GetInvoked = 0
		svc.EditInvoked = 0
		payload := []byte(`{"title": "", "description": "desc"}`)

		req := httptest.NewRequest("PATCH", "/100", bytes.NewBuffer(payload))
		rr := execute(h, req)

		assert.Equal(t, http.StatusNotFound, rr.Code)
		assert.Equal(t, 1, svc.GetInvoked)
		assert.Equal(t, 0, svc.EditInvoked)
	})

	t.Run("invalid todo", func(t *testing.T) {
		svc.GetInvoked = 0
		svc.EditInvoked = 0
		payload := []byte(`{"title": "", "description": "desc"}`)

		req := httptest.NewRequest("PATCH", "/1", bytes.NewBuffer(payload))
		rr := execute(h, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Equal(t, 1, svc.GetInvoked)
		assert.Equal(t, 1, svc.EditInvoked)
	})

	t.Run("valid todo", func(t *testing.T) {
		svc.GetInvoked = 0
		svc.EditInvoked = 0
		payload := []byte(`{"title": "hai", "description": "desc"}`)

		req := httptest.NewRequest("PATCH", "/1", bytes.NewBuffer(payload))
		rr := execute(h, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, 1, svc.GetInvoked)
		assert.Equal(t, 1, svc.EditInvoked)
	})
}

func TestMarkAsDone(t *testing.T) {
	svc := &mockService{
		GetFn: func(id int) (*gotodo.Todo, error) {
			if id != 1 {
				return nil, errors.New("not found")
			}

			return &gotodo.Todo{1, "title", "description", false}, nil
		},
		MarkAsDoneFn: func(todo *gotodo.Todo) error {
			if todo.Title == "" {
				return errors.New("invalid todo")
			}

			todo.Done = true

			return nil
		},
	}
	h := handler.NewServer(svc)

	t.Run("invalid id", func(t *testing.T) {
		svc.GetInvoked = 0
		svc.MarkAsDoneInvoked = 0

		req := httptest.NewRequest("PUT", "/abc/done", nil)
		rr := execute(h, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, 0, svc.GetInvoked)
		assert.Equal(t, 0, svc.MarkAsDoneInvoked)
	})

	t.Run("non-existent todo", func(t *testing.T) {
		svc.GetInvoked = 0
		svc.MarkAsDoneInvoked = 0

		req := httptest.NewRequest("PUT", "/100/done", nil)
		rr := execute(h, req)

		assert.Equal(t, http.StatusNotFound, rr.Code)
		assert.Equal(t, 1, svc.GetInvoked)
		assert.Equal(t, 0, svc.MarkAsDoneInvoked)
	})

	t.Run("existent todo", func(t *testing.T) {
		svc.GetInvoked = 0
		svc.MarkAsDoneInvoked = 0

		req := httptest.NewRequest("PUT", "/1/done", nil)
		rr := execute(h, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, 1, svc.GetInvoked)
		assert.Equal(t, 1, svc.MarkAsDoneInvoked)
	})
}

func TestDelete(t *testing.T) {
	svc := &mockService{
		GetFn: func(id int) (*gotodo.Todo, error) {
			if id != 1 {
				return nil, errors.New("not found")
			}

			return &gotodo.Todo{1, "title", "description", false}, nil
		},
		DeleteFn: func(todo *gotodo.Todo) error {
			return nil
		},
	}
	h := handler.NewServer(svc)

	t.Run("invalid id", func(t *testing.T) {
		svc.GetInvoked = 0
		svc.DeleteInvoked = 0

		req := httptest.NewRequest("DELETE", "/abc", nil)
		rr := execute(h, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, 0, svc.GetInvoked)
		assert.Equal(t, 0, svc.DeleteInvoked)
	})

	t.Run("non-existent todo", func(t *testing.T) {
		svc.GetInvoked = 0
		svc.DeleteInvoked = 0

		req := httptest.NewRequest("DELETE", "/100", nil)
		rr := execute(h, req)

		assert.Equal(t, http.StatusNotFound, rr.Code)
		assert.Equal(t, 1, svc.GetInvoked)
		assert.Equal(t, 0, svc.DeleteInvoked)
	})

	t.Run("existent todo", func(t *testing.T) {
		svc.GetInvoked = 0
		svc.DeleteInvoked = 0

		req := httptest.NewRequest("DELETE", "/1", nil)
		rr := execute(h, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, 1, svc.GetInvoked)
		assert.Equal(t, 1, svc.DeleteInvoked)
	})
}

func TestDeleteFinished(t *testing.T) {
	svc := &mockService{
		DeleteFinishedFn: func() {},
	}
	h := handler.NewServer(svc)

	t.Run("delete all", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/", nil)
		rr := execute(h, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("delete pending only", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/?done=false", nil)
		rr := execute(h, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("delete finished only", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/?done=true", nil)
		rr := execute(h, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})
}

// Implements gotodo.Service interface
type mockService struct {
	GetFn      func(id int) (*gotodo.Todo, error)
	GetInvoked int

	GetAllFn      func() []*gotodo.Todo
	GetAllInvoked int

	GetPendingFn      func() []*gotodo.Todo
	GetPendingInvoked int

	GetFinishedFn      func() []*gotodo.Todo
	GetFinishedInvoked int

	AddFn      func(title string, description string) (*gotodo.Todo, error)
	AddInvoked int

	EditFn      func(todo *gotodo.Todo) error
	EditInvoked int

	MarkAsDoneFn      func(todo *gotodo.Todo) error
	MarkAsDoneInvoked int

	DeleteFn      func(todo *gotodo.Todo) error
	DeleteInvoked int

	DeleteFinishedFn      func()
	DeleteFinishedInvoked int
}

func (s *mockService) Get(id int) (*gotodo.Todo, error) {
	s.GetInvoked++
	return s.GetFn(id)
}

func (s *mockService) GetAll() []*gotodo.Todo {
	s.GetAllInvoked++
	return s.GetAllFn()
}

func (s *mockService) GetPending() []*gotodo.Todo {
	s.GetPendingInvoked++
	return s.GetPendingFn()
}

func (s *mockService) GetFinished() []*gotodo.Todo {
	s.GetFinishedInvoked++
	return s.GetFinishedFn()
}

func (s *mockService) Add(title string, description string) (*gotodo.Todo, error) {
	s.AddInvoked++
	return s.AddFn(title, description)
}

func (s *mockService) Edit(todo *gotodo.Todo) error {
	s.EditInvoked++
	return s.EditFn(todo)
}

func (s *mockService) MarkAsDone(todo *gotodo.Todo) error {
	s.MarkAsDoneInvoked++
	return s.MarkAsDoneFn(todo)
}

func (s *mockService) Delete(todo *gotodo.Todo) error {
	s.DeleteInvoked++
	return s.DeleteFn(todo)
}

func (s *mockService) DeleteFinished() {
	s.DeleteFinishedInvoked++
	s.DeleteFinishedFn()
}
