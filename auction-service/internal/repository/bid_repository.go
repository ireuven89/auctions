package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ireuven89/auctions/auction-service/domain"
	"go.uber.org/zap"
)

type BidRepository struct {
	logger *zap.Logger
	db     *sql.DB
}

func (r *BidRepository) Find(ctx context.Context, id string) (domain.Bid, error) {
	_, err := r.db.ExecContext(ctx, "select id, auction_id, bidder_id, price from bids where id = ?", id)

	if err != nil {
		return domain.Bid{}, fmt.Errorf("BidRepository.Find %w", err)
	}

	return domain.Bid{}, nil
}

func (r *BidRepository) Create(ctx context.Context, bid domain.Bid) error {
	_, err := r.db.ExecContext(ctx, "insert into bids (id, auction_id, bidder_id, price) values (?,?,?,?,?)", bid.ID, bid.BidderID, bid.AuctionID, bid.Winner, bid.CreateAt)

	if err != nil {
		return fmt.Errorf("BidRepository.Create %w", err)
	}

	return nil
}

func (r *BidRepository) Update(ctx context.Context, bid domain.Bid) error {
	_, err := r.db.ExecContext(ctx, "update bids set price = ? where id = ?", bid.Price, bid.ID)

	if err != nil {
		return fmt.Errorf("BidRepository.Update %w", err)
	}

	return nil
}

func (r *BidRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, "delete from bids where id = ?", id)

	if err != nil {
		return fmt.Errorf("BidRepository.Delete failed deleting %w", err)
	}

	return nil
}
