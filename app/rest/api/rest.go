package api

import (
	"context"
	"fmt"
	"github.com/reddi/agassi/app/store/service"
	"net/http"
	"sync"
	"time"

	log "github.com/go-pkgz/lgr"

	"github.com/go-chi/chi/v5"
)

const hardBodyLimit = 1024 * 64 // limit size of body

// Rest is a rest access server
type Rest struct {
	DataService *service.DataStore

	pubRest public

	httpServer *http.Server
	lock       sync.Mutex
}

// Run the lister and request's router, activate rest server
func (s *Rest) Run(address string, port int) {
	if address == "*" {
		address = ""
	}

	log.Printf("[INFO] activate http rest server on %s:%d", address, port)

	s.lock.Lock()
	s.httpServer = s.makeHTTPServer(address, port, s.routes())
	s.httpServer.ErrorLog = log.ToStdLogger(log.Default(), "WARN")
	s.lock.Unlock()

	err := s.httpServer.ListenAndServe()
	log.Printf("[WARN] http server terminated, %s", err)
}

// Shutdown rest http server
func (s *Rest) Shutdown() {
	log.Print("[WARN] shutdown rest server")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	s.lock.Lock()
	if s.httpServer != nil {
		if err := s.httpServer.Shutdown(ctx); err != nil {
			log.Printf("[DEBUG] http shutdown error, %s", err)
		}
		log.Print("[DEBUG] shutdown http server completed")
	}
}

func (s *Rest) controllerGroups() public {
	pubGrp := public{
		dataService: s.DataService,
	}
	return pubGrp
}

func (s *Rest) routes() chi.Router {
	router := chi.NewRouter()

	s.pubRest = s.controllerGroups() // assign controllers for groups

	// api routes
	router.Route("/api/v1", func(rapi chi.Router) {
		// open routes
		rapi.Group(func(ropen chi.Router) {
			ropen.Post("/player", s.pubRest.createPlayerCtrl)
			ropen.Post("/coach", s.pubRest.createCoachCtrl)
		})
	})

	return router
}

func (s *Rest) makeHTTPServer(address string, port int, router http.Handler) *http.Server {
	return &http.Server{
		Addr:              fmt.Sprintf("%s:%d", address, port),
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       30 * time.Second,
	}
}
