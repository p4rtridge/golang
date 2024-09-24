package main

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx := context.Background()

	pg, err := pgxpool.New(ctx, os.Getenv("PG_URL"))
	if err != nil {
		log.Fatalln(err)
	}

	userId := 0
	balance := float32(0.0)
	err = pg.QueryRow(ctx, "SELECT id, balance FROM users where id = $1", 1001).Scan(&userId, &balance)
	if err != nil {
		log.Fatalln(err)
	}

	productId := 0
	productName := ""
	productQuantity := 0
	err = pg.QueryRow(ctx, "SELECT id, name, quantity FROM products where id = 5003562").Scan(&productId, &productName, &productQuantity)
	if err != nil {
		log.Fatalln(err)
	}

	// simulate order validation
	if true {
		orderId := 0

		err = pg.QueryRow(ctx, "INSERT INTO orders (user_id, total_price) VALUES ($1, $2) RETURNING id", userId, 100).Scan(&orderId)
		if err != nil {
			log.Fatalln(err)
		}

		// error
		_, err = pg.Exec(ctx, "INSERT INTO order_items (order_id, product_id, product_name, product_quantity, product_price) VALUES ($1, $2, $3, $4)", orderId, productId, productName, 100)
		if err != nil {
			log.Fatalln(err)
		}
	}
}
