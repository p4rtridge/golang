package postgres

import (
	"context"
	"order_service/internal/core"
	"order_service/services/user/entity"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var QUERY_GET_USER_BY_ID = "SELECT * FROM users WHERE id = $1"

type postgresRepo struct {
	db *pgxpool.Pool
}

func NewPostgresRepo(db *pgxpool.Pool) *postgresRepo {
	return &postgresRepo{
		db,
	}
}

func (repo *postgresRepo) GetUserById(ctx context.Context, id int) (*entity.User, error) {
	var data entity.User

	err := repo.db.QueryRow(ctx, QUERY_GET_USER_BY_ID, id).Scan(&data.Id, &data.Username, &data.Password, &data.Balance, &data.CreatedAt, &data.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, core.ErrRecordNotFound
		}
		return nil, err
	}

	return &data, nil
}
