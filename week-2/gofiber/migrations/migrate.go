package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Get files in target dir
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("[Error]: Get working directory error: %v", err)
	}

	files, err := os.ReadDir(fmt.Sprintf("%s/migrations", pwd))
	if err != nil {
		log.Fatalf("[Error]: Read directory error: %v", err)
	}

	// Connect to postgres
	pool, err := pgxpool.New(ctx, os.Getenv("PG_URL"))
	if err != nil {
		log.Fatalf("[Error]: Connect to pg error: %v", err)
	}
	defer pool.Close()

	// Migrate read files
	for _, file := range files {
		if filepath.Ext(file.Name()) != ".sql" {
			continue
		}

		fmt.Printf("[%s]: Reading file...\n", file.Name())
		bytes, err := os.ReadFile(fmt.Sprintf("%s/migrations/%s", pwd, file.Name()))
		if err != nil {
			log.Fatalf("[Error]: Read file error: %v", err)
		}

		lines := strings.Split(string(bytes), ";\n")
		lines = lines[:len(lines)-1]

		for _, line := range lines {
			_, err := pool.Exec(ctx, line)
			if err != nil {
				log.Fatalf("[Error]: Migrate error: %v", err)
			}

		}

		fmt.Printf("[%s]: Done\n", file.Name())
	}
}
