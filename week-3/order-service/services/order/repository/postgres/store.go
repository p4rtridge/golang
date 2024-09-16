package postgres

import (
	"context"
	"order_service/services/order/entity"

	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresRepo struct {
	db *pgxpool.Pool
}

func NewPostgresRepo(db *pgxpool.Pool) *postgresRepo {
	return &postgresRepo{
		db,
	}
}

func (repo *postgresRepo) CreateOrder(ctx context.Context, data *entity.Order) error {
}
