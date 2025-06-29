package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/ireuven89/auctions/auction-service/auction"

	"go.uber.org/zap"
)

type AuctionDB struct {
	ID          string     `db:"id"`
	Name        string     `db:"name"`
	Description string     `db:"description"`
	UserId      string     `db:"user_id"`
	Active      bool       `db:"active"`
	EndTime     *time.Time `db:"end_time"`
	CreatedAt   *time.Time `db:"created_at"`
	UpdatedAt   *time.Time `db:"updated_at"`
}

func toAuction(db AuctionDB) auction.Auction {

	return auction.Auction{
		ID:        db.ID,
		Name:      db.Name,
		UserId:    db.UserId,
		Active:    db.Active,
		EndTime:   db.EndTime.Unix(),
		CreatedAt: *db.CreatedAt,
		UpdatedAt: *db.UpdatedAt,
	}
}

type Repository interface {
	Find(ctx context.Context, id string) (auction.Auction, error)
	FindAll(ctx context.Context, request auction.AuctionRequest) ([]auction.Auction, error)
	Update(ctx context.Context, auction auction.AuctionRequest) error
	WithTransactionContext(ctx context.Context, fn func(txRepo Repository) error) error
	Create(ctx context.Context, auction auction.AuctionRequest) error
	Delete(ctx context.Context, id string) error
	DeleteMany(ctx context.Context, ids []interface{}) error
}

type AuctionRepository struct {
	logger *zap.Logger
	db     *sql.DB
	tx     *sql.Tx
}

func NewRepository(db *sql.DB, logger *zap.Logger) Repository {

	return &AuctionRepository{
		logger: logger,
		db:     db,
	}
}

func (r *AuctionRepository) WithTransaction(tx *sql.Tx) Repository {
	return &AuctionRepository{
		db: r.db,
		tx: tx,
	}
}

func (r *AuctionRepository) WithTransactionContext(ctx context.Context, fn func(txRepo Repository) error) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	txRepo := r.WithTransaction(tx)

	if err = fn(txRepo); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (r *AuctionRepository) Find(ctx context.Context, id string) (auction.Auction, error) {
	var result AuctionDB
	start := time.Now()

	q := "select id, name, description, user_id, active, end_time, created_at, updated_at from auctions where id = ?"

	r.logger.Debug("AuctionRepository.Find ", zap.Any("query", q), zap.Any("args", id))

	defer func() {
		log.Printf("Query took %s: %s", time.Since(start), q)
	}()

	row := r.db.QueryRowContext(ctx, q, id)

	if row.Err() != nil {
		r.logger.Error("AuctionRepository.Find failed fetching result ", zap.Error(row.Err()), zap.String("id", id))
		return auction.Auction{}, row.Err()
	}

	if err := row.Scan(&result.ID, &result.Name, &result.Description, &result.UserId, &result.Active, &result.EndTime, &result.CreatedAt, &result.UpdatedAt); err != nil {
		r.logger.Error("failed getting db result", zap.Error(err))
		return auction.Auction{}, err
	}

	return toAuction(result), nil
}

func (r *AuctionRepository) FindAll(ctx context.Context, request auction.AuctionRequest) ([]auction.Auction, error) {

	whereParams := prepareSearchQuery(request)
	q := fmt.Sprintf("SELECT id, name, descrption, user_id, active, end_time, created_at, updated_at from auctions where %s", whereParams)

	r.logger.Debug("AuctionRepository.FindAll", zap.String("query", q))

	rows, err := r.db.QueryContext(ctx, q)

	if err != nil {
		r.logger.Error("AuctionRepository.FindAll failed to query", zap.Error(err))
		return nil, fmt.Errorf("AuctionRepository.FindAll failed to detch results %w", err)
	}

	result := make([]auction.Auction, 0)
	for rows.Next() {
		var auctionDB AuctionDB
		var endTime time.Time
		if err = rows.Scan(&auctionDB.ID, &auctionDB.Name, &auctionDB.Description, &auctionDB.UserId, &auctionDB.Active, &endTime, &auctionDB.CreatedAt, &auctionDB.UpdatedAt); err != nil {
			r.logger.Error("FindAll failed to cast results", zap.Error(err))
			return nil, fmt.Errorf("AuctionRepository.FindAll %w", err)
		}
		auctionDB.EndTime = &endTime
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

	if auction.Description != "" {
		sets = append(sets, "description = ?")
		args = append(args, auction.Description)
	}

	if auction.EndTime != 0 {
		sets = append(sets, "end_time = ?")
		args = append(args, time.UnixMilli(auction.EndTime))
	}

	if auction.Active != nil {
		sets = append(sets, "active = ?")
		args = append(args, *auction.Active)
	}

	if len(sets) == 0 {
		return "", nil, fmt.Errorf("AuctionRepository.Update failed to update query")
	}

	query += strings.Join(sets, ", ") + " WHERE id = ?"
	args = append(args, auction.ID)

	return query, args, nil
}

func (r *AuctionRepository) Create(ctx context.Context, auction auction.AuctionRequest) error {
	q := "insert into auctions (id, name,user_id, active, description, end_time, created_at, updated_at) values(?, ?, ?, ?, ?, ?, ?, ?)"

	r.logger.Debug("AuctionRepository.Create", zap.String("query", q), zap.Any("args", auction))

	var endTimeStamp interface{}
	if auction.EndTime == 0 {
		endTimeStamp = nil
	} else {
		endTimeStamp = time.UnixMilli(auction.EndTime)
	}

	_, err := r.db.ExecContext(ctx, q, auction.ID, auction.Name, auction.UserId, auction.Active, auction.Description, endTimeStamp, time.UnixMilli(auction.CreatedAt), time.UnixMilli(auction.UpdatedAt))

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
