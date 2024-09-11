package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/partridge1307/gofiber/entity"
)

const (
	QUERY_GET_USERS        = "SELECT * FROM users"
	QUERY_GET_USER_BY_ID   = "SELECT * FROM users WHERE id = $1"
	QUERY_GET_USER_BY_NAME = "SELECT * FROM users WHERE username = $1"
	QUERY_CREATE_USER      = "INSERT INTO users (username, password) VALUES ($1, $2)"
)

type PostgresRepo struct {
	db *pgxpool.Pool
}

func ConnectToPostgres(conn_url string) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pgPool, err := pgxpool.New(ctx, conn_url)
	if err != nil {
		return nil, err
	}

	return pgPool, nil
}

func NewPostgresRepo(db *pgxpool.Pool) *PostgresRepo {
	return &PostgresRepo{
		db,
	}
}

func (repo *PostgresRepo) GetUsers(ctx context.Context) (*[]entity.User, error) {
	rows, err := repo.db.Query(ctx, QUERY_GET_USERS)
	if err != nil {
		return nil, err
	}

	users, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (entity.User, error) {
		user := entity.User{}

		err := row.Scan(&user.ID, &user.Username, &user.Password, &user.Balance)
		if err != nil {
			return user, err
		}

		return user, nil
	})
	if err != nil {
		return nil, err
	}

	return &users, nil
}

func (repo *PostgresRepo) GetUser(ctx context.Context, u interface{}) (*entity.User, error) {
	var row pgx.Row

	switch u.(type) {
	case int:
		row = repo.db.QueryRow(ctx, QUERY_GET_USER_BY_ID, u)
	case string:
		row = repo.db.QueryRow(ctx, QUERY_GET_USER_BY_NAME, u)
	default:
		return nil, errors.New("invalid type. Expected 'int' or 'string'")
	}

	user := &entity.User{}

	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.Balance)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (repo *PostgresRepo) CreateUser(ctx context.Context, user *entity.User) error {
	row, err := repo.db.Exec(ctx, QUERY_CREATE_USER, user.Username, user.Password)
	if err != nil {
		return err
	}

	if row.RowsAffected() != 1 {
		return errors.New("record does not created")
	}

	return nil
}
