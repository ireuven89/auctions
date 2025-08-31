package internal

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/ireuven89/auctions/auction-service/domain"
	"github.com/ireuven89/auctions/auction-service/internal/mocks"
	"github.com/ireuven89/auctions/auction-service/internal/service"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

// MockRepository mocks the db.Repository interface
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Find(ctx context.Context, id string) (domain.Auction, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(domain.Auction), args.Error(1)
}

func (m *MockRepository) FindAll(ctx context.Context, request domain.AuctionRequest) ([]domain.Auction, error) {
	args := m.Called(ctx, request)
	return args.Get(0).([]domain.Auction), args.Error(1)
}

func (m *MockRepository) Update(ctx context.Context, req domain.AuctionRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func (m *MockRepository) Create(ctx context.Context, req domain.AuctionRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func (m *MockRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRepository) DeleteMany(ctx context.Context, ids []interface{}) error {
	args := m.Called(ctx, ids)
	return args.Error(0)
}

func TestCreateAuction(t *testing.T) {
	mockRepo := new(MockRepository)
	itemMockRepo := new(mocks.ItemRepositoryMock)
	logger := zap.NewNop()
	svc := service.NewService(mockRepo, itemMockRepo, logger)

	req := domain.AuctionRequest{Description: "Test Auction", MinIncrement: 1.0, InitialOffer: 1.0}

	// Alternative: More flexible context matching
	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("domain.AuctionRequest")).Return(nil)
	itemMockRepo.On("CreateBulk", mock.Anything, mock.AnythingOfType("[]domain.ItemRequest")).Return(nil)

	id, err := svc.Create(context.Background(), req)

	assert.NoError(t, err)
	assert.NotEmpty(t, id)
	mockRepo.AssertExpectations(t)
}

func TestFetchAuction_NotFound(t *testing.T) {
	mockRepo := new(MockRepository)
	itemMockRepo := new(mocks.ItemRepositoryMock)

	logger := zap.NewNop()
	svc := service.NewService(mockRepo, itemMockRepo, logger)

	mockRepo.On("Find", mock.Anything, "not_found").Return(domain.Auction{}, sql.ErrNoRows)

	res, err := svc.Fetch(context.Background(), "not_found")

	assert.Nil(t, res)
	assert.Equal(t, &domain.AppError{Kind: "not_found", Message: "sql: no rows in result set not found"}, err)
	mockRepo.AssertExpectations(t)
}

func TestUpdateAuction(t *testing.T) {
	mockRepo := new(MockRepository)
	itemMockRepo := new(mocks.ItemRepositoryMock)

	logger := zap.NewNop()
	svc := service.NewService(mockRepo, itemMockRepo, logger)

	req := domain.AuctionRequest{ID: uuid.New().String(), Description: "Updated Auction", CreatedAt: time.Time{}, UpdatedAt: time.Time{}}
	mockRepo.On("Update", context.Background(), mock.AnythingOfType("domain.AuctionRequest")).Return(nil)

	err := svc.Update(context.Background(), req)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDeleteAuction(t *testing.T) {
	mockRepo := new(MockRepository)
	itemMockRepo := new(mocks.ItemRepositoryMock)

	logger := zap.NewNop()
	svc := service.NewService(mockRepo, itemMockRepo, logger)

	id := uuid.New().String()
	mockRepo.On("Delete", mock.Anything, id).Return(nil)

	err := svc.Delete(context.Background(), id)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
