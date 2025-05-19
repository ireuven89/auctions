package mocks

import (
	"context"

	"github.com/ireuven89/auctions/auction-service/auction"
)

type MockAuctionService struct {
	FetchFunc      func(ctx context.Context, id string) (*auction.Auction, error)
	CreateFunc     func(ctx context.Context, auction auction.AuctionRequest) (string, error)
	UpdateFunc     func(ctx context.Context, auction2 auction.AuctionRequest) error
	DeleteFunc     func(ctx context.Context, id string) error
	DeleteManyFunc func(ctx context.Context, ids []string) error
	SearchFunc     func(ctx context.Context, request auction.AuctionRequest) ([]auction.Auction, error)
}

func (m *MockAuctionService) Fetch(ctx context.Context, id string) (*auction.Auction, error) {
	return m.FetchFunc(ctx, id)
}

func (m *MockAuctionService) Create(ctx context.Context, auction auction.AuctionRequest) (string, error) {
	return m.CreateFunc(ctx, auction)
}

func (m *MockAuctionService) Update(ctx context.Context, request auction.AuctionRequest) error {
	return m.UpdateFunc(ctx, request)
}

func (m *MockAuctionService) Delete(ctx context.Context, id string) error {
	return m.DeleteFunc(ctx, id)
}

func (m *MockAuctionService) DeleteMany(ctx context.Context, ids []string) error {
	return m.DeleteManyFunc(ctx, ids)
}

func (m *MockAuctionService) Search(ctx context.Context, request auction.AuctionRequest) ([]auction.Auction, error) {
	return m.SearchFunc(ctx, request)
}
