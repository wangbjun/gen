package server

import (
	"context"
	"errors"
	"fmt"
	_ "gen/api"
	"gen/config"
	"gen/log"
	"gen/registry"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"os"
	"sync"
)

// Server is responsible for managing the lifecycle of services.
type Server struct {
	context          context.Context
	shutdownFn       context.CancelFunc
	childRoutines    *errgroup.Group
	log              *zap.Logger
	cfg              *config.Cfg
	shutdownOnce     sync.Once
	shutdownFinished chan struct{}
	isInitialized    bool
	mtx              sync.Mutex
	serviceRegistry  serviceRegistry
}

type serviceRegistry interface {
	IsDisabled(srv registry.Service) bool
	GetServices() []*registry.Descriptor
}

type globalServiceRegistry struct{}

func (r *globalServiceRegistry) IsDisabled(srv registry.Service) bool {
	return registry.IsDisabled(srv)
}

func (r *globalServiceRegistry) GetServices() []*registry.Descriptor {
	return registry.GetServices()
}

// New returns a new instance of Server.
func New(cfgFile string) (*Server, error) {
	rootCtx, shutdownFn := context.WithCancel(context.Background())
	childRoutines, childCtx := errgroup.WithContext(rootCtx)

	s := &Server{
		context:          childCtx,
		shutdownFn:       shutdownFn,
		shutdownFinished: make(chan struct{}),
		childRoutines:    childRoutines,
		log:              log.New("server"),
		cfg:              config.NewConfig(cfgFile),
		serviceRegistry:  &globalServiceRegistry{},
	}
	if err := s.init(); err != nil {
		return nil, err
	}
	return s, nil
}

// init initializes the server and its services.
func (s *Server) init() error {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	if s.isInitialized {
		return nil
	}
	s.isInitialized = true

	if err := s.cfg.Load(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to load cfg: %s\n", err.Error())
		os.Exit(1)
	}

	s.log = log.New("server")

	services := s.serviceRegistry.GetServices()
	if err := s.buildServiceGraph(services); err != nil {
		return err
	}

	return nil
}

// Run initializes and starts services. This will block until all services have
// exited. To initiate shutdown, call the Shutdown method in another goroutine.
func (s *Server) Run() error {
	defer close(s.shutdownFinished)

	if err := s.init(); err != nil {
		return err
	}

	services := s.serviceRegistry.GetServices()

	// Start background services.
	for _, svc := range services {
		service, ok := svc.Instance.(registry.BackgroundService)
		if !ok {
			continue
		}
		if s.serviceRegistry.IsDisabled(svc.Instance) {
			continue
		}
		// Variable is needed for accessing loop variable in callback
		descriptor := svc
		s.childRoutines.Go(func() error {
			select {
			case <-s.context.Done():
				return s.context.Err()
			default:
			}
			err := service.Run(s.context)
			// Do not return context.Canceled error since errgroup.Group only
			// returns the first error to the caller - thus we can miss a more
			// interesting error.
			if err != nil && !errors.Is(err, context.Canceled) {
				s.log.Error("Stopped " + descriptor.Name)
				return fmt.Errorf("%s run error: %w", descriptor.Name, err)
			}
			s.log.Debug("Stopped " + descriptor.Name)
			return nil
		})
	}
	s.log.Debug("Waiting on services...")
	return s.childRoutines.Wait()
}

// Shutdown initiates Grafana graceful shutdown. This shuts down all
// running background services. Since Run blocks Shutdown supposed to
// be run from a separate goroutine.
func (s *Server) Shutdown(ctx context.Context, reason string) error {
	var err error
	s.shutdownOnce.Do(func() {
		s.log.Info("Shutdown started")
		// Call cancel func to stop services.
		s.shutdownFn()
		// Wait for server to shut down
		select {
		case <-s.shutdownFinished:
			s.log.Debug("Finished waiting for server to shut down")
		case <-ctx.Done():
			s.log.Warn("Timed out while waiting for server to shut down")
			err = fmt.Errorf("timeout waiting for shutdown")
		}
	})

	return err
}

// buildServiceGraph builds a graph of services and their dependencies.
func (s *Server) buildServiceGraph(services []*registry.Descriptor) error {
	// Specify service dependencies.
	objs := []interface{}{
		s.cfg,
		s,
	}
	return registry.BuildServiceGraph(objs, services)
}
