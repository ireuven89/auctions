package internal

import (
	"context"
	"github.com/ireuven89/auctions/bidder-service/bidder"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type mockService struct {
	mock.Mock
}

func (m *mockService) GetBidder(ctx context.Context, id string) (bidder.Bidder, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(bidder.Bidder), args.Error(1)
}

func (m *mockService) CreateBidder(ctx context.Context, b bidder.Bidder) (string, error) {
	args := m.Called(ctx, b)
	return args.String(0), args.Error(1)
}

func (m *mockService) UpdateBidder(ctx context.Context, b bidder.Bidder) error {
	return m.Called(ctx, b).Error(0)
}

func (m *mockService) DeleteBidder(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}

func (m *mockService) SearchBidders(ctx context.Context, request bidder.BiddersRequest) ([]bidder.Bidder, error) {
	args := m.Called(ctx, request)

	return args.Get(0).([]bidder.Bidder), args.Error(1)
}

func TestMakeEndpointGetBidder(t *testing.T) {
	mockSvc := new(mockService)
	ctx := context.Background()

	expectedBidder := bidder.Bidder{ID: "1", Name: "John"}
	mockSvc.On("GetBidder", ctx, "1").Return(expectedBidder, nil)

	endpoint := MakeEndpointGetBidder(mockSvc)

	resp, err := endpoint(ctx, GetBidderRequestModel{id: "1"})
	assert.NoError(t, err)

	responseModel, ok := resp.(GetBidderResponseModel)
	assert.True(t, ok)
	assert.Equal(t, expectedBidder, responseModel.bidder)
}

func TestMakeEndpointCreateBidder(t *testing.T) {
	mockSvc := new(mockService)
	ctx := context.Background()

	newBidder := bidder.Bidder{Name: "Alice"}
	mockSvc.On("CreateBidder", ctx, newBidder).Return("123", nil)

	endpoint := MakeEndpointCreateBidder(mockSvc)

	resp, err := endpoint(ctx, CreateBidderRequestModel{bidder: newBidder})
	assert.NoError(t, err)

	responseModel, ok := resp.(CreateBidderResponseModel)
	assert.True(t, ok)
	assert.Equal(t, "123", responseModel.id)
}

func TestMakeEndpointUpdateBidder(t *testing.T) {
	mockSvc := new(mockService)
	ctx := context.Background()

	updated := bidder.Bidder{ID: "1", Name: "Updated"}
	mockSvc.On("UpdateBidder", ctx, updated).Return(nil)

	endpoint := MakeEndpointUpdateBidder(mockSvc)

	resp, err := endpoint(ctx, UpdateBidderRequestModel{bidder: updated})
	assert.NoError(t, err)
	assert.Nil(t, resp)
}

func TestMakeEndpointDeleteBidder(t *testing.T) {
	mockSvc := new(mockService)
	ctx := context.Background()

	mockSvc.On("DeleteBidder", ctx, "1").Return(nil)

	endpoint := MakeEndpointDeleteBidder(mockSvc)

	resp, err := endpoint(ctx, DeleteBidderRequestModel{id: "1"})
	assert.NoError(t, err)
	assert.Nil(t, resp)
}
