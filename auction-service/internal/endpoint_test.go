package internal

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/ireuven89/auctions/auction-service/auction"
	"github.com/ireuven89/auctions/auction-service/internal/mocks"
)

func TestMakeEndpointGetAuction(t *testing.T) {
	mockService := &mocks.MockAuctionService{
		FetchFunc: func(ctx context.Context, id string) (*auction.Auction, error) {
			if id == "123" {
				return &auction.Auction{ID: "123", Name: "Test Auction"}, nil
			}
			return &auction.Auction{}, fmt.Errorf("not found")
		},
	}

	endpoint := MakeEndpointGetAuction(mockService)

	// Case: success
	req := GetAuctionRequestModel{id: "123"}
	resp, err := endpoint(context.Background(), req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	responseModel, ok := resp.(GetAuctionResponseModel)
	if !ok {
		t.Fatalf("response is not of type GetAuctionResponseModel")
	}

	if responseModel.auction.ID != "123" {
		t.Errorf("expected auction ID '123', got '%s'", responseModel.auction.ID)
	}

	// Case: error
	req = GetAuctionRequestModel{id: "not-exist"}
	_, err = endpoint(context.Background(), req)
	if err == nil {
		t.Fatalf("expected error for missing auction, got nil")
	}
}

func TestMakeEndpointUpdateAuction(t *testing.T) {
	mockService := &mocks.MockAuctionService{
		UpdateFunc: func(ctx context.Context, a auction.AuctionRequest) error {
			if a.ID == "" {
				return errors.New("missing ID")
			}
			return nil
		},
	}

	endpoint := MakeEndpointUpdateAuction(mockService)

	// Test success case
	req := UpdateAuctionRequestModel{
		AuctionRequest: auction.AuctionRequest{ID: "123", Name: "Updated Auction"},
	}
	resp, err := endpoint(context.Background(), req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp != nil {
		t.Errorf("expected nil response, got %v", resp)
	}

	// Test failure case
	req = UpdateAuctionRequestModel{
		AuctionRequest: auction.AuctionRequest{ID: "", Name: "Invalid"},
	}
	_, err = endpoint(context.Background(), req)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestMakeEndpointCreateAuction(t *testing.T) {
	mockService := &mocks.MockAuctionService{
		CreateFunc: func(ctx context.Context, a auction.AuctionRequest) (string, error) {
			if a.ID == "" {
				return "", errors.New("missing ID")
			}
			return a.ID, nil
		},
	}

	endpoint := MakeEndpointCreateAuction(mockService)

	// Test success case
	req := CreateAuctionRequestModel{
		AuctionRequest: auction.AuctionRequest{ID: "123", Name: "Updated Auction"},
	}
	resp, err := endpoint(context.Background(), req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	expected := CreateAuctionResponseModel{id: "123"}

	if resp != expected {
		t.Errorf("expected response %v response, got %v", expected, resp)
	}

	// Test failure case
	req = CreateAuctionRequestModel{
		AuctionRequest: auction.AuctionRequest{ID: "", Name: "Invalid"},
	}
	_, err = endpoint(context.Background(), req)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
