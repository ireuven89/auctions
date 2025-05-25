package db

import (
	"context"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ireuven89/auctions/auction-service/auction"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

type TestUpdate struct {
	Name          string
	Request       auction.AuctionRequest
	ExpectedQuery string
	ExpectedArgs  []interface{}
	ExpectedErr   bool
}

func TestUpdateQuery(t *testing.T) {
	tests := []TestUpdate{
		{
			Name:          "single query",
			Request:       auction.AuctionRequest{ID: "testdata-id", Name: "name"},
			ExpectedQuery: "UPDATE auctions SET name = ? WHERE id = ?",
			ExpectedArgs:  []interface{}{"name", "testdata-id"},
			ExpectedErr:   false,
		},
		{
			Name:          "testdata empty",
			Request:       auction.AuctionRequest{ID: "testdata-id"},
			ExpectedQuery: "",
			ExpectedArgs:  nil,
			ExpectedErr:   true,
		},
		{
			Name:          "another query",
			Request:       auction.AuctionRequest{ID: "testdata-id", Name: "name"},
			ExpectedQuery: "UPDATE auctions SET name = ? WHERE id = ?",
			ExpectedArgs:  []interface{}{"name", "testdata-id"},
			ExpectedErr:   false,
		},
	}

	for _, test := range tests {
		t.Logf("testing %s", test.Name)
		q, args, err := buildUpdateQuery(test.Request)
		assert.Equal(t, err != nil, test.ExpectedErr, test.Name)
		assert.Equal(t, q, test.ExpectedQuery, test.Name)
		assert.Equal(t, args, test.ExpectedArgs, test.Name)
	}
}

func TestRepo_FindAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	s := strings.TrimSpace("     -042")
	res, err := strconv.ParseInt(s, 10, 64)

	println(res)

	logger := zaptest.NewLogger(t)
	r := &AuctionRepository{
		db:     db,
		logger: logger,
	}

	// Prepare fake request
	req := auction.AuctionRequest{Name: "car"}

	// The query should match the generated WHERE clause
	expectedQuery := regexp.QuoteMeta("SELECT id, name, bidder_id from auctions where name LIKE '%car%'")

	rows := sqlmock.NewRows([]string{"id", "name", "bidder_id"}).
		AddRow("a1", "car auction", "b123").
		AddRow("a2", "sports car auction", "b456")

	mock.ExpectQuery(expectedQuery).WillReturnRows(rows)

	auctions, err := r.FindAll(context.Background(), req)
	require.NoError(t, err)
	require.Len(t, auctions, 2)

	require.Equal(t, "a1", auctions[0].ID)
	require.Equal(t, "car auction", auctions[0].Name)
	require.Equal(t, "b123", auctions[0].BidderId)
}
