package internal

import (
	"context"
	"database/sql"
	"github.com/go-redis/redismock/v9"
	"github.com/ireuven89/auctions/auction-service/db"
	"testing"

	"github.com/google/uuid"
	"github.com/ireuven89/auctions/auction-service/auction"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

// MockRepository mocks the db.Repository interface
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) WithTransactionContext(ctx context.Context, fn func(txRepo db.Repository) error) error {
	args := m.Called(ctx, fn)

	return args.Error(0)
}

func (m *MockRepository) Find(ctx context.Context, id string) (auction.Auction, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(auction.Auction), args.Error(1)
}

func (m *MockRepository) FindAll(ctx context.Context, request auction.AuctionRequest) ([]auction.Auction, error) {
	args := m.Called(ctx, request)
	return args.Get(0).([]auction.Auction), args.Error(1)
}

func (m *MockRepository) Update(ctx context.Context, req auction.AuctionRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func (m *MockRepository) Create(ctx context.Context, req auction.AuctionRequest) error {
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
	redisClient, _ := redismock.NewClientMock()
	logger := zap.NewNop()
	svc := NewService(mockRepo, redisClient, logger)

	req := auction.AuctionRequest{Name: "Test Auction", Description: "Test Description", UserId: uuid.New().String()}
	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("auction.AuctionRequest")).Return(nil)

	id, err := svc.Create(context.Background(), req)

	assert.NoError(t, err)
	assert.NotEmpty(t, id)
	mockRepo.AssertExpectations(t)
}

func TestFetchAuction_NotFound(t *testing.T) {
	mockRepo := new(MockRepository)
	redisClient, _ := redismock.NewClientMock()
	logger := zap.NewNop()
	svc := NewService(mockRepo, redisClient, logger)

	mockRepo.On("Find", mock.Anything, "not_found").Return(auction.Auction{}, sql.ErrNoRows)

	res, err := svc.Fetch(context.Background(), "not_found")

	assert.Nil(t, res)
	assert.Equal(t, auction.ErrNotFound, err)
	mockRepo.AssertExpectations(t)
}

func TestUpdateAuction(t *testing.T) {
	mockRepo := new(MockRepository)
	redisClient, _ := redismock.NewClientMock()
	logger := zap.NewNop()
	svc := NewService(mockRepo, redisClient, logger)

	req := auction.AuctionRequest{ID: uuid.New().String(), Name: "Updated Auction"}
	mockRepo.On("Update", mock.Anything, req).Return(nil)

	err := svc.Update(context.Background(), req)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDeleteAuction(t *testing.T) {
	mockRepo := new(MockRepository)
	redisClient, _ := redismock.NewClientMock()
	logger := zap.NewNop()
	svc := NewService(mockRepo, redisClient, logger)

	id := uuid.New().String()
	mockRepo.On("Delete", mock.Anything, id).Return(nil)

	err := svc.Delete(context.Background(), id)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
