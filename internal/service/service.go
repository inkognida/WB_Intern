package services

import (
	"WB_Intern/internal/handlers"
	"WB_Intern/internal/model"
	"context"
	"github.com/sirupsen/logrus"
	"github.com/go-chi/chi"
	"net/http"
	"os"
	"os/signal"
	"time"
)

const (
	subject = "order"
)

type orderRepository interface {
	Close(context.Context) error
}

type orderCache interface {
	GetOrderById(string) (*model.Order, error)
}

type orderBroker interface {
	Subscribe(context.Context, string) error
	Close(context.Context) error
}

type Service struct {
	logger *logrus.Logger
	repo orderRepository
	cache orderCache
	broker orderBroker

}

func NewService(logger *logrus.Logger, repo orderRepository, cache orderCache, broker orderBroker) *Service {
	return &Service{
		logger: logger,
		repo: repo,
		cache: cache,
		broker: broker,
	}
}

const (
	addr = ":8080"
)

func (s *Service) ServerExec() {
	router := chi.NewRouter()

	server := &http.Server{
		Addr:              addr,
		Handler:           router,
		ReadTimeout:       10*time.Second,
		WriteTimeout:      10*time.Second,
	}
	handler := handlers.NewOrderHandler(s.logger, s.cache)

	router.Get("/", handler.ShowForm)
	router.Post("/order", handler.ShowOrder)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Fatalln("Server error: %s", err)
		}
	}()

	s.logger.Log(logrus.InfoLevel, "Server started")
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt)

	// Block until a signal
	<-stopChan
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)

	defer func() {
		cancel()
		err := s.repo.Close(context.Background())
		if err != nil {
			s.logger.Fatalln(err)
		}
		err = s.broker.Close(context.Background())
		if err != nil {
			s.logger.Fatalln(err)
		}
	}()

	if err := server.Shutdown(ctx); err != nil {
		s.logger.Fatalln("Server shutdown error: %s", err)
	}
	s.logger.Infoln("Server stopped")
}

func (s *Service) Run() {
	// broker action
	if err := s.broker.Subscribe(context.Background(), subject); err != nil {
		s.logger.Fatalln(err)
	}

	// server action
	s.ServerExec()
}

