package queries

import (
	"WB_Intern/internal/model"
	"context"
)

const insertNewOrder = `insert into orders (data) values ($1)`
func (q *Queries) SaveOrder(ctx context.Context, order model.Order) error {
	q.mu.Lock()

	_, err := q.pool.Exec(ctx, insertNewOrder, order)
	if err != nil {
		q.mu.Unlock()
		return err
	}

	q.mu.Unlock()
	return nil
}

const selectOrders = `select data from orders`
func (q *Queries) GetOrders(ctx context.Context) (map[string]*model.Order, error) {
	orders := make(map[string]*model.Order)

	// TODO CHECK
	q.mu.Lock()
	defer q.mu.Unlock()

	rows, err := q.pool.Query(ctx, selectOrders)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		order := &model.Order{}
		if err = rows.Scan(order); err != nil {
			return nil, err
		}
		orders[order.OrderUID] = order
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

