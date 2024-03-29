package http

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/KWRI/demo-service/core/cfg/di"
	"github.com/KWRI/demo-service/core/http/health"
	"github.com/KWRI/demo-service/core/log"
)

//
// ServiceProvider provides a HTTP service.
//
type ServiceProvider interface {
	//
	// Run runs a new HTTP service.
	//
	Run()
}

//
// Service is a HTTP service.
//
type Service struct {
	name         string
	address      string
	handler      http.Handler
	logger       log.Logger
	diContainer  di.ContainerProvider
	readTimeout  time.Duration
	writeTimeout time.Duration
}

//
// NewService returns a new HTTP service instance.
//
func NewService(
	name string,
	address string,
	readTimeout time.Duration,
	writeTimeout time.Duration,
	handler http.Handler,
	logger log.Logger,
	diContainer di.ContainerProvider,
) *Service {

	return &Service{
		name:         name,
		logger:       logger,
		handler:      handler,
		address:      address,
		diContainer:  diContainer,
		readTimeout:  readTimeout,
		writeTimeout: writeTimeout,
	}
}

//
// Run runs a new HTTP service.
//
func (s *Service) Run() {

	s.logger.Info(
		"%s service with build info {date:%s, branch:%s, commit:%s} start listening address:%v",
		s.name,
		health.BuildDate,
		health.BuildBranch,
		health.BuildCommit,
		s.address,
	)

	server := http.Server{
		Addr:         s.address,
		Handler:      s.handler,
		ReadTimeout:  s.readTimeout,
		WriteTimeout: s.writeTimeout,
		IdleTimeout:  s.writeTimeout,
	}

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-ch
		s.logger.Info("%s service starting graceful shutdown...", s.name)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			s.logger.Error("%+v", err)
		}

		if err := s.diContainer.Shutdown(); err != nil {
			s.logger.Error("%+v", err)
		}

		s.logger.Info("%s service has been stopped", s.name)
	}()

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		s.logger.Error(`%+v`, err)
	}
}
