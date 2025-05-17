package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/ireuven89/auctions/bidder-service/bidder"
	"go.uber.org/zap"
	"strings"
)

type Repository interface {
	Find(ctx context.Context, id string) (bidder.Bidder, error)
	FindAll(ctx context.Context, request bidder.BiddersRequest) ([]bidder.Bidder, error)
	Create(ctx context.Context, bidder bidder.Bidder) error
	Update(ctx context.Context, bidder bidder.Bidder) error
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

func (r *Repo) Find(ctx context.Context, id string) (bidder.Bidder, error) {
	var result BidderDb
	q := "select id, name from bidders where id = ?"

	row := r.db.QueryRowContext(ctx, q, id)

	if row.Err() != nil {
		return bidder.Bidder{}, row.Err()
	}

	if err := row.Scan(&result.ID, &result.Name); err != nil {
		return bidder.Bidder{}, err
	}

	return toBidder(result), nil
}

func (r *Repo) Create(ctx context.Context, bidder bidder.Bidder) error {
	q := "insert into bidders (id, name) values(?, ?)"

	if _, err := r.db.ExecContext(ctx, q, bidder.ID, bidder.Name); err != nil {
		r.logger.Error("Create failed inserting to ", zap.Error(err))
		return err
	}
	return nil
}

func (r *Repo) Update(ctx context.Context, bidder bidder.Bidder) error {
	q := "update bidders set name = ? where id = ?"

	if _, err := r.db.ExecContext(ctx, q, bidder.Name, bidder.ID); err != nil {
		r.logger.Error("Update ailed updating", zap.Error(err), zap.String("resource", bidder.ID))
		return fmt.Errorf("Repository.Update %w", err)
	}

	return nil
}

func (r *Repo) Delete(ctx context.Context, id string) error {
	q := "delete from bidders-service where id = ?"

	if _, err := r.db.ExecContext(ctx, q, id); err != nil {
		r.logger.Error("failed deleting bidder ", zap.Error(err), zap.String("id", id))
		return err
	}

	return nil
}

func (r *Repo) FindAll(ctx context.Context, request bidder.BiddersRequest) ([]bidder.Bidder, error) {
	whereParams := prepareWhereQuery(request)
	q := fmt.Sprintf("SELECT id, name from bidders where %s", whereParams)

	rows, err := r.db.QueryContext(ctx, q)

	if err != nil {
		r.logger.Error("failed fetching bidders ", zap.Error(err))
		return nil, fmt.Errorf("FindAll failed fetcing reuqest %w", err)
	}

	biddersDB, err := parseResults(rows)

	if err != nil {
		return nil, fmt.Errorf("failed parsing results %w", err)
	}

	var result []bidder.Bidder

	for _, bidderDB := range biddersDB {
		result = append(result, toBidder(bidderDB))
	}

	return result, nil
}

func parseResults(rows *sql.Rows) ([]BidderDb, error) {
	var res []BidderDb

	for rows.Next() {
		var bidder BidderDb
		if err := rows.Scan(&bidder.ID, &bidder.Name); err != nil {
			return nil, err
		}
		res = append(res, bidder)
	}
	return res, nil
}

func prepareWhereQuery(request bidder.BiddersRequest) string {
	var where strings.Builder

	if request.Name != "" {
		where.WriteString(fmt.Sprintf("name LIKE '%%%s%%'", request.Name))
	}

	return where.String()
}
