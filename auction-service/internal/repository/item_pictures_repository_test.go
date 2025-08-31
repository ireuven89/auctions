package repository

import (
	"context"
	"database/sql/driver"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
	//"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ireuven89/auctions/auction-service/domain"
	"github.com/stretchr/testify/require"
)

type TestCreateItemPicture struct {
	Name          string
	input         domain.ItemPicture
	expectedQuery string
	result        driver.Result
	wantErr       bool
}

func TestItemPictureRepository_CreateItemPicture(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	repo := &ItemPictureRepository{
		db:     db,
		logger: zaptest.NewLogger(t),
	}

	ctx := context.Background()

	tests := []TestCreateItemPicture{
		{
			Name: "success",
			input: domain.ItemPicture{
				ItemID:       "mock_item-id",
				ID:           "mock_id",
				DownloadLink: "mockUrl",
			},
			expectedQuery: "(?i)INSERT INTO item_pictures \\(id, item_id, download_link\\) values \\(\\?, \\?, \\?\\)",
			result:        sqlmock.NewResult(1, 1),
			wantErr:       false,
		},
	}

	for _, test := range tests {
		mock.ExpectExec(test.expectedQuery).
			WithArgs(test.input.ID, test.input.ItemID, test.input.DownloadLink).
			WillReturnResult(test.result)
		err = repo.CreateItemPicture(ctx, test.input)
		assert.Equal(t, err != nil, test.wantErr, test.Name)
		assert.NoError(t, mock.ExpectationsWereMet())
	}
}

func TestCreateItemPicture1(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &ItemPictureRepository{db: db, logger: zaptest.NewLogger(t)}
	ctx := context.Background()

	picture := domain.ItemPicture{
		ID:           "pic123",
		ItemID:       "item456",
		DownloadLink: "https://example.com/image.jpg",
	}

	mock.ExpectExec("INSERT INTO item_pictures").
		WithArgs(picture.ID, picture.ItemID, picture.DownloadLink).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.CreateItemPicture(ctx, picture)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}
