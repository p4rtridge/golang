package repositories

import (
	"context"
	"user-service/entities"

	"github.com/jmoiron/sqlx"
)

type AuthRepository interface {
	AddAuth(ctx context.Context, data entities.Auth) error
	GetAuth(ctx context.Context, username string) (*entities.User, error)
}

type postgresRepo struct {
	db *sqlx.DB
}

func NewAuthRepo(db *sqlx.DB) AuthRepository {
	return &postgresRepo{
		db,
	}
}

func (repo *postgresRepo) AddAuth(ctx context.Context, data entities.Auth) error {
	query := `INSERT INTO users (username, password, role) VALUES (?, ?, ?)`
	query = sqlx.Rebind(sqlx.DOLLAR, query)

	res, err := repo.db.ExecContext(ctx, query, data.Username, data.Password, data.Role)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return entities.ErrCannotAdd
	}

	return nil
}

func (repo *postgresRepo) GetAuth(ctx context.Context, username string) (*entities.User, error) {
	query := `SELECT id, username, password FROM users WHERE username = ?`
	query = sqlx.Rebind(sqlx.DOLLAR, query)

	var user entities.User

	err := repo.db.GetContext(ctx, &user, query, username)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
