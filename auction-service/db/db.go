package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pressly/goose/v3"
	"github.com/sethvargo/go-retry"
)

var migrationFolder = "/migrations"
var databaseName = "auctions"
var gooseUp = goose.Up

func MustNewDB(host, user, password string, port int) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/", user, password, host, port)

	backoff := retry.WithMaxRetries(5, retry.NewExponential(3*time.Second))
	var attempt int

	db, err := retry.DoValue[*sql.DB](context.Background(), backoff, func(ctx context.Context) (*sql.DB, error) {
		db, err := sql.Open("mysql", dsn)

		if err != nil {
			fmt.Printf("failed connecting db %v attempt %d", err, attempt)
			attempt++
			return nil, retry.RetryableError(err)
		}

		if err = db.Ping(); err != nil {
			fmt.Printf("failed pinggin db %v attempt %d", err, attempt)
			attempt++
			return nil, retry.RetryableError(err)
		}

		_, err = db.Exec(fmt.Sprintf("create database if not exists %s ", databaseName))

		if err != nil {
			fmt.Printf("failed creating db %v attempt %d", err, attempt)
			attempt++
			return nil, err
		}

		db, err = sql.Open("mysql", dsn+databaseName)
		if err != nil {
			fmt.Printf("failed connecting db %v with databse name attempt %d", err, attempt)
			attempt++
			return nil, fmt.Errorf("failed opening db %w", err)
		}

		return db, nil
	})

	if err = migrate(db); err != nil {
		return nil, fmt.Errorf("failed migrating %w", err)
	}

	return db, nil
}

func migrate(db *sql.DB) error {

	if err := goose.SetDialect("mysql"); err != nil {
		fmt.Printf("migrate %v", err)
		return err
	}

	ctx := context.Background()

	migrationPath := os.Getenv("MIGRATIONS_DIR")

	if migrationPath == "" {
		migrationPath = migrationFolder
	}

	backoff := retry.WithMaxRetries(3, retry.NewConstant(1*time.Second))
	attemtps := 1
	err := retry.Do(ctx, backoff, func(ctx context.Context) error {
		err := gooseUp(db, migrationPath)

		if err != nil {
			fmt.Printf("migrate failed %v", err)
			return retry.RetryableError(err)
		}

		return nil
	})

	if err != nil {
		fmt.Errorf("migrate failed %v retry number %d", err, attemtps)
		return fmt.Errorf("migrate %w", err)
	}

	return nil
}
