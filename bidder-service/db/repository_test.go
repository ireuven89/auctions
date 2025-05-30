package db

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ireuven89/auctions/bidder-service/bidder"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type testUpdateQuery struct {
	Name          string
	request       bidder.Bidder
	ExpectedQuery string
	ExpectedArgs  []interface{}
	WantErr       bool
}

func TestPrepareUpdateQuery(t *testing.T) {
	tests := []testUpdateQuery{
		{
			Name:          "failed query",
			request:       bidder.Bidder{ID: "testdata-id", Name: "name", Item: "item"},
			ExpectedQuery: "UPDATE bidders SET name = ?, item = ? WHERE id = ?",
			ExpectedArgs:  []interface{}{"name", "item", "testdata-id"},
			WantErr:       false,
		},
		{
			Name:          "failed query",
			request:       bidder.Bidder{ID: "testdata-id"},
			ExpectedQuery: "",
			ExpectedArgs:  nil,
			WantErr:       true,
		},
	}

	for _, test := range tests {
		q, args, err := prepareUpdateQuery(test.request)
		assert.Equal(t, err != nil, test.WantErr)
		assert.Equal(t, test.ExpectedQuery, q)
		assert.Equal(t, test.ExpectedArgs, args)
	}
}

type TestDeleteBidder struct {
	id        string
	wantErr   bool
	sqlResult sql.Result
}

func TestDelete(t *testing.T) {
	db, mock, err := sqlmock.New()
	logger := zap.NewNop()
	repo := NewRepository(db, logger)

	assert.NoError(t, err)

	tests := []TestDeleteBidder{
		{
			id:        "id",
			wantErr:   false,
			sqlResult: sqlmock.NewResult(0, 1),
		},
		{
			id:        "",
			wantErr:   false,
			sqlResult: sqlmock.NewResult(0, 0),
		},
	}

	ctx := context.Background()

	for _, test := range tests {
		mock.ExpectExec("delete from bidders where id = ?").WithArgs(test.id).
			WillReturnResult(test.sqlResult)
		err = repo.Delete(ctx, test.id)
		assert.Equal(t, test.wantErr, err != nil, fmt.Sprintf("want err %v but got %v error is %v", test.wantErr, err != nil, err))
		assert.NoError(t, mock.ExpectationsWereMet(), "expectations where not met")

	}
}
