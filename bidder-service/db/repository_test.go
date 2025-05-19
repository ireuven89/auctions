package db

import (
	"testing"

	"github.com/ireuven89/auctions/bidder-service/bidder"
	"github.com/stretchr/testify/assert"
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
			request:       bidder.Bidder{ID: "test-id", Name: "name", Item: "item"},
			ExpectedQuery: "UPDATE bidders SET name = ?, item = ? WHERE id = ?",
			ExpectedArgs:  []interface{}{"name", "item", "test-id"},
			WantErr:       false,
		},
		{
			Name:          "failed query",
			request:       bidder.Bidder{ID: "test-id"},
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
