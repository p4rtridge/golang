package postgres

import (
	"context"
	"order_service/services/order/entity"

	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, data *entity.Order) error
}

type postgresRepo struct {
	db *pgxpool.Pool
}

func NewOrderRepo(db *pgxpool.Pool) OrderRepository {
	return &postgresRepo{
		db,
	}
}

func (repo *postgresRepo) CreateOrder(ctx context.Context, data *entity.Order) error {
	return nil
}
