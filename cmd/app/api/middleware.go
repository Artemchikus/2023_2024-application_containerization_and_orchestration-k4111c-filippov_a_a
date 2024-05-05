package api

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
)

func (s *server) setRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New().String()
		w.Header().Set("X-Request-ID", id)

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), CtxKeyRequestID, id)))
	})
}

func (s *server) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := log.New(os.Stdout, "", log.LstdFlags)

		logger.Printf("started %s %s, remote_addr: %s, request_id: %s", r.Method, r.RequestURI, r.Context().Value(CtxKeyRequestID), r.RemoteAddr)

		start := time.Now().UTC()
		rw := &responseWriter{w, http.StatusOK}

		defer func() {
			respTime := time.Since(start).Milliseconds()

			logger.Printf("complited with %d %s in %vms, request_id: %s, remote_addr: %s", rw.code, http.StatusText(rw.code), respTime, r.Context().Value(CtxKeyRequestID), r.RemoteAddr)

			s.metrix.RespTime.Set(respTime)
		}()
		
		next.ServeHTTP(rw, r)
	})
}

func (s *server) collectStatistics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.metrix.CurrentReq.Inc()
		defer s.metrix.CurrentReq.Dec()

		s.metrix.AvgReq.Set(s.metrix.CurrentReq.Get().CurrentReq)
		defer s.metrix.AvgReq.Set(s.metrix.CurrentReq.Get().CurrentReq)

		next.ServeHTTP(w, r)
	})
}
