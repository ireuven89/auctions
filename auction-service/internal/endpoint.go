package internal

import (
	"context"
	"fmt"
	"mime/multipart"

	"github.com/ireuven89/auctions/auction-service/domain"
	"github.com/ireuven89/auctions/auction-service/internal/service"

	"github.com/go-kit/kit/endpoint"
)

type GetAuctionRequestModel struct {
	id string
}

type GetAuctionResponseModel struct {
	auction *domain.Auction
}

func MakeEndpointGetAuction(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(GetAuctionRequestModel)
		if !ok {
			return nil, fmt.Errorf("MakeEndpointGetAuction.failed parsing request")
		}

		res, err := s.Fetch(ctx, req.id)

		if err != nil {

			return nil, fmt.Errorf("MakeEndpointGetAuction %w", err)
		}

		return GetAuctionResponseModel{
			auction: res,
		}, nil
	}
}

type GetAuctionsRequestModel struct {
	domain.AuctionRequest
}

type GetAuctionsResponseModel struct {
	auctions []domain.Auction
}

func MakeEndpointGetAuctions(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(GetAuctionsRequestModel)
		if !ok {
			return nil, fmt.Errorf("MakeEndpointGetAuctions.failed parsing request")
		}

		res, err := s.Search(ctx, req.AuctionRequest)

		if err != nil {

			return nil, fmt.Errorf("MakeEndpointGetAuctions %w", err)
		}

		return GetAuctionsResponseModel{
			auctions: res,
		}, nil
	}
}

type CreateAuctionRequestModel struct {
	domain.AuctionRequest
}

type CreateAuctionResponseModel struct {
	id string
}

func MakeEndpointCreateAuction(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(CreateAuctionRequestModel)
		if !ok {
			return nil, fmt.Errorf("MakeEndpointGetAuction.failed parsing request")
		}

		res, err := s.Create(ctx, req.AuctionRequest)

		if err != nil {
			return nil, fmt.Errorf("MakeEndpointCreateAuction %w", err)
		}

		return CreateAuctionResponseModel{id: res}, nil
	}
}

type UpdateAuctionRequestModel struct {
	domain.AuctionRequest
}

func MakeEndpointUpdateAuction(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(UpdateAuctionRequestModel)
		if !ok {
			return nil, fmt.Errorf("MakeEndpointUpdateAuction.failed parsing request")
		}

		err = s.Update(ctx, req.AuctionRequest)

		if err != nil {
			return nil, fmt.Errorf("MakeEndpointUpdateAuction %w", err)
		}

		return nil, nil
	}
}

type DeleteAuctionRequestModel struct {
	id string
}

func MakeEndpointDeleteAuction(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(DeleteAuctionRequestModel)
		if !ok {
			return nil, fmt.Errorf("MakeEndpointDeleteAuction.failed parsing request")
		}

		err = s.Delete(ctx, req.id)

		if err != nil {
			return nil, fmt.Errorf("MakeEndpointDeleteAuction %w", err)
		}

		return nil, nil
	}
}

type AuctionItemsRequestModel struct {
	AuctionID string
	items     []domain.Item
}

func MakeEndpointCreateAuctionItems(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(AuctionItemsRequestModel)
		if !ok {
			return nil, fmt.Errorf("MakeEndpointCreateAuctionItems.failed parsing request")
		}

		err = s.CreateAuctionItems(ctx, req.AuctionID, req.items)

		if err != nil {
			return nil, fmt.Errorf("MakeEndpointCreateAuctionItems %w", err)
		}

		return nil, nil
	}
}

type AuctionPicturesRequestModel struct {
	ItemID string
	Files  []*multipart.FileHeader
}

func MakeEndpointCreateAuctionItemPictures(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(AuctionPicturesRequestModel)
		if !ok {
			return nil, fmt.Errorf("MakeEndpointDeleteAuction.failed parsing request")
		}

		err = s.CreateAuctionPictures(ctx, req.ItemID, req.Files)

		if err != nil {
			return nil, fmt.Errorf("MakeEndpointCreateAuctionItems %w", err)
		}

		return nil, nil
	}
}

type CreateItemRequestModel struct {
	req domain.ItemRequest
}
