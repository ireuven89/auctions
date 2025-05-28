package repository

import (
	"context"
	"database/sql"
	"github.com/ireuven89/auctions/auction-service/domain"
	"go.uber.org/zap"
)

type BidRepository struct {
	logger *zap.Logger
	db     *sql.DB
}

func (r *BidRepository) Find(ctx context.Context, id string) (domain.Bid, error) {

	return domain.Bid{}, nil
}

func (r *BidRepository) Create(ctx context.Context, bid domain.Bid) error {
	_, err := r.db.ExecContext(ctx, "insert into bids (id, auction_id, bidder_id, price, winner) values (?,?,?,?,?)", bid.ID, bid.BidderID, bid.AuctionID, bid.Winner, bid.CreateAt)

	if err != nil {
		return err
	}

	return nil
}
