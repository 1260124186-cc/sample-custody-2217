package httpapi

import (
	"net/http"
	"strings"

	"example.com/sample-custody/internal/audit"
	"example.com/sample-custody/internal/metrics"
	"example.com/sample-custody/internal/service"
	"example.com/sample-custody/internal/store"
)

type Server struct {
	services service.Services
	store    *store.Store
	metrics  *metrics.Counters
	audit    *audit.Log
	mux      *http.ServeMux
}

func New(services service.Services, store *store.Store, metrics *metrics.Counters, auditLog *audit.Log) *Server {
	server := &Server{
		services: services,
		store:    store,
		metrics:  metrics,
		audit:    auditLog,
		mux:      http.NewServeMux(),
	}
	server.routes()
	return server
}

func (s *Server) routes() {
	s.mux.HandleFunc("GET /healthz", s.health)
	s.mux.HandleFunc("GET /samples", s.listSamples)
	s.mux.HandleFunc("POST /samples", s.registerSample)
	s.mux.HandleFunc("GET /samples/{id}", s.sampleDetail)
	s.mux.HandleFunc("POST /samples/{id}/transfers", s.transferSample)
	s.mux.HandleFunc("POST /samples/{id}/seal", s.sealSample)
	s.mux.HandleFunc("POST /batches", s.createBatch)
	s.mux.HandleFunc("GET /batches/{id}", s.batchDetail)
	s.mux.HandleFunc("POST /batches/{id}/complete", s.completeBatch)
	s.mux.HandleFunc("GET /manifests/{id}", s.manifest)
	s.mux.HandleFunc("GET /summary", s.summary)
	s.mux.HandleFunc("GET /search", s.searchSamples)
	s.mux.HandleFunc("GET /exports/samples", s.exportSamples)
}

func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.URL.Path != "/" {
		request.URL.Path = strings.TrimRight(request.URL.Path, "/")
		if request.URL.Path == "" {
			request.URL.Path = "/"
		}
	}
	s.mux.ServeHTTP(writer, request)
}
