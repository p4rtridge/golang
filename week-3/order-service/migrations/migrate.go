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

		for _, statement := range parseSQLStatement(strings.Split(string(bytes), "\n")) {
			_, err := pool.Exec(ctx, statement)
			if err != nil {
				log.Fatalf("[Error]: Error while executing statement: %v", err)
			}
		}

		fmt.Printf("[%s]: Done\n", file.Name())
	}
}

func parseSQLStatement(lines []string) []string {
	isFunc := false
	statements := make([]string, 0)
	var statement strings.Builder

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine == "" {
			continue
		}

		if strings.HasPrefix(strings.ToUpper(trimmedLine), "CREATE OR REPLACE FUNCTION") {
			isFunc = true
		}

		statement.WriteString(line + "\n")

		if isFunc && strings.Contains(strings.ToUpper(trimmedLine), "LANGUAGE") {
			statements = append(statements, statement.String())
			statement.Reset()
			isFunc = false
			continue
		}

		if !isFunc && strings.HasSuffix(trimmedLine, ";") {
			statements = append(statements, statement.String())
			statement.Reset()
		}
	}

	return statements
}

func normalizeStr(s string) string {
	return strings.ToUpper(strings.TrimSpace(s))
}
