package db

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/ireuven89/auctions/bidder-service/bidder"
	"go.uber.org/zap"
)

type Repository interface {
	Find(ctx context.Context, id string) (bidder.Bidder, error)
	FindAll(ctx context.Context, request bidder.BiddersRequest) ([]bidder.Bidder, error)
	Create(ctx context.Context, bidder bidder.Bidder) error
	Update(ctx context.Context, bidder bidder.Bidder) error
	Delete(ctx context.Context, id string) error
	DeleteMany(ctx context.Context, ids []interface{}) error
}

type BidderRepository struct {
	logger *zap.Logger
	db     *sql.DB
}

func NewRepository(db *sql.DB, logger *zap.Logger) Repository {

	return &BidderRepository{
		logger: logger,
		db:     db,
	}
}

func (r *BidderRepository) Find(ctx context.Context, id string) (bidder.Bidder, error) {
	var result BidderDb
	q := "select id, name, item from bidders where id = ?"

	row := r.db.QueryRowContext(ctx, q, id)

	if row.Err() != nil {
		return bidder.Bidder{}, row.Err()
	}

	if err := row.Scan(&result.ID, &result.Name, &result.Item); err != nil {
		return bidder.Bidder{}, err
	}

	return toBidder(result), nil
}

func (r *BidderRepository) Create(ctx context.Context, bidder bidder.Bidder) error {
	q := "insert into bidders (id, name, item) values(?, ?, ?)"

	if _, err := r.db.ExecContext(ctx, q, bidder.ID, bidder.Name, bidder.Item); err != nil {
		r.logger.Error("Create failed inserting to ", zap.Error(err))
		return err
	}
	return nil
}

func (r *BidderRepository) Update(ctx context.Context, bidder bidder.Bidder) error {
	q, args, err := prepareUpdateQuery(bidder)

	if err != nil {
		return fmt.Errorf("BidderRepository.Update failed preparing query %w", err)
	}

	if _, err = r.db.ExecContext(ctx, q, args...); err != nil {
		r.logger.Error("Update failed updating", zap.Error(err), zap.String("resource", bidder.ID))
		return fmt.Errorf("Repository.Update %w", err)
	}

	return nil
}

func prepareUpdateQuery(bidder bidder.Bidder) (string, []interface{}, error) {
	query := "UPDATE bidders SET "
	var sets []string
	var args []interface{}

	if bidder.Name != "" {
		sets = append(sets, "name = ?")
		args = append(args, bidder.Name)
	}

	if bidder.Item != "" {
		sets = append(sets, "item = ?")
		args = append(args, bidder.Item)
	}

	if len(sets) == 0 {
		return "", nil, fmt.Errorf("no fields to update")
	}

	query += strings.Join(sets, ", ") + " WHERE id = ?"
	args = append(args, bidder.ID)

	return query, args, nil
}

func (r *BidderRepository) Delete(ctx context.Context, id string) error {
	q := "delete from bidders-service where id = ?"

	if _, err := r.db.ExecContext(ctx, q, id); err != nil {
		r.logger.Error("failed deleting bidder ", zap.Error(err), zap.String("id", id))
		return err
	}

	return nil
}

func (r *BidderRepository) FindAll(ctx context.Context, request bidder.BiddersRequest) ([]bidder.Bidder, error) {
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

func (r *BidderRepository) DeleteMany(ctx context.Context, ids []interface{}) error {
	q, args := prepareInQuery("id", ids)

	if _, err := r.db.ExecContext(ctx, q, args...); err != nil {
		r.logger.Error("BidderRepository.DeleteMany failed deleting bidders", zap.Error(err))

		return fmt.Errorf("BidderRepository.DeleteMany %w", err)
	}

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
