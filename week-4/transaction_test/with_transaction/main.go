package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx := context.Background()

	var e error
	pg, err := pgxpool.New(ctx, os.Getenv("PG_URL"))
	if err != nil {
		log.Fatalln(err)
	}

	tx, err := pg.Begin(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	defer func() {
		if e != nil {
			tx.Commit(ctx)
		} else {
			tx.Rollback(ctx)
		}

		time.Sleep(30 * time.Second)
	}()

	userId := 0
	balance := float32(0.0)
	err = tx.QueryRow(ctx, "SELECT id, balance FROM users where id = $1", 1001).Scan(&userId, &balance)
	if err != nil {
		e = err
		log.Fatalln(err)
	}

	productId := 0
	productName := ""
	productQuantity := 0
	err = tx.QueryRow(ctx, "SELECT id, name, quantity FROM products where id = 5003562").Scan(&productId, &productName, &productQuantity)
	if err != nil {
		e = err
		log.Fatalln(err)
	}

	// simulate order validation
	if true {
		orderId := 0

		err = tx.QueryRow(ctx, "INSERT INTO orders (user_id, total_price) VALUES ($1, $2) RETURNING id", userId, 100).Scan(&orderId)
		if err != nil {
			e = err
			log.Fatalln(err)
		}

		_, err = tx.Exec(ctx, "INSERT INTO order_items (order_id, product_id, product_name, product_price) VALUES ($1, $2, $3, $4)", orderId, productId, productName, 100)
		if err != nil {
			e = err
			log.Fatalln(err)
		}
	}
}
