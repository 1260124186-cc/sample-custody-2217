package app

import (
	"net/http"

	"example.com/sample-custody/internal/audit"
	"example.com/sample-custody/internal/clock"
	"example.com/sample-custody/internal/httpapi"
	"example.com/sample-custody/internal/ids"
	"example.com/sample-custody/internal/metrics"
	"example.com/sample-custody/internal/service"
	"example.com/sample-custody/internal/store"
)

type Application struct {
	Handler http.Handler
}

func New() Application {
	repository := store.New()
	log := audit.NewLog()
	counters := metrics.NewCounters()
	services := service.New(repository, log, clock.System{}, ids.New(), counters)
	return Application{
		Handler: httpapi.New(services, repository, counters, log),
	}
}
