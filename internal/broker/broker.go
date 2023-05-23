package broker

import (
	"WB_Intern/internal/model"
	"context"
	"encoding/json"
	"github.com/nats-io/stan.go"
	"github.com/sirupsen/logrus"
)

type broker struct {
	logger *logrus.Logger
	conn stan.Conn
	repo orderRepository
	cache orderCache
}


type Broker interface {
	Subscribe(context.Context, string) error
	Close() error
}

type orderRepository interface {
	SaveOrder(context.Context, model.Order) error
}

type orderCache interface {
	SaveOrder(*model.Order) error
}

const (
	ClusterID = "test-cluster"
	ClientID = "sub"

)

func NewBroker(repo orderRepository, cache orderCache, logger *logrus.Logger) Broker {
	conn, err := stan.Connect(ClusterID, ClientID, stan.NatsURL(stan.DefaultNatsURL))
	if err != nil {
		logger.Fatalln(err)
	}

	return &broker{
		logger: logger,
		conn:   conn,
		repo:   repo,
		cache:  cache,
	}
}

func (b *broker) Close() error {
	err := b.conn.Close()
	if err != nil {
		b.logger.Fatalln(err)
	}

	return nil
}

// Subscribe ; ctx is for future purpose
func (b *broker) Subscribe(ctx context.Context, subject string) error {
	_, err := b.conn.Subscribe(subject, func(msg *stan.Msg) {
		if err := msg.Ack(); err != nil {
			b.logger.Infoln(err)
			return
		}

		order := &model.Order{}
		if err := json.Unmarshal(msg.Data, order); err != nil {
			b.logger.Infoln(err)
			return
		}

		if err := b.cache.SaveOrder(order); err != nil {
			b.logger.Infoln(err)
			return
		}

		if err := b.repo.SaveOrder(context.Background(), *order); err != nil {
			b.logger.Infoln(err)
			return
		}

		b.logger.Log(logrus.InfoLevel, "Order: ", order.OrderUID," saved to db and cache")
	}, stan.SetManualAckMode())

	if err != nil {
		b.logger.Infoln(err)
		return err

	}

	return nil
}

