package api

import (
	"context"
	_ "embed"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/parMaster/logserver/app/config"
	"github.com/parMaster/logserver/app/store"
	"github.com/parMaster/logserver/app/web"
)

// As soon as we need another API (GRPC, GraphQL, etc), we can:
// 1. Create an interface for it
// 2. Create a factory function that returns an instance of the interface
// 3. Replace all references to ApiServer with the interface
// 4. Replace the factory function with a switch statement that returns the appropriate instance
// 5. Replace the NewApiServer call with a call to the factory function

type ApiServer struct {
	ctx    context.Context
	config config.Config
	store  store.Storer
}

func NewApiServer(ctx context.Context, config config.Config) *ApiServer {
	l := &ApiServer{
		ctx:    ctx,
		config: config,
	}

	// Initialize database
	err := store.Load(ctx, config, &l.store)
	if err != nil {
		log.Fatalf("Can't configure database %e", err)
	}

	return l
}

func (l *ApiServer) Start() error {
	httpServer := &http.Server{
		Addr:              l.config.BindAddr,
		Handler:           l.router(),
		ReadHeaderTimeout: time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       time.Second,
	}

	log.Printf("[INFO] Starting http server on %s", l.config.BindAddr)

	go func() {
		<-l.ctx.Done()
		log.Printf("[INFO] Terminating http server")

		if err := httpServer.Close(); err != nil {
			log.Printf("[ERROR] failed to close http server, %v", err)
		}
	}()

	httpServer.ListenAndServe()

	return nil
}

func (l *ApiServer) router() http.Handler {
	router := chi.NewRouter()

	router.Get("/api/v1/check", l.HandleCheck)

	router.Get("/web/chart_tpl.min.js", func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte(web.Chart_tpl_min_js))
	})

	router.Get("/view", func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte(web.View_html))
	})

	router.Get("/viewData/{module}", func(rw http.ResponseWriter, r *http.Request) {
		if l.store == nil {
			rw.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		module := chi.URLParam(r, "module")
		if module == "" {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		rw.Header().Set("Content-Type", "application/json")
		rw.Header().Set("Access-Control-Allow-Origin", "*")
		out, err := l.store.View(module)
		if err != nil {
			log.Printf("[ERROR] Failed to get view: %v", err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		json.NewEncoder(rw).Encode(out)
	})

	return router
}

func (l *ApiServer) HandleCheck(w http.ResponseWriter, r *http.Request) {
	log.Printf("GET /api/v1/check")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{"status": "ok"}); err != nil {
		log.Printf("[ERROR] %s", err.Error())
	}
}
