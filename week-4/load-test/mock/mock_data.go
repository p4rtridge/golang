package main

import (
	"context"
	"fmt"
	"load_test/pkg"
	"log"
	"os"
	"runtime"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	USER_NUM         = 1000
	USER_BALANCE     = float32(50.0)
	PRODUCT_NAME     = "your mom"
	PRODUCT_PRICE    = float32(25.0)
	PRODUCT_QUANTITY = 100
)

const (
	QUERY_ADD_USERS   = "INSERT INTO users (username, password, balance) VALUES ($1, $2, $3)"
	QUERY_ADD_PRODUCT = "INSERT INTO products (name, quantity, price) VALUES ($1, $2, $3)"
)

var ctx context.Context = context.Background()

func main() {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	pool, err := pgxpool.New(ctx, os.Getenv("PG_URL"))
	if err != nil {
		log.Fatal(err)
	}

	hasher := pkg.NewHasher(64*1024, 3, 16, 32, uint8(runtime.NumCPU()))

	for i := 0; i < USER_NUM; i++ {
		password, err := hasher.HashPassword("super_safe_password")
		if err != nil {
			log.Fatal(err)
		}

		result, err := pool.Exec(ctx, QUERY_ADD_USERS, fmt.Sprintf("user-%d", i), password, USER_BALANCE)
		if err != nil {
			fmt.Errorf("[%d] err: %v", i, err)
		}
		if result.RowsAffected() < 1 {
			fmt.Errorf("[%d] row cannot be insert", i)
		}

		fmt.Printf("[%d] done", i)
	}

	_, err = pool.Exec(ctx, QUERY_ADD_PRODUCT, PRODUCT_NAME, PRODUCT_QUANTITY, PRODUCT_PRICE)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("done")
}
