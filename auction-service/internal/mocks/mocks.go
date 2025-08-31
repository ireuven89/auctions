package mocks

import (
	"context"
	"mime/multipart"

	"github.com/ireuven89/auctions/auction-service/domain"
	"github.com/stretchr/testify/mock"
)

type MockAuctionService struct {
	FetchFunc                 func(ctx context.Context, id string) (*domain.Auction, error)
	CreateFunc                func(ctx context.Context, auction domain.AuctionRequest) (string, error)
	UpdateFunc                func(ctx context.Context, auction2 domain.AuctionRequest) error
	DeleteFunc                func(ctx context.Context, id string) error
	DeleteManyFunc            func(ctx context.Context, ids []string) error
	SearchFunc                func(ctx context.Context, request domain.AuctionRequest) ([]domain.Auction, error)
	CreateAuctionItemsFunc    func(ctx context.Context, itemId string, items []domain.Item) error
	CreateAuctionPicturesFunc func(ctx context.Context, id string, request []*multipart.FileHeader) error
	PlaceBidFunc              func(ctx context.Context, bid domain.PlaceBidRequest) error
}

func (m *MockAuctionService) CreateAuctionPictures(ctx context.Context, id string, request []*multipart.FileHeader) error {
	return m.CreateAuctionPicturesFunc(ctx, id, request)
}

func (m *MockAuctionService) Fetch(ctx context.Context, id string) (*domain.Auction, error) {
	return m.FetchFunc(ctx, id)
}

func (m *MockAuctionService) Create(ctx context.Context, auction domain.AuctionRequest) (string, error) {
	return m.CreateFunc(ctx, auction)
}

func (m *MockAuctionService) Update(ctx context.Context, request domain.AuctionRequest) error {
	return m.UpdateFunc(ctx, request)
}

func (m *MockAuctionService) Delete(ctx context.Context, id string) error {
	return m.DeleteFunc(ctx, id)
}

func (m *MockAuctionService) DeleteMany(ctx context.Context, ids []string) error {
	return m.DeleteManyFunc(ctx, ids)
}

func (m *MockAuctionService) Search(ctx context.Context, request domain.AuctionRequest) ([]domain.Auction, error) {
	return m.SearchFunc(ctx, request)
}

func (m *MockAuctionService) PlaceBid(ctx context.Context, bid domain.PlaceBidRequest) error {
	return m.PlaceBidFunc(ctx, bid)
}

func (m *MockAuctionService) CreateAuctionItems(ctx context.Context, itemId string, items []domain.Item) error {
	return m.CreateAuctionItemsFunc(ctx, itemId, items)
}

type ItemRepositoryMock struct {
	mock.Mock
}

func (m *ItemRepositoryMock) Find(ctx context.Context, id string) (domain.Item, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(domain.Item), args.Error(1)
}

func (m *ItemRepositoryMock) FindWithPictures(ctx context.Context, auctionId string) ([]domain.ItemPicture, error) {
	args := m.Called(ctx, auctionId)
	return args.Get(0).([]domain.ItemPicture), args.Error(1)
}

func (m *ItemRepositoryMock) Update(ctx context.Context, item domain.ItemRequest) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

func (m *ItemRepositoryMock) Create(ctx context.Context, item domain.ItemRequest) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

func (m *ItemRepositoryMock) CreateBulk(ctx context.Context, items []domain.Item) error {
	args := m.Called(ctx, items)
	return args.Error(0)
}

func (m *ItemRepositoryMock) Delete(ctx context.Context, id string) error {
	args := m.Called(id)
	return args.Error(0)
}

// MockAuctionRepository mocks the db.Repository interface
type MockAuctionRepository struct {
	mock.Mock
}

func (m *MockAuctionRepository) Find(ctx context.Context, id string) (domain.Auction, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(domain.Auction), args.Error(1)
}

func (m *MockAuctionRepository) FindAll(ctx context.Context, request domain.AuctionRequest) ([]domain.Auction, error) {
	args := m.Called(ctx, request)
	return args.Get(0).([]domain.Auction), args.Error(1)
}

func (m *MockAuctionRepository) Update(ctx context.Context, req domain.AuctionRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func (m *MockAuctionRepository) Create(ctx context.Context, req domain.AuctionRequest) error {
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

func (m *ItemRepositoryMock) FindByAuctionId(ctx context.Context, auctionId string) ([]domain.Item, error) {
	args := m.Called(ctx, auctionId)
	return args.Get(0).([]domain.Item), args.Error(1)
}

func (m *ItemRepositoryMock) CreateItemPicture(ctx context.Context, picture domain.ItemPicture) error {
	args := m.Called(picture)

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

type ItemPicturesRepoMock struct {
	mock.Mock
}

func (m *ItemPicturesRepoMock) CreateItemPicture(ctx context.Context, picture domain.ItemPicture) error {
	args := m.Called(ctx, picture)

	return args.Error(0)
}
func (m *ItemPicturesRepoMock) DeleteItemPicture(ctx context.Context, id string) error {
	args := m.Called(ctx, id)

	return args.Error(0)
}
func (m *ItemPicturesRepoMock) CreateItemPictureBulk(ctx context.Context, pictures []domain.ItemPicture) error {
	args := m.Called(ctx, pictures)

	return args.Error(0)
}
