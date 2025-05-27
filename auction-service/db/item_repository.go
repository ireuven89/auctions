package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ireuven89/auctions/auction-service/domain"
	"github.com/ireuven89/auctions/auction-service/internal"
	"go.uber.org/zap"
)

type ItemDB struct {
	ID        string `db:"id"`
	Name      string `db:"name"`
	AuctionID string `db:"auction_id"`
}

type ItemWithPictureDB struct {
	ItemID       string `db:"item_id"`
	AuctionID    string `db:"auction_id"`
	Name         string `db:"name"`
	PictureID    string `db:"picture_id"`
	DownloadLink string `db:"download_link"`
}

func toItem(db ItemDB) domain.Item {

	return domain.Item{
		ID:          db.ID,
		Description: db.Name,
		AuctionID:   db.AuctionID,
	}
}

func toItemPicture(db ItemWithPictureDB) domain.ItemPicture {

	return domain.ItemPicture{
		ItemID:      db.ItemID,
		DownloadUrl: db.DownloadLink,
	}
}

type ItemRepo struct {
	db     *sql.DB
	logger *zap.Logger
}

func NewItemRepo(logger *zap.Logger, db *sql.DB) internal.ItemRepository {

	return &ItemRepo{
		logger: logger,
		db:     db,
	}
}

func (r *ItemRepo) GetItem(ctx context.Context, id string) (domain.Item, error) {
	var itemDB ItemDB
	row := r.db.QueryRowContext(ctx, "select id, auction_id, name from items where id = ?", id)

	if row.Err() != nil {
		return domain.Item{}, fmt.Errorf("ItemRepo.GetItem %w", row.Err())
	}

	if err := row.Scan(&itemDB.ID, &itemDB.AuctionID, &itemDB.Name); err != nil {
		return domain.Item{}, fmt.Errorf("ItemRepo.GetItem %w", err)
	}

	item := toItem(itemDB)

	return item, nil
}
func (r *ItemRepo) GetItemsBuAuction(ctx context.Context, auctionId string) ([]domain.Item, error) {
	var itemDB ItemDB
	var result []domain.Item
	rows, err := r.db.QueryContext(ctx, "select id, auction_id, name from items where auction_id = ?", auctionId)

	if err != nil {
		return nil, fmt.Errorf("ItemRepo.GetItem %w", err)
	}

	for rows.Next() {
		if err = rows.Scan(&itemDB.ID, &itemDB.AuctionID, &itemDB.Name); err != nil {
			return nil, fmt.Errorf("ItemRepo.GetItem %w", err)
		}
		result = append(result, toItem(itemDB))
	}

	return result, nil
}

func (r *ItemRepo) GeItemsWByAuctionWithPictures(ctx context.Context, auctionId string) ([]domain.Item, error) {
	var itemDB ItemWithPictureDB
	var result []domain.ItemPicture

	q := `select it.id as item_id , it.auction_id, it.name, itp.download_link 
		  from items it
		  join items_pictures itp on it.id = itp.item_id 
		  where it.auction_id = ?`

	rows, err := r.db.QueryContext(ctx, q, auctionId)

	if err != nil {
		return nil, fmt.Errorf("ItemRepo.GeItemsWByAuctionWithPictures %w", err)
	}

	for rows.Next() {
		if err = rows.Scan(&itemDB.ItemID, &itemDB.AuctionID, &itemDB.Name, &itemDB.DownloadLink); err != nil {
			return nil, fmt.Errorf("ItemRepo.GetItem %w", err)
		}
		result = append(result, toItemPicture(itemDB))
	}

	return nil, err
}
