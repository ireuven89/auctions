package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/ireuven89/auctions/auction-service/domain"
	"go.uber.org/zap"
)

type ItemPictureDB struct {
	ID           string `db:"id"`
	DownloadLink string `db:"download_link"`
	ItemId       string `db:"item_id"`
	AuctionID    string `db:"auction_id"`
}

func toItemPicture(db ItemPictureDB) domain.ItemPicture {

	return domain.ItemPicture{
		ID:          db.ID,
		DownloadUrl: db.DownloadLink,
		ItemID:      db.ItemId,
	}
}

type ItemPictureRepository struct {
	db     *sql.DB
	logger *zap.Logger
}

func (r *ItemPictureRepository) CreateItemPicture(ctx context.Context, picture domain.ItemPicture) error {
	q := "INSERT INTO items_picture (id, item_id, download_link) values (?, ?, ?)"

	if _, err := r.db.ExecContext(ctx, q, picture.ID, picture.ItemID, picture.DownloadUrl); err != nil {
		return fmt.Errorf("ItemPictureRepository.CreateItemPicture %w", err)
	}
	return nil
}
func (r *ItemPictureRepository) DeleteItemPicture(ctx context.Context, id string) error {
	if _, err := r.db.ExecContext(ctx, "delete from item_pictures where id = ?", id); err != nil {
		return fmt.Errorf("ItemPictureRepository.DeleteItemPicture %w", err)
	}
	return nil
}
func (r *ItemPictureRepository) CreateItemPictureBulk(ctx context.Context, pictures []domain.ItemPicture) error {
	//items pictures insert
	pictuersQ, values := prepareInsertItemsPictures(pictures)
	if _, err := r.db.ExecContext(ctx, pictuersQ, values...); err != nil {
		return fmt.Errorf("ItemRepository.Create %w", err)
	}

	return nil
}

func prepareInsertItemsPictures(itemPictures []domain.ItemPicture) (string, []interface{}) {
	var values []interface{}
	var placeHolders []string

	for i, item := range itemPictures {
		placeholder := fmt.Sprintf("($%d, $%d, $%d)", i*3+1, i*3+2, i*3+3)
		placeHolders = append(placeHolders, placeholder)
		values = append(values, item.ID, item.ItemID, item.DownloadUrl)
	}

	query := fmt.Sprintf("INSERT INTO items_picture (id, item_id, download_link) VALUES %s",
		strings.Join(placeHolders, ", "))

	return query, values

}
