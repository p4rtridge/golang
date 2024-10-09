package postgres

import (
	"context"
	"order_service/internal/core"
	authEntity "order_service/services/auth/entity"
	userEntity "order_service/services/user/entity"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepository interface {
	AddAuth(ctx context.Context, data authEntity.Auth) error
	GetAuth(ctx context.Context, username string) (*userEntity.User, error)
}

const (
	QUERY_GET_USER_BY_USERNAME = "SELECT id, username, password, role, balance, created_at, updated_at FROM users WHERE username = $1"
	QUERY_INSERT_USER          = "INSERT INTO users (username, password, role) VALUES ($1, $2, $3)"
)

type postgresRepo struct {
	db *pgxpool.Pool
}

func NewAuthRepo(db *pgxpool.Pool) AuthRepository {
	return &postgresRepo{
		db,
	}
}

func (repo *postgresRepo) GetAuth(ctx context.Context, username string) (*userEntity.User, error) {
	var u userEntity.User

	err := repo.db.QueryRow(ctx, QUERY_GET_USER_BY_USERNAME, username).Scan(&u.Id, &u.Username, &u.Password, &u.Role, &u.Balance, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, core.ErrRecordNotFound
		}

		return nil, err
	}

	return &u, nil
}

func (repo *postgresRepo) AddAuth(ctx context.Context, data authEntity.Auth) error {
	_, err := repo.db.Exec(ctx, QUERY_INSERT_USER, data.Username, data.Password, data.Role)
	if err != nil {
		return err
	}

	return nil
}
