package postgres

import (
	"context"
	orderEntity "order_service/services/order/entity"
	productEntity "order_service/services/product/entity"
	userEntity "order_service/services/user/entity"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, data *orderEntity.Order, userData *userEntity.User, productData *[]productEntity.Product) error
	GetOrders(ctx context.Context) (*[]orderEntity.Order, error)
}

const (
	QUERY_INSERT_ORDER            = "INSERT INTO orders (user_id, total_price) VALUES ($1, $2) RETURNING id"
	QUERY_INSERT_ORDER_ITEM       = "INSERT INTO order_items (order_id, product_id, product_name, product_price, quantity) VALUES ($1, $2, $3, $4, $5)"
	QUERY_UPDATE_USER_BALANCE     = "UPDATE users SET balance = COALESCE($2, balance) WHERE id = $1"
	QUERY_UPDATE_PRODUCT_QUANTITY = "UPDATE products SET quantity = COALESCE($2, quantity) WHERE id = $1"
	QUERY_GET_ORDERS              = "SELECT o.id AS order_id, o.user_id, o.status, oi.product_id, oi.product_name, oi.product_price, oi.quantity, o.total_price, o.created_at, o.updated_at FROM orders AS o JOIN order_items AS oi ON o.id = oi.order_id"
	QUERY_UPDATE_ORDER_STATUS     = "UPDATE orders SET status = $2 WHERE order_id = $1"
)

type postgresRepo struct {
	db *pgxpool.Pool
}

func NewOrderRepo(db *pgxpool.Pool) OrderRepository {
	return &postgresRepo{
		db,
	}
}

func (repo *postgresRepo) CreateOrder(ctx context.Context, data *orderEntity.Order, userData *userEntity.User, productData *[]productEntity.Product) error {
	tx, err := repo.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
			_, err = repo.db.Exec(ctx, QUERY_UPDATE_ORDER_STATUS, orderEntity.CANCELED)
		} else {
			tx.Commit(ctx)
		}
	}()

	var newOrderId int
	err = tx.QueryRow(ctx, QUERY_INSERT_ORDER, userData.Id, data.TotalPrice).Scan(&newOrderId)
	if err != nil {
		return err
	}

	for idx, item := range data.Items {
		_, err = tx.Exec(ctx, QUERY_INSERT_ORDER_ITEM, newOrderId, item.ProductId, item.ProductName, item.ProductPrice, item.Quantity)
		if err != nil {
			return err
		}

		_, err = tx.Exec(ctx, QUERY_UPDATE_PRODUCT_QUANTITY, (*productData)[idx].Id, (*productData)[idx].Quantity)
		if err != nil {
			return err
		}
	}

	_, err = tx.Exec(ctx, QUERY_UPDATE_USER_BALANCE, userData.Id, userData.Balance)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, QUERY_UPDATE_ORDER_STATUS, newOrderId, orderEntity.DONE)
	if err != nil {
		return err
	}

	return nil
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
		var status string

		err := rows.Scan(&orderId, &userId, &status, &productId, &productName, &productPrice, &quantity, &totalPrice, &createdAt, &updatedAt)
		if err != nil {
			return nil, err
		}

		item := orderEntity.NewOrderItem(orderId, productId, productName, productPrice, quantity)

		if _, exists := ordersMap[orderId]; !exists {
			ordersMap[orderId] = &orderEntity.Order{
				Id:         orderId,
				UserId:     userId,
				TotalPrice: totalPrice,
				Status:     orderEntity.OrderStatus(status),
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
