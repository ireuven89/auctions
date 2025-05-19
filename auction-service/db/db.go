package db

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pressly/goose/v3"
	"github.com/sethvargo/go-retry"
	"os"
	"time"
)

var migrationFolder = "/migrations"
var databaseName = "auctions"

func MustNewDB(host, user, password string, port int) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/", user, password, host, port)
	println("Test blocked push")

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
			return nil, err
		}

		return db, nil
	})

	if err = migrate(db); err != nil {
		return nil, err
	}

	return db, nil
}

func migrate(db *sql.DB) error {

	if err := goose.SetDialect("mysql"); err != nil {
		return err
	}

	ctx := context.Background()

	migrationPath := os.Getenv("MIGRATIONS_DIR")

	if migrationPath == "" {
		migrationPath = migrationFolder
	}

	backoff := retry.WithMaxRetries(3, retry.NewConstant(1*time.Second))
	err := retry.Do(ctx, backoff, func(ctx context.Context) error {
		err := goose.Up(db, migrationPath)

		if err != nil {
			return retry.RetryableError(err)
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
