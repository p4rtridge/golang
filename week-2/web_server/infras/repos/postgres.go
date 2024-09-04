package repos

import (
	"context"
	"time"
	"web_server/config"
	"web_server/core/entities"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// frameworks and drivers layer
type PostgresRepo struct {
	DB *pgxpool.Pool
}

func NewPostgresRepo(cfg *config.Config, times int) (*PostgresRepo, error) {
	var e error

	for i := 0; i < times; i++ {
		dbPool, err := pgxpool.New(context.Background(), cfg.URL)
		if err == nil {
			return &PostgresRepo{
				DB: dbPool,
			}, err
		}

		time.Sleep(2 * time.Second)

		e = err
	}

	return nil, e
}

func (repo *PostgresRepo) GetTask(id int) (*entities.Task, error) {
	rows, _ := repo.DB.Query(context.Background(), "SELECT * FROM tasks WHERE id = $1", id)

	task, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entities.Task])
	if err != nil {
		return nil, err
	}

	return &task, nil
}

const QUERY string = "INSERT INTO tasks (title, completed) VALUES (@title, @completed) RETURNING id, title, completed"

func (repo *PostgresRepo) SaveTask(task *entities.Task) (*entities.Task, error) {
	args := pgx.NamedArgs{
		"title":     task.Title,
		"completed": task.Completed,
	}

	var newTask entities.Task
	err := repo.DB.QueryRow(context.Background(), QUERY, args).Scan(&newTask.Id, &newTask.Title, &newTask.Completed)
	if err != nil {
		return nil, err
	}

	return &newTask, nil
}
