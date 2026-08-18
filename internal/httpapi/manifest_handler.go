package httpapi

import "net/http"

func (s *Server) manifest(writer http.ResponseWriter, request *http.Request) {
	manifest, err := s.services.Queries.Manifest(request.PathValue("id"))
	if err != nil {
		writeError(writer, err)
		return
	}
	writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	writer.WriteHeader(http.StatusOK)
	_, _ = writer.Write([]byte(manifest))
}
