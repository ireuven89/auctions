package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/ireuven89/auctions/auction-service/domain"
	"go.uber.org/zap"
)

type ItemDB struct {
	ID          string `db:"id"`
	AuctionID   string `db:"auction_id"`
	Description string `db:"description"`
}

func ToItem(item ItemDB) domain.Item {
	return domain.Item{
		ID:          item.ID,
		Description: item.Description,
		AuctionID:   item.AuctionID,
	}
}

type ItemPictureDB struct {
	ID           string `db:"id"`
	DownloadLink string `db:"download_link"`
	ItemId       string `db:"item_id"`
	AuctionID    string `db:"auction_id"`
}

func toItemResponse(item ItemPictureDB) domain.ItemPictureResponse {

	return domain.ItemPictureResponse{
		ID:           item.ID,
		AuctionID:    item.AuctionID,
		DownloadLink: item.ItemId,
	}
}

type ItemRepository struct {
	logger *zap.Logger
	db     *sql.DB
}

func NewItemRepo(db *sql.DB, logger *zap.Logger) *ItemRepository {

	return &ItemRepository{
		logger: logger,
		db:     db,
	}
}

func (r *ItemRepository) Find(ctx context.Context, id string) (domain.Item, error) {
	var itemDB ItemDB

	q := `select id, description, auction_id from items where id = ?`

	row := r.db.QueryRowContext(ctx, q, id)

	if row.Err() != nil {
		return domain.Item{}, fmt.Errorf("ItemRepository.Find %w", row.Err())
	}

	if err := row.Scan(&itemDB.ID, &itemDB.Description, &itemDB.AuctionID); err != nil {
		return domain.Item{}, fmt.Errorf("ItemRepository.Find %w", err)
	}

	return ToItem(itemDB), nil
}
func (r *ItemRepository) FindWithPictures(ctx context.Context, auctionId string) ([]domain.ItemPictureResponse, error) {

	q := `select it.id, it.description, it.auction_id, itp.download_link 
		  from items it
		  join items_pictures itp  on itp.item_id = it.id 
          where it.auction_id = ?`

	row, err := r.db.QueryContext(ctx, q, auctionId)

	if err != nil {
		return nil, fmt.Errorf("ItemRepository.FindWithPictures %w", row.Err())
	}

	var itemPictureDB ItemPictureDB
	var response []domain.ItemPictureResponse

	for row.Next() {
		if err = row.Scan(&itemPictureDB.ID, *&itemPictureDB, &itemPictureDB.AuctionID, &itemPictureDB.DownloadLink); err != nil {
			return nil, fmt.Errorf("ItemRepository.FindWithPictures %w", err)
		}
		response = append(response, toItemResponse(itemPictureDB))
	}

	return response, nil
}
func (r *ItemRepository) Update(ctx context.Context, request domain.ItemRequest) error {
	tx, err := r.db.Begin()

	if err != nil {
		return fmt.Errorf("ItemRepository.Update %w", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	if request.Description != "" {
		tx.QueryContext(ctx, `update items set name = $1 where id = $2`, request.Description, request.ID)
	}

	if len(request.Pictures) > 0 {
		pictureStatement, err := tx.Prepare(`update items set download_link = $1 where id = $2`)
		if err != nil {
			return fmt.Errorf("ItemRepository.FindWithPictures %w", err)
		}
		for _, picture := range request.Pictures {
			pictureStatement.ExecContext(ctx, picture.DownloadLink, picture.ID)
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("ItemRepository.FindWithPictures %w", err)
	}

	return nil
}

func (r *ItemRepository) Create(ctx context.Context, item domain.ItemRequest) error {
	tx, err := r.db.Begin()

	defer tx.Rollback()

	if err != nil {
		return fmt.Errorf("ItemRepository.Create %w", err)
	}

	//items insert
	q := `insert into items (id, name, auction_id) values(?, ?, ?)`
	if _, err := tx.ExecContext(ctx, q, item.ID, item.Description, item.AuctionID); err != nil {
		return fmt.Errorf("ItemRepository.Create %w", err)
	}
	//items pictures insert
	pictuersQ := prepareInsertItemsPictures(item.Pictures)
	if _, err = tx.ExecContext(ctx, pictuersQ, item.ID, item.Description, item.AuctionID); err != nil {
		return fmt.Errorf("ItemRepository.Create %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("ItemRepository.Create %w", err)
	}

	return nil
}

func prepareInsertItemsPictures(itemPictures []domain.ItemPicture) string {
	var values []interface{}
	var placeHolders []string
	//q := `insert into items_pictures (id, item_id, downloaLink) values `

	for i, item := range itemPictures {
		placeholder := fmt.Sprintf("($%d, $%d, $%d)", i*3+1, i*3+2, i*3+3)
		placeHolders = append(placeHolders, placeholder)
		values = append(values, item.ID, item.ItemID, item.DownloadLink)
	}

	query := fmt.Sprintf("INSERT INTO items_picture (id, item_id, download_link) VALUES %s",
		strings.Join(placeHolders, ", "))

	return query

}
func (r *ItemRepository) Delete(ctx context.Context, id string) error {
	q := `delete from items where id = ?`
	if _, err := r.db.ExecContext(ctx, q, id); err != nil {
		return fmt.Errorf("ItemRepository.Delete %w", err)
	}

	return nil
}

func (r *ItemRepository) FindByAuctionId(ctx context.Context, auctionId string) ([]domain.Item, error) {
	q := `select id, descrtption, auction_id, opening_price,  from items where id = ?`

	rows, err := r.db.QueryContext(ctx, q, auctionId)

	if err != nil {
		return nil, err
	}

	var result []domain.Item
	var itemDBs []ItemDB

	for rows.Next() {
		var curr ItemDB
		if err = rows.Scan(&curr.ID, &curr.Description, &curr.AuctionID); err != nil {
			return nil, err
		}
		itemDBs = append(itemDBs, curr)
	}

	for _, item := range itemDBs {
		result = append(result, ToItem(item))
	}

	return result, nil
}

func (r *ItemRepository) CreateBulk(ctx context.Context, request []domain.ItemRequest) error {
	tx, err := r.db.Begin()

	if err != nil {
		fmt.Errorf("ItemRepository.CreateBulk %w", err)
	}

	defer func() {
		tx.Rollback()
	}()

	//create items
	qPrefix := `INSERT INTO items (id, description, auction_id) VALUES %s`
	var itemValues []interface{}
	var placeHolders []string

	for _, item := range request {
		placeHolders = append(placeHolders, "(?, ?, ?)")
		itemValues = append(itemValues, item.ID, item.Description, item.AuctionID)
	}

	q := fmt.Sprintf(qPrefix, strings.Join(placeHolders, ","))
	_, err = tx.ExecContext(ctx, q, itemValues...)
	if err != nil {
		return fmt.Errorf("ItemRepository.CreateBulk failed to insert items: %w", err)
	}

	//item picture insert
	var itemPictureValues []interface{}
	placeHolders = []string{}

	for _, item := range request {
		for _, picture := range item.Pictures {
			placeHolders = append(placeHolders, "(?, ?, ?)")
			itemPictureValues = append(itemPictureValues, picture.ID, picture.ItemID, picture.DownloadLink)
		}
	}

	query := fmt.Sprintf("INSERT INTO item_pictures (id, item_id, download_link) VALUES %s",
		strings.Join(placeHolders, ","))

	tx.ExecContext(ctx, query, itemPictureValues)

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("ItemRepository.CreateBulk %w", err)
	}

	return nil
}
