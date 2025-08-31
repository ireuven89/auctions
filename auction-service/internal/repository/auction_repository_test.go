package repository

import (
	"context"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/ireuven89/auctions/auction-service/domain"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

type TestUpdate struct {
	Name          string
	Request       domain.AuctionRequest
	ExpectedQuery string
	ExpectedArgs  []interface{}
	ExpectedErr   bool
}

func TestUpdateQuery(t *testing.T) {
	tests := []TestUpdate{
		{
			Name:          "single query",
			Request:       domain.AuctionRequest{ID: "testdata-id", Description: "name"},
			ExpectedQuery: "UPDATE auctions SET description = ? WHERE id = ?",
			ExpectedArgs:  []interface{}{"name", "testdata-id"},
			ExpectedErr:   false,
		},
		{
			Name:          "testdata empty",
			Request:       domain.AuctionRequest{ID: "testdata-id"},
			ExpectedQuery: "",
			ExpectedArgs:  nil,
			ExpectedErr:   true,
		},
		{
			Name:          "another query",
			Request:       domain.AuctionRequest{ID: "testdata-id", Description: "name"},
			ExpectedQuery: "UPDATE auctions SET description = ? WHERE id = ?",
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
	req := domain.AuctionRequest{Description: "car"}

	// The query should match the generated WHERE clause
	expectedQuery := regexp.QuoteMeta("SELECT id, description, category, regions, status, initial_offer, created_at, updated_at from auctions where description LIKE '%car%'")

	rows := sqlmock.NewRows([]string{"id", "description", "category", "regions", "status", "initial_offer", "created_at", "updated_at"}).
		AddRow("a1", "car auction", "Vintage", []byte{}, domain.Active.String(), 2.0, time.Now(), time.Now()).
		AddRow("a2", "sports car auction", "Clothing", []byte{}, domain.Active.String(), 2.0, time.Now(), time.Now())

	mock.ExpectQuery(expectedQuery).WillReturnRows(rows)

	auctions, err := r.FindAll(context.Background(), req)
	require.NoError(t, err)
	require.Len(t, auctions, 2)

	require.Equal(t, "a1", auctions[0].ID)
	require.Equal(t, "car auction", auctions[0].Description)
}
