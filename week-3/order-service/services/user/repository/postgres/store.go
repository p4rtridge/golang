package postgres

import (
	"context"
	"fmt"
	"order_service/internal/core"
	"order_service/services/user/entity"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	GetUsers(ctx context.Context) (*[]entity.User, error)
	GetUserById(ctx context.Context, userId int) (*entity.User, error)
	AddUserBalanceById(ctx context.Context, userId int, balance float32) error
}

var (
	QUERY_GET_USER_BY_ID            = "SELECT * FROM users WHERE id = $1"
	QUERY_GET_USERS                 = "SELECT * FROM users"
	QUERY_UPDATE_USER_BALANCE_BY_ID = "UPDATE users SET balance = COALESCE(balance, 0.0) + $2 WHERE id = $1"
)

type postgresRepo struct {
	db *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) UserRepository {
	return &postgresRepo{
		db,
	}
}

func (repo *postgresRepo) GetUsers(ctx context.Context) (*[]entity.User, error) {
	rows, _ := repo.db.Query(ctx, QUERY_GET_USERS)

	datas, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (entity.User, error) {
		var data entity.User

		err := row.Scan(&data.Id, &data.Username, &data.Password, &data.Balance, &data.CreatedAt, &data.UpdatedAt)
		if err != nil {
			return entity.User{}, err
		}

		return data, nil
	})
	if err != nil {
		fmt.Println("get users", err)
		return nil, err
	}

	return &datas, nil
}

func (repo *postgresRepo) GetUserById(ctx context.Context, userId int) (*entity.User, error) {
	var data entity.User

	err := repo.db.QueryRow(ctx, QUERY_GET_USER_BY_ID, userId).Scan(&data.Id, &data.Username, &data.Password, &data.Balance, &data.CreatedAt, &data.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, core.ErrRecordNotFound
		}
		return nil, err
	}

	return &data, nil
}

func (repo *postgresRepo) AddUserBalanceById(ctx context.Context, userId int, balance float32) error {
	_, err := repo.db.Exec(ctx, QUERY_UPDATE_USER_BALANCE_BY_ID, userId, balance)
	if err != nil {
		return err
	}

	return nil
}
