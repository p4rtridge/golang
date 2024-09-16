package postgres

import (
	"context"
	"errors"
	orderEntity "order_service/services/order/entity"
	productEntity "order_service/services/product/entity"
	userEntity "order_service/services/user/entity"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, order *orderEntity.Order, callbackFn func(order *orderEntity.Order, user *userEntity.User, products *[]productEntity.Product) (bool, error)) error
	GetOrders(ctx context.Context) (*[]orderEntity.Order, error)
}

const (
	QUERY_GET_ORDERS                  = "SELECT o.id AS order_id, o.user_id, oi.product_id, oi.product_name, oi.product_price, oi.quantity, o.total_price, o.created_at, o.updated_at FROM orders AS o JOIN order_items AS oi ON o.id = oi.order_id"
	QUERY_GET_USER_LOCK               = "SELECT * FROM users WHERE id = $1 FOR UPDATE"
	QUERY_GET_PRODUCT_LOCK            = "SELECT * FROM products WHERE id = $1 FOR UPDATE"
	QUERY_CREATE_ORDER_WITH_RETURN_ID = "INSERT INTO orders (user_id, total_price) VALUES ($1, $2) RETURNING id"
	QUERY_CREATE_ORDER_ITEM           = "INSERT INTO order_items (order_id, product_id, product_name, product_price, quantity) VALUES ($1, $2, $3, $4, $5)"
	QUERY_UPDATE_USER_BALANCE         = "UPDATE users SET balance = COALESCE($2, balance) WHERE id = $1"
	QUERY_UPDATE_PRODUCT_QUANTITY     = "UPDATE products SET quantity = COALESCE($2, quantity) WHERE id = $1"
)

type postgresRepo struct {
	db *pgxpool.Pool
}

func NewOrderRepo(db *pgxpool.Pool) OrderRepository {
	return &postgresRepo{
		db,
	}
}

func (repo *postgresRepo) CreateOrder(ctx context.Context, order *orderEntity.Order, callbackFn func(order *orderEntity.Order, user *userEntity.User, products *[]productEntity.Product) (bool, error)) error {
	return runInTransaction(ctx, repo.db, func(tx pgx.Tx) error {
		var user userEntity.User

		err := tx.QueryRow(ctx, QUERY_GET_USER_LOCK, order.UserId).Scan(&user.Id, &user.Username, &user.Password, &user.Balance, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return err
		}

		products := make([]productEntity.Product, 0, len(order.Items))
		for _, item := range order.Items {
			var product productEntity.Product

			err := tx.QueryRow(ctx, QUERY_GET_PRODUCT_LOCK, item.ProductId).Scan(&product.Id, &product.Name, &product.Quantity, &product.Price, &product.CreatedAt, &product.UpdatedAt)
			if err != nil {
				return err
			}

			products = append(products, product)
		}

		accept, err := callbackFn(order, &user, &products)
		if err != nil {
			return err
		}
		if !accept {
			return nil
		}

		var newOrderId int

		err = tx.QueryRow(ctx, QUERY_CREATE_ORDER_WITH_RETURN_ID, order.UserId, order.TotalPrice).Scan(&newOrderId)
		if err != nil {
			return err
		}
		order.SetId(newOrderId)

		for idx, item := range order.Items {
			_, err = tx.Exec(ctx, QUERY_CREATE_ORDER_ITEM, order.Id, item.ProductId, item.ProductName, item.ProductPrice, item.Quantity)
			if err != nil {
				return err
			}

			_, err = tx.Exec(ctx, QUERY_UPDATE_PRODUCT_QUANTITY, products[idx].Id, products[idx].Quantity)
			if err != nil {
				return err
			}
		}

		_, err = tx.Exec(ctx, QUERY_UPDATE_USER_BALANCE, user.Id, user.Balance)
		if err != nil {
			return err
		}

		return nil
	})
}

func (repo *postgresRepo) GetOrders(ctx context.Context) (*[]orderEntity.Order, error) {
	rows, err := repo.db.Query(ctx, QUERY_GET_ORDERS)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ordersMap := make(map[int]*orderEntity.Order)

	for rows.Next() {
		var orderId, userId, productId, quantity int
		var productName string
		var productPrice, totalPrice float32
		var createdAt time.Time
		var updatedAt *time.Time

		err := rows.Scan(&orderId, &userId, &productId, &productName, &productPrice, &quantity, &totalPrice, &createdAt, &updatedAt)
		if err != nil {
			return nil, err
		}

		item := orderEntity.NewOrderItem(orderId, productId, productName, productPrice, quantity)

		if _, exists := ordersMap[orderId]; !exists {
			ordersMap[orderId] = &orderEntity.Order{
				Id:         orderId,
				UserId:     userId,
				TotalPrice: totalPrice,
				Items:      []orderEntity.OrderItem{item},
			}
		} else {
			ordersMap[orderId].Items = append(ordersMap[orderId].Items, item)
		}
	}

	var orders []orderEntity.Order

	for _, order := range ordersMap {
		orders = append(orders, *order)
	}

	return &orders, nil
}

func runInTransaction(ctx context.Context, db *pgxpool.Pool, fn func(tx pgx.Tx) error) error {
	tx, err := db.Begin(ctx)
	if err != nil {
		return err
	}

	err = fn(tx)

	if err == nil {
		return tx.Commit(ctx)
	}

	rollbackErr := tx.Rollback(ctx)
	if rollbackErr != nil {
		return errors.Join(err, rollbackErr)
	}

	return err
}
