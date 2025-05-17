package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/ireuven89/auctions/auction-service/auction"
	"go.uber.org/zap"
	"log"
	"strings"
	"time"
)

type AuctionDB struct {
	ID       string `db:"id"`
	Name     string `db:"name"`
	BidderId string `db:"bidder_id"`
}

func toAuction(db AuctionDB) auction.Auction {

	return auction.Auction{
		ID:   db.ID,
		Name: db.Name,
	}
}

type Repository interface {
	Find(ctx context.Context, id string) (auction.Auction, error)
	Update(ctx context.Context, auction auction.AuctionRequest) error
	Create(ctx context.Context, auction auction.AuctionRequest) error
	Delete(ctx context.Context, id string) error
}

type Repo struct {
	logger *zap.Logger
	db     *sql.DB
}

func NewRepository(db *sql.DB, logger *zap.Logger) Repository {

	return &Repo{
		logger: logger,
		db:     db,
	}
}

func (r *Repo) Find(ctx context.Context, id string) (auction.Auction, error) {
	var result AuctionDB
	start := time.Now()

	q := "select id, name from auctions where id = ?"

	r.logger.Debug(q)

	defer func() {
		log.Printf("Query took %s: %s", time.Since(start), q)
	}()

	row := r.db.QueryRowContext(ctx, q, id)

	if row.Err() != nil {
		r.logger.Error("failed fetching result ", zap.Error(row.Err()), zap.String("id", id))
		return auction.Auction{}, row.Err()
	}

	if err := row.Scan(&result.ID, &result.Name); err != nil {
		r.logger.Error("failed getting db result", zap.Error(err))
		return auction.Auction{}, err
	}

	return toAuction(result), nil
}

func (r *Repo) Update(ctx context.Context, auction auction.AuctionRequest) error {
	q, args, err := buildUpdateQuery(auction)

	r.logger.Debug(q)

	if err != nil {
		r.logger.Error("Update. failed to update query ", zap.Error(err))
		return err
	}

	_, err = r.db.ExecContext(ctx, q, args...)

	if err != nil {
		r.logger.Error("Update. failed update query ", zap.Error(err))
		return err
	}

	return nil
}

func buildUpdateQuery(auction auction.AuctionRequest) (string, []interface{}, error) {
	query := "UPDATE auctions SET "
	var sets []string
	var args []interface{}

	if auction.Name != "" {
		sets = append(sets, "name = ?")
		args = append(args, auction.Name)
	}

	if len(sets) == 0 {
		return "", nil, fmt.Errorf("no fields to update")
	}

	query += strings.Join(sets, ", ") + " WHERE id = ?"
	args = append(args, auction.ID)

	return query, args, nil
}

func (r *Repo) Create(ctx context.Context, auction auction.AuctionRequest) error {
	q := "insert into auctions (id, name) values(?, ?)"

	_, err := r.db.ExecContext(ctx, q, auction.ID, auction.Name)

	if err != nil {
		r.logger.Error("Create.failed to insert ", zap.Error(err))
		return err
	}

	return nil
}

func (r *Repo) Delete(ctx context.Context, id string) error {
	q := "delete from auctions where id = ?"

	_, err := r.db.ExecContext(ctx, q, id)

	if err != nil {
		return err
	}

	return nil
}
