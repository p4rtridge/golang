package postgres

import (
	"context"
	"fmt"
	"order_service/internal/core"
	"order_service/services/product/entity"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	QUERY_INSERT_PRODUCT       = "INSERT INTO products (name, quantity, price) VALUES ($1, $2, $3)"
	QUERY_GET_PRODUCTS         = "SELECT * FROM products"
	QUERY_GET_PRODUCT_BY_ID    = "SELECT * FROM products WHERE id = $1"
	QUERY_UPDATE_PRODUCT_BY_ID = "UPDATE products SET name = COALESCE($2, name), quantity = COALESCE($3, quantity), price = COALESCE($4, price), updated_at = $5 WHERE id = $1"
	QUERY_DELETE_PRODUCT_BY_ID = "DELETE FROM products WHERE id = $1"
)

type postgresRepo struct {
	db *pgxpool.Pool
}

func NewPostgresRepo(db *pgxpool.Pool) *postgresRepo {
	return &postgresRepo{
		db,
	}
}

func (repo *postgresRepo) CreateProduct(ctx context.Context, data *entity.Product) error {
	_, err := repo.db.Exec(ctx, QUERY_INSERT_PRODUCT, data.Name, data.Quantity, data.Price)
	if err != nil {
		fmt.Println("create err", err)
		return err
	}

	return nil
}

func (repo *postgresRepo) GetProducts(ctx context.Context) (*[]entity.Product, error) {
	rows, _ := repo.db.Query(ctx, QUERY_GET_PRODUCTS)

	datas, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (entity.Product, error) {
		var data entity.Product

		err := row.Scan(&data.Id, &data.Name, &data.Quantity, &data.Price, &data.CreatedAt, &data.UpdatedAt)
		if err != nil {
			return entity.Product{}, err
		}

		return data, nil
	})
	if err != nil {
		fmt.Println("get products", err)
		return nil, err
	}

	return &datas, nil
}

func (repo *postgresRepo) GetProduct(ctx context.Context, productID int) (*entity.Product, error) {
	var data entity.Product

	err := repo.db.QueryRow(ctx, QUERY_GET_PRODUCT_BY_ID, productID).Scan(&data.Id, &data.Name, &data.Quantity, &data.Price, &data.CreatedAt, &data.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, core.ErrRecordNotFound
		}
		return nil, err
	}

	return &data, nil
}

func (repo *postgresRepo) UpdateProduct(ctx context.Context, productID int, data *entity.Product) error {
	newName := pgtype.Text{Valid: false}
	newQuantity := pgtype.Int4{Valid: false}
	newPrice := pgtype.Float4{Valid: false}

	if data.Name != "" {
		newName = pgtype.Text{String: data.Name, Valid: true}
	}

	if data.Quantity != 0 {
		newQuantity = pgtype.Int4{Int32: int32(data.Quantity), Valid: true}
	}

	if data.Price != 0.0 {
		newPrice = pgtype.Float4{Float32: data.Price, Valid: true}
	}

	_, err := repo.db.Exec(ctx, QUERY_UPDATE_PRODUCT_BY_ID, productID, newName, newQuantity, newPrice, time.Now())
	if err != nil {
		fmt.Println("product update err", err)
		return err
	}

	return nil
}

func (repo *postgresRepo) DeleteProduct(ctx context.Context, productId int) error {
	_, err := repo.db.Exec(ctx, QUERY_DELETE_PRODUCT_BY_ID, productId)
	if err != nil {
		fmt.Println("product delete err", err)
		return err
	}

	return nil
}
