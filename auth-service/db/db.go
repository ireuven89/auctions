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
	"github.com/redis/go-redis/v9"
	"github.com/sethvargo/go-retry"
)

const dbConnFormat = "%s:%s@tcp(%s:%d)/"
const migrationFolder = "/migrations"

func MustNewDB(host, user, password string, port int) (*sql.DB, error) {
	dsn := fmt.Sprintf(dbConnFormat, user, password, host, port)

	db, err := sql.Open("mysql", dsn)

	if err != nil {
		return nil, fmt.Errorf("MustNewDB %w", err)
	}

	//create db
	_, err = db.Exec("create database if not exists auth")

	if err != nil {
		return nil, fmt.Errorf("MustNewDB %w", err)
	}

	//open conn with DB
	dsn = fmt.Sprintf(dbConnFormat+"auth", user, password, host, port)
	db, err = sql.Open("mysql", dsn)

	if err != nil {
		return nil, fmt.Errorf("MustNewDB %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed intiating DB %w", err)
	}

	if err = migrate(db); err != nil {
		return nil, fmt.Errorf("failed migrating DB %w", err)
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
		return fmt.Errorf("failed migrating db %w", err)
	}

	return nil
}

func MustNewRedis(host, password string) (*redis.Client, error) {
	opt := redis.Options{Addr: host, Password: password}

	c := redis.NewClient(&opt)

	if statusCmd := c.Ping(context.Background()); statusCmd.Err() != nil {
		return nil, fmt.Errorf("MustNewRedis failed %w", statusCmd.Err())
	}

	return c, nil
}
