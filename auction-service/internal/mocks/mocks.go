package mocks

import (
	"context"

	"github.com/ireuven89/auctions/auction-service/domain"
	"github.com/stretchr/testify/mock"

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

// MockAuctionRepository mocks the db.Repository interface
type MockAuctionRepository struct {
	mock.Mock
}

func (m *MockAuctionRepository) Find(ctx context.Context, id string) (auction.Auction, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(auction.Auction), args.Error(1)
}

func (m *MockAuctionRepository) FindAll(ctx context.Context, request auction.AuctionRequest) ([]auction.Auction, error) {
	args := m.Called(ctx, request)
	return args.Get(0).([]auction.Auction), args.Error(1)
}

func (m *MockAuctionRepository) Update(ctx context.Context, req auction.AuctionRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func (m *MockAuctionRepository) Create(ctx context.Context, req auction.AuctionRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func (m *MockAuctionRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockAuctionRepository) DeleteMany(ctx context.Context, ids []interface{}) error {
	args := m.Called(ctx, ids)
	return args.Error(0)
}

type MockItemRepository struct {
	GetItemFunc                       func(ctx context.Context, id string) (domain.Item, error)
	GetItemsBuAuctionFunc             func(ctx context.Context, auctionId string) ([]domain.Item, error)
	GeItemsWByAuctionWithPicturesFunc func(ctx context.Context, auctionId string) ([]domain.Item, error)
}

func (m *MockItemRepository) GetItem(ctx context.Context, id string) (domain.Item, error) {

	return m.GetItemFunc(ctx, id)
}
func (m *MockItemRepository) GetItemsBuAuction(ctx context.Context, id string) ([]domain.Item, error) {

	return m.GeItemsWByAuctionWithPicturesFunc(ctx, id)
}
func (m *MockItemRepository) GeItemsWByAuctionWithPictures(ctx context.Context, id string) ([]domain.Item, error) {

	return m.GeItemsWByAuctionWithPicturesFunc(ctx, id)
}
