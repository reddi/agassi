package cmd

import (
	"context"
	log "github.com/go-pkgz/lgr"
	"github.com/reddi/agassi/app/rest/api"
	"github.com/reddi/agassi/app/store/service"
	"os"
	"os/signal"
	"syscall"
)

type ServerCommand struct {
	Port    int    `long:"port" env:"AGASSI_PORT" default:"8080" description:"port"`
	Address string `long:"address" env:"AGASSI_ADDRESS" default:"" description:"listening address"`
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

func (s *ServerCommand) newServerApp(ctx context.Context) (*serverApp, error) {
	dataService := &service.DataStore{}
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
