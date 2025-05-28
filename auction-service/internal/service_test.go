package internal

import (
	"context"
	"database/sql"
	"testing"

	"github.com/ireuven89/auctions/auction-service/internal/mocks"

	"github.com/google/uuid"
	"github.com/ireuven89/auctions/auction-service/auction"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func TestCreateAuction(t *testing.T) {
	mockRepo := new(mocks.MockAuctionRepository)
	mockitemRepo := new(mocks.MockItemRepository)
	logger := zap.NewNop()
	svc := NewService(mockRepo, mockitemRepo, logger)

	req := auction.AuctionRequest{Name: "Test Auction"}
	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("auction.AuctionRequest")).Return(nil)

	id, err := svc.Create(context.Background(), req)

	assert.NoError(t, err)
	assert.NotEmpty(t, id)
	mockRepo.AssertExpectations(t)
}

func TestFetchAuction_NotFound(t *testing.T) {
	mockRepo := new(mocks.MockAuctionRepository)
	mockitemRepo := new(mocks.MockItemRepository)
	logger := zap.NewNop()
	svc := NewService(mockRepo, mockitemRepo, logger)

	mockRepo.On("Find", mock.Anything, "not_found").Return(auction.Auction{}, sql.ErrNoRows)

	res, err := svc.Fetch(context.Background(), "not_found")

	assert.Nil(t, res)
	assert.Equal(t, auction.ErrNotFound, err)
	mockRepo.AssertExpectations(t)
}

func TestUpdateAuction(t *testing.T) {
	mockRepo := new(mocks.MockAuctionRepository)
	mockitemRepo := new(mocks.MockItemRepository)
	logger := zap.NewNop()
	svc := NewService(mockRepo, mockitemRepo, logger)

	req := auction.AuctionRequest{ID: uuid.New().String(), Name: "Updated Auction"}
	mockRepo.On("Update", mock.Anything, req).Return(nil)

	err := svc.Update(context.Background(), req)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDeleteAuction(t *testing.T) {
	mockRepo := new(mocks.MockAuctionRepository)
	mockitemRepo := new(mocks.MockItemRepository)
	logger := zap.NewNop()
	svc := NewService(mockRepo, mockitemRepo, logger)

	id := uuid.New().String()
	mockRepo.On("Delete", mock.Anything, id).Return(nil)

	err := svc.Delete(context.Background(), id)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
