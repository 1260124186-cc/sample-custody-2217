package httpapi

import "net/http"

func (s *Server) health(writer http.ResponseWriter, request *http.Request) {
	counts := s.store.Counts()
	writeJSON(writer, http.StatusOK, map[string]any{
		"status":   "ok",
		"service":  "sample-custody",
		"counts":   counts,
		"metrics":  s.metrics.Snapshot(),
		"activity": s.audit.SummaryByAction(),
	})
}
