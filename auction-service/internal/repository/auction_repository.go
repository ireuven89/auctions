package repository

import (
	"github.com/ireuven89/auctions/auction-service/domain"
	"github.com/ireuven89/auctions/auction-service/internal/service"

	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"go.uber.org/zap"
)

type AuctionDB struct {
	ID           string    `db:"id"`
	Description  string    `db:"description"`
	Regions      []byte    `db:"regions"`
	InitialOffer float64   `db:"initial_offer"`
	CurrentBid   float64   `db:"current_highest"`
	Status       string    `db:"status"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

func toAuction(db AuctionDB) domain.Auction {

	return domain.Auction{
		ID:           db.ID,
		Description:  db.Description,
		Regions:      db.Regions,
		InitialOffer: db.InitialOffer,
		CurrentBid:   db.CurrentBid,
		Status:       domain.FromString(db.Status),
		CreatedAt:    db.CreatedAt,
		UpdatedAt:    db.UpdatedAt,
	}
}

type AuctionRepository struct {
	logger *zap.Logger
	db     *sql.DB
}

func NewAuctionRepo(db *sql.DB, logger *zap.Logger) service.Repository {

	return &AuctionRepository{
		logger: logger,
		db:     db,
	}
}

func (r *AuctionRepository) Find(ctx context.Context, id string) (domain.Auction, error) {
	var result AuctionDB
	start := time.Now()

	q := "select id, description from auctions where id = ?"

	r.logger.Debug("AuctionRepository.Find ", zap.Any("query", q), zap.Any("args", id))

	defer func() {
		log.Printf("Query took %s: %s", time.Since(start), q)
	}()

	row := r.db.QueryRowContext(ctx, q, id)

	if row.Err() != nil {
		r.logger.Error("AuctionRepository.Find failed fetching result ", zap.Error(row.Err()), zap.String("id", id))
		return domain.Auction{}, row.Err()
	}

	if err := row.Scan(&result.ID, &result.Description); err != nil {
		r.logger.Error("failed getting db result", zap.Error(err))
		return domain.Auction{}, err
	}

	return toAuction(result), nil
}

func (r *AuctionRepository) FindAll(ctx context.Context, request domain.AuctionRequest) ([]domain.Auction, error) {
	var result []domain.Auction
	whereParams := prepareSearchQuery(request)
	q := fmt.Sprintf("SELECT id, description, regions, status, initalOffer, created_at, updatead_at from auctions where %s", whereParams)

	r.logger.Debug("AuctionRepository.FindAll", zap.String("query", q))

	rows, err := r.db.QueryContext(ctx, q)

	if err != nil {
		r.logger.Error("AuctionRepository.FindAll failed to query", zap.Error(err))
		return nil, fmt.Errorf("AuctionRepository.FindAll failed to detch results %w", err)
	}

	for rows.Next() {
		var auctionDB AuctionDB
		if err = rows.Scan(&auctionDB.ID, &auctionDB.Description, &auctionDB.Regions, &auctionDB.Status, &auctionDB.InitialOffer, &auctionDB.CreatedAt, &auctionDB.UpdatedAt); err != nil {
			r.logger.Error("FindAll failed to cast results", zap.Error(err))
			return nil, fmt.Errorf("AuctionRepository.FindAll %w", err)
		}
		result = append(result, toAuction(auctionDB))
	}

	return result, nil
}

func prepareSearchQuery(request domain.AuctionRequest) string {
	var where strings.Builder

	if request.Description != "" {
		where.WriteString(fmt.Sprintf("description LIKE '%%%s%%'", request.Description))
	}

	return where.String()
}

func (r *AuctionRepository) Update(ctx context.Context, auction domain.AuctionRequest) error {
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

func buildUpdateQuery(auction domain.AuctionRequest) (string, []interface{}, error) {
	query := "UPDATE auctions SET "
	var sets []string
	var args []interface{}

	if auction.Description != "" {
		sets = append(sets, "description = ?")
		args = append(args, auction.Description)
	}

	if len(auction.Regions) != 0 {
		args = append(args, auction.Regions)
	}

	if auction.Status != "" {
		args = append(args, auction.Status)
	}

	if len(sets) == 0 {
		return "", nil, fmt.Errorf("no fields to update")
	}

	query += strings.Join(sets, ", ") + " WHERE id = ?"
	args = append(args, auction.ID)

	return query, args, nil
}

func (r *AuctionRepository) Create(ctx context.Context, auction domain.AuctionRequest) error {
	q := "insert into auctions (id, description, seller_id, regions, status, initial_offer, created_at, updated_at) values(?, ?, ?, ?, ?, ?, ?, ?)"

	r.logger.Debug("AuctionRepository.Create", zap.String("query", q), zap.Any("args", auction))

	_, err := r.db.ExecContext(ctx, q, auction.ID, auction.Description, auction.SellerId, auction.Regions, auction.Status, auction.InitialOffer, auction.CreatedAt, auction.UpdatedAt)

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
