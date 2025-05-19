package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pressly/goose/v3"
	"github.com/sethvargo/go-retry"
)

var migrationFolder = "/migrations"
var databaseName = "bidders"

func MustNewDB(host, user, password string, port int) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/", user, password, host, port)
	fmt.Printf("dsn is %s\n", dsn)

	backoff := retry.WithMaxRetries(5, retry.NewExponential(3*time.Second))
	var attempt int

	db, err := retry.DoValue[*sql.DB](
		context.Background(),
		backoff,
		func(ctx context.Context) (*sql.DB, error) {
			db, err := sql.Open("mysql", dsn)
			if err != nil {
				fmt.Printf("open failed: %v attempt %d\n", err, attempt)
				attempt++
				return nil, retry.RetryableError(fmt.Errorf("open failed: %w attempt %d\n", err, attempt))
			}

			if err := db.PingContext(ctx); err != nil {
				fmt.Printf("open failed: %v attemped %d\n", err, attempt)
				attempt++
				return nil, retry.RetryableError(fmt.Errorf("ping failed: %w attempt %d", err, attempt))
			}

			_, err = db.Exec(fmt.Sprintf("create database if not exists %s", databaseName))

			if err != nil {
				return nil, err
			}
			//this connects to after db creation
			db, err = sql.Open("mysql", dsn+databaseName)

			if err != nil {
				return nil, retry.RetryableError(err)
			}

			return db, nil
		},
	)

	if err != nil {
		return nil, err
	}

	if err = migrate(db); err != nil {
		return nil, err
	}

	return db, nil
}

func migrate(db *sql.DB) error {

	if err := goose.SetDialect("mysql"); err != nil {
		return err
	}
	goose.SetLogger(log.New(os.Stdout, "goose: ", log.LstdFlags))

	ctx := context.Background()

	backoff := retry.WithMaxRetries(3, retry.NewConstant(1*time.Second))
	var attempt int

	migrationPath := os.Getenv("MIGRATIONS_DIR")

	if migrationPath == "" {
		migrationPath = migrationFolder
	}

	err := retry.Do(ctx, backoff, func(ctx context.Context) error {

		err := goose.Up(db, migrationPath)

		if err != nil {
			fmt.Printf("failed runnning migration %v attempt %d\n", err, attempt)
			attempt++
			return retry.RetryableError(err)
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
