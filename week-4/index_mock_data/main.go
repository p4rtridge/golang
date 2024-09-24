package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"os"
	"sync"

	"github.com/go-faker/faker/v4"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Product struct {
	Uuid     string  `faker:"uuid_digit"`
	Name     string  `faker:"word"`
	Quantity float64 `faker:"amount"`
	Price    int     `faker:"boundary_start=10, boundary_end=101"`
}

const (
	COUNT    = 5_000_000
	WP_COUNT = 200
)

func main() {
	ctx := context.Background()

	cfg, err := pgxpool.ParseConfig(os.Getenv("PG_URL"))
	if err != nil {
		log.Fatalln(err)
	}

	cfg.MaxConns = 100

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		log.Fatalln(err)
	}

	wg := wp(pool)

	wg.Wait()
	fmt.Println("Done")
}

func wp(pool *pgxpool.Pool) *sync.WaitGroup {
	var wg sync.WaitGroup

	wg.Add(WP_COUNT)
	for i := 0; i < WP_COUNT; i++ {
		go process(i, pool, &wg)
	}

	return &wg
}

func process(workerId int, pool *pgxpool.Pool, wg *sync.WaitGroup) {
	ctx, cancel := context.WithCancel(context.Background())
	defer wg.Done()
	defer cancel()

	for i := 0; i < int(math.Ceil(COUNT/WP_COUNT)); i++ {
		var product Product
		err := faker.FakeData(&product)
		if err != nil {
			log.Fatalln(err)
		}
		_, err = pool.Exec(ctx, "INSERT INTO products (name, quantity, price) VALUES ($1, $2, $3)", fmt.Sprintf("%s-%s", product.Name, product.Uuid), int(product.Quantity), float32(product.Price))
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Printf("[Worker %d] inserted one record.\n", workerId)
	}
}
