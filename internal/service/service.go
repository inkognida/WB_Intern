package services

import (
	"WB_Intern/internal/handlers"
	"WB_Intern/internal/model"
	"context"
	"github.com/sirupsen/logrus"
	"github.com/go-chi/chi"
	"net/http"
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
	Close() error
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

func (s *Service) Run() {
	// broker action
	if err := s.broker.Subscribe(context.Background(), subject); err != nil {
		s.logger.Fatalln(err)
	}

	// server action
	router := chi.NewRouter()
	handler := handlers.NewOrderHandler(s.logger, s.cache)

	router.Get("/", handler.ShowForm)
	router.Post("/order", handler.ShowOrder)

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		s.logger.Fatalln("Failed to listen and serve", err)
	}
}

