package api

import (
	_ "find-ship/cmd/app/api/docs"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Find Ship API
// @version 1.0
// @description This is a vk channles site parser.

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host
func (s *server) configureRouter() {
	s.router.Use(s.setRequestID)
	s.router.Use(s.logRequest)
	s.router.Use(s.collectStatistics)

	s.router.HandleFunc("/health", s.healthCheck).Methods(http.MethodGet)
	s.router.HandleFunc("/metrics", s.getMetrics).Methods(http.MethodGet)

	v1 := s.router.PathPrefix("/api/v1").Subrouter()
	v1.HandleFunc("/ship", s.getShips).Methods(http.MethodGet)
	v1.HandleFunc("/ship/{id:[0-9]+}", s.getShipByID).Methods(http.MethodGet)

	v1.HandleFunc("/ship", s.postShip).Methods(http.MethodPost)
	v1.HandleFunc("/ship/{id:[0-9]+}", s.deleteShip).Methods(http.MethodDelete)
	v1.HandleFunc("/ship/{id:[0-9]+}", s.patchShip).Methods(http.MethodPatch)

	v1.PathPrefix("/docs/").Handler(httpSwagger.WrapHandler).Methods(http.MethodGet)
}
