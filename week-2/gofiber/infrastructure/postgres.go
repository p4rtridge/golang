package infrastructure

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/partridge1307/gofiber/config"
	"github.com/partridge1307/gofiber/entities"
)

const (
	QUERY_GET_USERS        = "SELECT * FROM users"
	QUERY_GET_USER_BY_ID   = "SELECT * FROM users WHERE id = $1"
	QUERY_GET_USER_BY_NAME = "SELECT * FROM users WHERE username = $1"
	QUERY_CREATE_USER      = "INSERT INTO users (username, password) VALUES ($1, $2)"
)

type postgresRepo struct {
	db *pgxpool.Pool
}

func NewPostgresRepo(cfg *config.Config) (*postgresRepo, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pgConfig, err := pgxpool.ParseConfig(cfg.URL)
	if err != nil {
		return nil, err
	}

	var e error
	for i := 0; i < cfg.MAX_RETRIES; i++ {
		pool, err := pgxpool.NewWithConfig(ctx, pgConfig)
		if err == nil {
			return &postgresRepo{
				db: pool,
			}, nil
		}

		e = err

		time.Sleep(2 * time.Second)
	}

	return nil, e
}

func (r *postgresRepo) GetUsers(ctx context.Context) (*[]entities.User, error) {
	rows, err := r.db.Query(ctx, QUERY_GET_USERS)
	if err != nil {
		return nil, err
	}

	users, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (entities.User, error) {
		var user entities.User

		err := row.Scan(&user.ID, &user.Username, &user.Password, &user.Balance)
		if err != nil {
			return entities.User{}, err
		}

		return user, nil
	})
	if err != nil {
		return nil, err
	}

	return &users, nil
}

func (r *postgresRepo) GetUser(ctx context.Context, u interface{}) (*entities.User, error) {
	var row pgx.Row
	switch u.(type) {
	case int:
		row = r.db.QueryRow(ctx, QUERY_GET_USER_BY_ID, u)
	case string:
		row = r.db.QueryRow(ctx, QUERY_GET_USER_BY_NAME, u)
	}

	if row == nil {
		return nil, errors.New("invalid type, 'string' or 'int' expected")
	}

	var user entities.User
	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.Balance)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *postgresRepo) CreateUser(ctx context.Context, user *entities.User) error {
	result, err := r.db.Exec(ctx, QUERY_CREATE_USER, user.Username, user.Password)
	if err != nil {
		return err
	}

	if result.RowsAffected() != 1 {
		return errors.New("could not insert user")
	}

	return nil
}
