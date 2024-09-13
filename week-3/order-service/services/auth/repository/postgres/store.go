package postgres

import (
	"context"
	"order_service/internal/core"
	authEntity "order_service/services/auth/entity"
	"order_service/services/user/entity"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	QUERY_GET_USER_BY_USERNAME = "SELECT * FROM users WHERE username = $1"
	QUERY_INSERT_USER          = "INSERT INTO users (username, password) VALUES ($1, $2)"
)

type postgresRepo struct {
	db *pgxpool.Pool
}

func NewPostgresRepo(db *pgxpool.Pool) *postgresRepo {
	return &postgresRepo{
		db,
	}
}

func (repo *postgresRepo) GetAuth(ctx context.Context, username string) (*entity.User, error) {
	var data entity.User

	err := repo.db.QueryRow(ctx, QUERY_GET_USER_BY_USERNAME, username).Scan(&data.Id, &data.Username, &data.Password, &data.Balance, &data.CreatedAt, &data.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, core.ErrRecordNotFound
		}

		return nil, err
	}

	return &data, nil
}

func (repo *postgresRepo) AddAuth(ctx context.Context, data *authEntity.Auth) error {
	_, err := repo.db.Exec(ctx, QUERY_INSERT_USER, data.Username, data.Password)
	if err != nil {
		return err
	}

	return nil
}
