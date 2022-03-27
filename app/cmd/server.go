package cmd

import (
	"context"
	"fmt"
	log "github.com/go-pkgz/lgr"
	"github.com/pkg/errors"
	"github.com/reddi/agassi/app/rest/api"
	"github.com/reddi/agassi/app/store/engine"
	"github.com/reddi/agassi/app/store/service"
	bolt "go.etcd.io/bbolt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type ServerCommand struct {
	Store StoreGroup `group:"store" namespace:"store" env-namespace:"STORE"`

	Site    string `long:"site" env:"SITE" default:"agassi" description:"site name"`
	Port    int    `long:"port" env:"AGASSI_PORT" default:"8080" description:"port"`
	Address string `long:"address" env:"AGASSI_ADDRESS" default:"" description:"listening address"`
}

// StoreGroup defines options group for store params
type StoreGroup struct {
	Type string `long:"type" env:"TYPE" description:"type of storage" choice:"bolt" default:"bolt"` // nolint
	Bolt struct {
		Path    string        `long:"path" env:"PATH" default:"./var" description:"parent directory for the bolt files"`
		Timeout time.Duration `long:"timeout" env:"TIMEOUT" default:"30s" description:"bolt timeout"`
	} `group:"bolt" namespace:"bolt" env-namespace:"BOLT"`
}

type serverApp struct {
	*ServerCommand
	restSrv     *api.Rest
	dataService *service.DataStore
	terminated  chan struct{}
}

func (s *ServerCommand) Execute(_ []string) error {
	log.Printf("[INFO] start server on port %s:%d", s.Address, s.Port)

	ctx, cancel := context.WithCancel(context.Background())
	go func() { // catch signal and invoke graceful termination
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
		<-stop
		log.Printf("[WARN] interrupt signal")
		cancel()
	}()

	app, err := s.newServerApp(ctx)
	if err != nil {
		log.Printf("[PANIC] failed to setup application, %+v", err)
		return err
	}
	if err = app.run(ctx); err != nil {
		log.Printf("[ERROR] agassi terminated with error %+v", err)
		return err
	}
	log.Printf("[INFO] agassi terminated")
	return nil
}

// makeDataStore creates store for all sites
func (s *ServerCommand) makeDataStore() (result engine.Interface, err error) {
	log.Printf("[INFO] make data store, type=%s", s.Store.Type)

	switch s.Store.Type {
	case "bolt":
		if err = makeDirs(s.Store.Bolt.Path); err != nil {
			return nil, errors.Wrap(err, "failed to create bolt store")
		}
		fileName := fmt.Sprintf("%s/%s.db", s.Store.Bolt.Path, s.Site)
		result, err = engine.NewBoltDB(bolt.Options{Timeout: s.Store.Bolt.Timeout}, fileName)
	default:
		return nil, errors.Errorf("unsupported store type %s", s.Store.Type)
	}
	return result, errors.Wrap(err, "can't initialize data store")
}

func (s *ServerCommand) newServerApp(ctx context.Context) (*serverApp, error) {
	storeEngine, err := s.makeDataStore()
	if err != nil {
		return nil, fmt.Errorf("failed to make data store engine: %w", err)
	}
	dataService := &service.DataStore{
		Engine: storeEngine,
	}
	srv := &api.Rest{
		DataService: dataService,
	}
	return &serverApp{
		ServerCommand: s,
		restSrv:       srv,
		dataService:   dataService,
		terminated:    make(chan struct{}),
	}, nil
}

// Run all application objects
func (a *serverApp) run(ctx context.Context) error {

	go func() {
		// shutdown on context cancellation
		<-ctx.Done()
		log.Print("[INFO] shutdown initiated")
		a.restSrv.Shutdown()
	}()

	a.restSrv.Run(a.Address, a.Port)

	if e := a.dataService.Close(); e != nil {
		log.Printf("[WARN] failed to close data store, %s", e)
	}

	close(a.terminated)
	return nil
}

// Wait for application completion (termination)
func (a *serverApp) Wait() {
	<-a.terminated
}
