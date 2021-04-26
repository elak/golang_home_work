package internalhttp

import (
	"encoding/json"
	"net/http"

	domain "github.com/elak/golang_home_work/hw16_playlist_maker/internal/storage"
	"github.com/gorilla/mux"
)

type HistoryHandler struct {
	app Application
}

func (h *HistoryHandler) HandleList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		SendError(w, http.StatusMethodNotAllowed, ErrWrongMethod)
	}

	list, err := h.app.ListHistory(r.Context())
	if err != nil {
		SendError(w, http.StatusInternalServerError, err)
	}
	Send(w, http.StatusOK, list)
}

func (h *HistoryHandler) HandleItem(w http.ResponseWriter, r *http.Request) {
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

func (h *HistoryHandler) handleCreate(w http.ResponseWriter, r *http.Request, id domain.UUID) {
	defer func() {
		err := r.Body.Close()
		if err != nil {
			SendError(w, http.StatusInternalServerError, err)
		}
	}()
	var item domain.History
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		SendError(w, http.StatusBadRequest, err)
		return
	}
	if id != item.VideoID {
		SendError(w, http.StatusBadRequest, err)
		return
	}

	err = h.app.CreateHistory(r.Context(), item)
	if err != nil {
		SendError(w, http.StatusInternalServerError, err)
		return
	}

	Send(w, http.StatusCreated, item)
}

func (h *HistoryHandler) handleRead(w http.ResponseWriter, r *http.Request, id domain.UUID) {
	item, err := h.app.ReadHistory(r.Context(), id)
	if err != nil {
		SendError(w, http.StatusInternalServerError, err)
		return
	}

	Send(w, http.StatusOK, item)
}

func (h *HistoryHandler) handleUpdate(w http.ResponseWriter, r *http.Request, id domain.UUID) {
	defer func() {
		err := r.Body.Close()
		if err != nil {
			SendError(w, http.StatusInternalServerError, err)
		}
	}()
	var item domain.History
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		SendError(w, http.StatusBadRequest, err)
		return
	}
	if id != item.VideoID {
		SendError(w, http.StatusBadRequest, err)
		return
	}

	err = h.app.UpdateHistory(r.Context(), item)
	if err != nil {
		SendError(w, http.StatusInternalServerError, err)
		return
	}

	Send(w, http.StatusOK, item)
}

func (h *HistoryHandler) handleDelete(w http.ResponseWriter, r *http.Request, id domain.UUID) {
	err := h.app.DeleteHistory(r.Context(), id)
	if err != nil {
		SendError(w, http.StatusInternalServerError, err)
		return
	}

	Send(w, http.StatusOK, map[string]domain.UUID{"id": id})
}
