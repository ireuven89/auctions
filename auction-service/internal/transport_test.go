package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ireuven89/auctions/auction-service/auction"
	"github.com/ireuven89/auctions/auction-service/internal/mocks"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
)

func TestGetAuction(t *testing.T) {
	s := &mocks.MockAuctionService{
		FetchFunc: func(ctx context.Context, id string) (*auction.Auction, error) {
			return &auction.Auction{ID: id, Name: "Test Auction"}, nil
		},
	}
	r := httprouter.New()
	NewTransport(s, r)

	req := httptest.NewRequest(http.MethodGet, "/auctions/123", nil)
	req = req.WithContext(context.WithValue(req.Context(), httprouter.ParamsKey, httprouter.Params{{Key: "id", Value: "123"}}))
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestCreateAuctionTransport(t *testing.T) {
	s := &mocks.MockAuctionService{
		CreateFunc: func(ctx context.Context, a auction.AuctionRequest) (string, error) {
			return "created-id", nil
		},
	}
	r := httprouter.New()
	NewTransport(s, r)

	body := map[string]interface{}{"id": "created-id", "name": "New Auction"}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/auctions", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)

	var result map[string]string
	json.NewDecoder(resp.Body).Decode(&result)
	assert.Equal(t, "created-id", result["id"])
}

func TestUpdateAuctionTransport(t *testing.T) {
	s := &mocks.MockAuctionService{
		UpdateFunc: func(ctx context.Context, a auction.AuctionRequest) error {
			return nil
		},
	}
	r := httprouter.New()
	NewTransport(s, r)

	body := map[string]interface{}{"name": "Updated Auction"}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPut, "/auctions/456", bytes.NewBuffer(b))
	req = req.WithContext(context.WithValue(req.Context(), httprouter.ParamsKey, httprouter.Params{{Key: "id", Value: "456"}}))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestDeleteAuctionTransport(t *testing.T) {
	s := &mocks.MockAuctionService{
		DeleteFunc: func(ctx context.Context, id string) error {
			return nil
		},
	}
	r := httprouter.New()
	NewTransport(s, r)

	req := httptest.NewRequest(http.MethodDelete, "/auctions/789", bytes.NewBuffer([]byte(`{"id":"789"}`)))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestSearchAuctionTransport(t *testing.T) {
	s := &mocks.MockAuctionService{
		SearchFunc: func(ctx context.Context, req auction.AuctionRequest) ([]auction.Auction, error) {
			return nil, nil
		},
	}
	r := httprouter.New()
	NewTransport(s, r)

	req := httptest.NewRequest(http.MethodGet, "/auctions?name=test", bytes.NewBuffer([]byte("name: test")))
	req = req.WithContext(context.WithValue(req.Context(), httprouter.ParamsKey, httprouter.Params{{Key: "id", Value: "456"}}))
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)
}
