package postgres

import (
	"context"
	"fmt"
	"order_service/pkg"
	orderEntity "order_service/services/order/entity"
	productEntity "order_service/services/product/entity"
	userEntity "order_service/services/user/entity"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, order *orderEntity.Order, callbackFn func(order *orderEntity.Order, user *userEntity.User, products *[]productEntity.Product) (bool, error)) error
	GetOrders(ctx context.Context, userId int) (*[]orderEntity.Order, error)
}

const (
	QUERY_GET_ORDERS                  = "SELECT o.id AS order_id, o.user_id, oi.product_id, oi.product_name, oi.product_price, oi.quantity, o.total_price, o.created_at, o.updated_at FROM orders AS o JOIN order_items AS oi ON o.id = oi.order_id WHERE o.user_id = $1"
	QUERY_GET_USER_LOCK               = "SELECT * FROM users WHERE id = $1"
	QUERY_GET_PRODUCT_LOCK            = "SELECT * FROM products WHERE id = $1"
	QUERY_CREATE_ORDER_WITH_RETURN_ID = "INSERT INTO orders (user_id, total_price) VALUES ($1, $2) RETURNING id"
	QUERY_CREATE_ORDER_ITEM           = "INSERT INTO order_items (order_id, product_id, product_name, product_price, quantity) VALUES ($1, $2, $3, $4, $5)"
	QUERY_UPDATE_USER_BALANCE         = "UPDATE users SET balance = balance - $2, updated_at = $3 WHERE id = $1"
	QUERY_UPDATE_PRODUCT_QUANTITY     = "UPDATE products SET quantity = quantity - $2, updated_at = $3 WHERE id = $1"
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
	return pkg.RunInTransaction(ctx, repo.db, func(tx pgx.Tx) error {
		var user userEntity.User

		err := tx.QueryRow(ctx, QUERY_GET_USER_LOCK, order.GetUserIdSafe()).Scan(&user.Id, &user.Username, &user.Password, &user.Balance, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return err
		}

		// fetch required datas
		orderItems := order.GetItemsSafe()
		products := make([]productEntity.Product, 0, len(orderItems))
		for _, item := range orderItems {
			var product productEntity.Product

			err := tx.QueryRow(ctx, QUERY_GET_PRODUCT_LOCK, item.GetProductId()).Scan(&product.Id, &product.Name, &product.Quantity, &product.Price, &product.CreatedAt, &product.UpdatedAt)
			if err != nil {
				return err
			}
			fmt.Println(product.GetQuantity())

			products = append(products, product)
		}

		// run business logic
		accept, err := callbackFn(order, &user, &products)
		if err != nil {
			return err
		}
		if !accept {
			return nil
		}

		var newOrderId int

		err = tx.QueryRow(ctx, QUERY_CREATE_ORDER_WITH_RETURN_ID, order.GetUserIdSafe(), order.GetTotalPriceSafe()).Scan(&newOrderId)
		if err != nil {
			return err
		}
		order.SetId(newOrderId)

		// handle product's stock and user's balance after ordered
		orderItems = order.GetItemsSafe()
		if len(orderItems) == 0 {
			return orderEntity.ErrInvalidMemory
		}

		for idx, item := range orderItems {
			_, err = tx.Exec(ctx, QUERY_CREATE_ORDER_ITEM, order.GetIdSafe(), item.GetProductId(), item.GetProductName(), item.GetProductPrice(), item.GetQuantity())
			if err != nil {
				return err
			}

			_, err = tx.Exec(ctx, QUERY_UPDATE_PRODUCT_QUANTITY, products[idx].GetId(), products[idx].GetQuantity(), time.Now())
			if err != nil {
				return err
			}
		}

		_, err = tx.Exec(ctx, QUERY_UPDATE_USER_BALANCE, user.GetId(), user.GetBalance(), time.Now())
		if err != nil {
			return err
		}

		return nil
	})
}

func (repo *postgresRepo) GetOrders(ctx context.Context, userId int) (*[]orderEntity.Order, error) {
	rows, err := repo.db.Query(ctx, QUERY_GET_ORDERS, userId)
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
