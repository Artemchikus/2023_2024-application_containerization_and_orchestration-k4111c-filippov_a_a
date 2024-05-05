package api

import (
	"context"
	"encoding/json"
	"errors"
	"find-ship/cmd/app/metrix"
	vkclient "find-ship/cmd/app/vk_client"
	"find-ship/config"
	"find-ship/storage"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type (
	server struct {
		router   *mux.Router
		storage  storage.Storage
		metrix   *metrix.Metrix
		vkClient vkclient.VK
	}
)

func Start(ctx context.Context, storage storage.Storage, conf *config.AppConfig) error {
	handler := &server{
		router:   mux.NewRouter(),
		storage:  storage,
		metrix:   metrix.New(ctx, storage),
		vkClient: vkclient.MustVKClient(),
	}
	handler.configureRouter()

	srv := &http.Server{
		Addr:    ":" + conf.Port,
		Handler: handler,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen and serve returned err: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("got interruption signal")

	if err := srv.Shutdown(context.TODO()); err != nil {
		log.Printf("server shutdown returned an err: %v\n", err)
	}

	storage.Disconnect()

	return nil
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	log.Printf("error: %v", err)

	resp := ErrorResponse{
		Error: err.Error(),
	}
	s.respond(w, r, code, resp)
}

func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)

	if data != nil {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(data)
	}
}
