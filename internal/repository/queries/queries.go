package queries

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"sync"
)

type Queries struct {
	pool *pgxpool.Pool
	mu *sync.Mutex
}

func New(pgxPool *pgxpool.Pool) *Queries {
	return &Queries{pool: pgxPool, mu: &sync.Mutex{}}
}

func (q *Queries) Close(ctx context.Context) error {
	q.pool.Close()
	return nil
}