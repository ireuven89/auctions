package db

import (
	"github.com/ireuven89/auctions/auction-service/auction"
	"github.com/stretchr/testify/assert"
	"testing"
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
			Request:       auction.AuctionRequest{ID: "test-id", Name: "name"},
			ExpectedQuery: "UPDATE auctions SET name = ? WHERE id = ?",
			ExpectedArgs:  []interface{}{"name", "test-id"},
			ExpectedErr:   false,
		},
		{
			Name:          "test empty",
			Request:       auction.AuctionRequest{ID: "test-id"},
			ExpectedQuery: "",
			ExpectedArgs:  nil,
			ExpectedErr:   true,
		},
		{
			Name:          "another query",
			Request:       auction.AuctionRequest{ID: "test-id", Name: "name"},
			ExpectedQuery: "UPDATE auctions SET name = ? WHERE id = ?",
			ExpectedArgs:  []interface{}{"name", "test-id"},
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
