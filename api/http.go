package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/MihaiBlebea/dog-ceo/template"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type server struct {
	templateService template.Service
	handler         http.Handler
	logger          *logrus.Logger
}

// Server interface
type Server interface {
	Handler() *http.Handler
}

// New returns a new http server service
func New(templateService template.Service, logger *logrus.Logger) Server {
	httpServer := server{
		templateService: templateService,
		logger:          logger,
	}

	r := mux.NewRouter()

	r.Methods("GET").Path("/").HandlerFunc(httpServer.indexHandler)

	r.PathPrefix("/static/").Handler(
		http.StripPrefix(
			"/static/",
			http.FileServer(
				http.Dir(
					httpServer.staticFolderPath(),
				),
			),
		),
	)

	httpServer.handler = r

	return &httpServer
}

func (h *server) Handler() *http.Handler {
	return &h.handler
}

func (h *server) staticFolderPath() string {
	p, err := os.Executable()
	if err != nil {
		h.logger.Fatal(err)
	}

	absPath := fmt.Sprintf(
		"%s/%s/",
		path.Dir(p),
		"static",
	)
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		h.logger.Fatal(err)
	}

	return absPath
}

func (h *server) indexHandler(w http.ResponseWriter, r *http.Request) {
	h.logger.WithFields(logrus.Fields{
		"url": r.URL.String(),
	}).Info("HTTP request started")
	start := time.Now()

	defer h.logger.WithFields(logrus.Fields{
		"duration": time.Since(start).Nanoseconds(),
	}).Info("HTTP request ended")

	path := h.staticFolderPath() + "html/index.gohtml"
	page, err := h.templateService.Load(path)

	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
	}

	page.Render(w, time.Since(start))
}
