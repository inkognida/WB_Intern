package repository

import (
	"WB_Intern/internal/model"
	"WB_Intern/internal/repository/queries"
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

// repo
type repo struct {
	logger *logrus.Logger
	pool   *pgxpool.Pool
	*queries.Queries
}

// connection to db, dsn - option
func dbSetup(logger *logrus.Logger) *pgxpool.Pool {
	dsn := "postgres://admin:admin@localhost:5442/orders?sslmode=disable"

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		logger.Fatalln(err)
	}

	pool, err := pgxpool.ConnectConfig(context.Background(), poolConfig)
	if err != nil {
		logger.Fatalln(err)
	}

	return pool
}

// NewRepo repo
func NewRepo(logger *logrus.Logger) Repository {
	pool := dbSetup(logger)

	return &repo{
		logger: logger,
		pool:   pool,
		Queries: queries.New(pool),
	}
}

type Repository interface {
	// SaveOrder save
	SaveOrder(context.Context, model.Order) error

	// GetOrders also for cache sync
	GetOrders(context.Context) (map[string]*model.Order, error)

	// Close stop db process
	Close(context.Context) error
}

