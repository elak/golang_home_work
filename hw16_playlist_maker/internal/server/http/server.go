package internalhttp

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	domain "github.com/elak/golang_home_work/hw16_playlist_maker/internal/storage"
	"github.com/gorilla/mux"
)

var ErrWrongMethod = errors.New("request method not supported")

type Server struct {
	Port     string
	Host     string
	Instance *http.Server
	App      Application

	categories CategoriesHandler
	groups     GroupsHandler
	videos     VideosHandler
	templates  TemplatesHandler
	history    HistoryHandler
}

type Application interface {
	CreateGroup(ctx context.Context, item domain.Group) error
	ReadGroup(ctx context.Context, id domain.UUID) (*domain.Group, error)
	UpdateGroup(ctx context.Context, item domain.Group) error
	DeleteGroup(ctx context.Context, id domain.UUID) error
	ListGroup(ctx context.Context) ([]*domain.Group, error)

	CreateTemplate(ctx context.Context, item domain.Template) error
	ReadTemplate(ctx context.Context, id domain.UUID) (*domain.Template, error)
	UpdateTemplate(ctx context.Context, item domain.Template) error
	DeleteTemplate(ctx context.Context, id domain.UUID) error
	ListTemplate(ctx context.Context) ([]*domain.Template, error)

	CreateCategory(ctx context.Context, item domain.Category) error
	ReadCategory(ctx context.Context, id domain.UUID) (*domain.Category, error)
	UpdateCategory(ctx context.Context, item domain.Category) error
	DeleteCategory(ctx context.Context, id domain.UUID) error
	ListCategory(ctx context.Context) ([]*domain.Category, error)

	CreateVideo(ctx context.Context, item domain.Video) error
	UpdateVideo(ctx context.Context, item domain.Video) error
	ReadVideo(ctx context.Context, id domain.UUID) (*domain.Video, error)
	DeleteVideo(ctx context.Context, id domain.UUID) error
	ListVideo(ctx context.Context) ([]*domain.Video, error)

	CreateHistory(ctx context.Context, item domain.History) error
	UpdateHistory(ctx context.Context, item domain.History) error
	ReadHistory(ctx context.Context, id domain.UUID) (*domain.History, error)
	DeleteHistory(ctx context.Context, id domain.UUID) error
	ListHistory(ctx context.Context) ([]*domain.History, error)

	MakePlayList(ctx context.Context, duration int, tlpID domain.UUID) ([]*domain.Video, error)
}

func NewServer(app Application, host string, port string) *Server {
	var s Server

	s.Host = host
	s.Port = port

	s.categories.app = app
	s.groups.app = app
	s.videos.app = app
	s.templates.app = app
	s.history.app = app

	r := mux.NewRouter()
	r.HandleFunc("/categories", s.categories.HandleList).Methods("GET")
	r.HandleFunc("/categories/{id}", s.categories.HandleItem)
	r.HandleFunc("/groups", s.groups.HandleList).Methods("GET")
	r.HandleFunc("/groups/{id}", s.groups.HandleItem)
	r.HandleFunc("/videos", s.videos.HandleList).Methods("GET")
	r.HandleFunc("/videos/{id}", s.videos.HandleItem)
	r.HandleFunc("/templates", s.templates.HandleList).Methods("GET")
	r.HandleFunc("/templates/{id}", s.templates.HandleItem)
	r.HandleFunc("/templates/{template}/make_playlist/{duration:[0-9]+}", s.HandlePlaylist).Methods("GET")
	r.HandleFunc("/videos/history", s.history.HandleList)
	r.HandleFunc("/videos/{id}/history", s.history.HandleItem)

	http.Handle("/", r)

	s.Instance = &http.Server{Addr: s.Host + ":" + s.Port}
	s.App = app

	return &s
}

func (s *Server) HandlePlaylist(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	templateId := domain.UUID(vars["template"])
	duration, err := strconv.Atoi(vars["duration"])
	if err != nil {
		SendError(w, http.StatusBadRequest, err)
		return
	}

	list, err := s.App.MakePlayList(r.Context(), duration, templateId)
	if err != nil {
		SendError(w, http.StatusInternalServerError, err)
		return
	}

	Send(w, http.StatusOK, list)
}

func (s *Server) Start(ctx context.Context) error {
	err := s.Instance.ListenAndServe()

	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	<-ctx.Done()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	return s.Instance.Shutdown(ctx)
}

func Send(w http.ResponseWriter, code int, payload interface{}) {
	if payload != nil {
		d, err := json.Marshal(payload)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(code)
		_, err = w.Write(d)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func SendError(w http.ResponseWriter, code int, err error) {
	Send(w, code, map[string]string{"error": err.Error()})
}
