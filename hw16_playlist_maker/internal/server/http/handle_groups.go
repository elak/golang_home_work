package internalhttp

import (
	"encoding/json"
	"net/http"

	domain "github.com/elak/golang_home_work/hw16_playlist_maker/internal/storage"
	"github.com/gorilla/mux"
)

type GroupsHandler struct {
	app Application
}

func (h *GroupsHandler) HandleList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		SendError(w, http.StatusMethodNotAllowed, ErrWrongMethod)
	}

	list, err := h.app.ListGroup(r.Context())
	if err != nil {
		SendError(w, http.StatusInternalServerError, err)
	}
	Send(w, http.StatusOK, list)
}

func (h *GroupsHandler) HandleItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := domain.UUID(vars["id"])
	switch r.Method {
	case http.MethodGet:
		h.handleRead(w, r, id)
	case http.MethodPost:
		h.handleCreate(w, r, id)
	case http.MethodPut:
		h.handleUpdate(w, r, id)
	case http.MethodDelete:
		h.handleDelete(w, r, id)
	default:
		SendError(w, http.StatusMethodNotAllowed, ErrWrongMethod)
	}
}

func (h *GroupsHandler) handleCreate(w http.ResponseWriter, r *http.Request, id domain.UUID) {
	defer func() {
		err := r.Body.Close()
		if err != nil {
			SendError(w, http.StatusInternalServerError, err)
		}
	}()
	var item domain.Group
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		SendError(w, http.StatusBadRequest, err)
		return
	}
	if id != item.ID {
		SendError(w, http.StatusBadRequest, err)
		return
	}

	err = h.app.CreateGroup(r.Context(), item)
	if err != nil {
		SendError(w, http.StatusInternalServerError, err)
		return
	}

	Send(w, http.StatusCreated, item)
}

func (h *GroupsHandler) handleRead(w http.ResponseWriter, r *http.Request, id domain.UUID) {
	item, err := h.app.ReadGroup(r.Context(), id)
	if err != nil {
		SendError(w, http.StatusInternalServerError, err)
		return
	}

	Send(w, http.StatusOK, item)
}

func (h *GroupsHandler) handleUpdate(w http.ResponseWriter, r *http.Request, id domain.UUID) {
	defer func() {
		err := r.Body.Close()
		if err != nil {
			SendError(w, http.StatusInternalServerError, err)
		}
	}()
	var item domain.Group
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		SendError(w, http.StatusBadRequest, err)
		return
	}
	if id != item.ID {
		SendError(w, http.StatusBadRequest, err)
		return
	}

	err = h.app.UpdateGroup(r.Context(), item)
	if err != nil {
		SendError(w, http.StatusInternalServerError, err)
		return
	}

	Send(w, http.StatusOK, item)
}

func (h *GroupsHandler) handleDelete(w http.ResponseWriter, r *http.Request, id domain.UUID) {
	err := h.app.DeleteGroup(r.Context(), id)
	if err != nil {
		SendError(w, http.StatusInternalServerError, err)
		return
	}

	Send(w, http.StatusOK, map[string]domain.UUID{"id": id})
}
