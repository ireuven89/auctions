package db

import (
	"github.com/ireuven89/auctions/auction-service/auction"

	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"go.uber.org/zap"
)

type AuctionDB struct {
	ID       string `db:"id"`
	Name     string `db:"name"`
	BidderId string `db:"bidder_id"`
}

func toAuction(db AuctionDB) auction.Auction {

	return auction.Auction{
		ID:       db.ID,
		Name:     db.Name,
		BidderId: db.BidderId,
	}
}

type Repository interface {
	Find(ctx context.Context, id string) (auction.Auction, error)
	FindAll(ctx context.Context, request auction.AuctionRequest) ([]auction.Auction, error)
	Update(ctx context.Context, auction auction.AuctionRequest) error
	Create(ctx context.Context, auction auction.AuctionRequest) error
	Delete(ctx context.Context, id string) error
	DeleteMany(ctx context.Context, ids []interface{}) error
}

type AuctionRepository struct {
	logger *zap.Logger
	db     *sql.DB
}

func NewRepository(db *sql.DB, logger *zap.Logger) Repository {

	return &AuctionRepository{
		logger: logger,
		db:     db,
	}
}

func (r *AuctionRepository) Find(ctx context.Context, id string) (auction.Auction, error) {
	var result AuctionDB
	start := time.Now()

	q := "select id, name from auctions where id = ?"

	r.logger.Debug("AuctionRepository.Find ", zap.Any("query", q), zap.Any("args", id))

	defer func() {
		log.Printf("Query took %s: %s", time.Since(start), q)
	}()

	row := r.db.QueryRowContext(ctx, q, id)

	if row.Err() != nil {
		r.logger.Error("AuctionRepository.Find failed fetching result ", zap.Error(row.Err()), zap.String("id", id))
		return auction.Auction{}, row.Err()
	}

	if err := row.Scan(&result.ID, &result.Name); err != nil {
		r.logger.Error("failed getting db result", zap.Error(err))
		return auction.Auction{}, err
	}

	return toAuction(result), nil
}

func (r *AuctionRepository) FindAll(ctx context.Context, request auction.AuctionRequest) ([]auction.Auction, error) {
	var result []auction.Auction
	whereParams := prepareSearchQuery(request)
	q := fmt.Sprintf("SELECT id, name, bidder_id from auctions where %s", whereParams)

	r.logger.Debug("AuctionRepository.FindAll", zap.String("query", q))

	rows, err := r.db.QueryContext(ctx, q)

	if err != nil {
		r.logger.Error("AuctionRepository.FindAll failed to query", zap.Error(err))
		return nil, fmt.Errorf("AuctionRepository.FindAll failed to detch results %w", err)
	}

	for rows.Next() {
		var auctionDB AuctionDB
		if err = rows.Scan(&auctionDB.ID, &auctionDB.Name, &auctionDB.BidderId); err != nil {
			r.logger.Error("FindAll failed to cast results", zap.Error(err))
			return nil, fmt.Errorf("AuctionRepository.FindAll %w", err)
		}
		result = append(result, toAuction(auctionDB))
	}

	return result, nil
}

func prepareSearchQuery(request auction.AuctionRequest) string {
	var where strings.Builder

	if request.Name != "" {
		where.WriteString(fmt.Sprintf("name LIKE '%%%s%%'", request.Name))
	}

	return where.String()
}

func (r *AuctionRepository) Update(ctx context.Context, auction auction.AuctionRequest) error {
	q, args, err := buildUpdateQuery(auction)

	r.logger.Debug("AuctionRepository.Update", zap.String("query", q), zap.Any("args", args))

	if err != nil {
		r.logger.Error("AuctionRepository.Update failed to update query ", zap.Error(err))
		return err
	}

	_, err = r.db.ExecContext(ctx, q, args...)

	if err != nil {
		r.logger.Error("AuctionRepositoryUpdate. failed update query ", zap.Error(err))
		return fmt.Errorf("AuctionRepository.Update failed updating %w", err)
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

func (r *AuctionRepository) Create(ctx context.Context, auction auction.AuctionRequest) error {
	q := "insert into auctions (id, name, bidder_id) values(?, ?, ?)"

	r.logger.Debug("AuctionRepository.Create", zap.String("query", q), zap.Any("args", auction))

	_, err := r.db.ExecContext(ctx, q, auction.ID, auction.Name)

	if err != nil {
		r.logger.Error("AuctionRepository.Create failed to insert ", zap.Error(err))
		return fmt.Errorf("AuctionRepository.Create failed creating %w", err)
	}

	return nil
}

func (r *AuctionRepository) Delete(ctx context.Context, id string) error {
	q := "delete from auctions where id = ?"
	startTime := time.Now()

	r.logger.Debug("AuctionRepository.Delete", zap.String("query", q), zap.Any("args", id))

	_, err := r.db.ExecContext(ctx, q, id)

	if err != nil {
		return fmt.Errorf("AuctionRepository.Delete failed deleting %w", err)
	}

	endTime := time.Now()
	r.logger.Debug(fmt.Sprintf("query took %v", endTime.Sub(startTime)))

	return nil
}
func (r *AuctionRepository) DeleteMany(ctx context.Context, ids []interface{}) error {
	q, args := prepareInQuery("id", ids)
	startTime := time.Now()
	r.logger.Debug("AuctionRepository.DeleteMany", zap.String("query", q), zap.Any("args", args))

	if _, err := r.db.ExecContext(ctx, q, args...); err != nil {
		r.logger.Error("AuctionRepository.DeleteMany failed deleting", zap.Error(err))
		return fmt.Errorf("AuctionRepository.DeleteMany failed deleting %w", err)
	}

	endTime := time.Now()
	r.logger.Debug(fmt.Sprintf("query took %v", endTime.Sub(startTime)))
	return nil
}

func prepareInQuery(col string, vals []interface{}) (string, []interface{}) {
	q := "delete from auctions where %s in (%s)"

	placeholders := make([]string, len(vals))
	args := make([]interface{}, len(vals))
	for i, val := range vals {
		placeholders[i] = "?"
		args[i] = val
	}

	query := fmt.Sprintf(q, col, strings.Join(placeholders, ","))

	return query, args
}
