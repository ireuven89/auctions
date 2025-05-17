package internal

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	"github.com/ireuven89/auctions/auction-service/auction"
)

type GetAuctionRequestModel struct {
	id string
}

type GetAuctionResponseModel struct {
	auction *auction.Auction
}

func MakeEndpointGetAuction(s Service) endpoint.Endpoint {
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

type CreateAuctionRequestModel struct {
	auction.AuctionRequest
}

type CreateAuctionResponseModel struct {
	id string
}

func MakeEndpointCreateAuction(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(CreateAuctionRequestModel)
		if !ok {
			return nil, fmt.Errorf("MakeEndpointGetAuction.failed parsing request")
		}

		res, err := s.Create(ctx, req.AuctionRequest)

		if err != nil {
			return nil, err
		}

		return CreateAuctionResponseModel{id: res}, nil
	}
}

type UpdateAuctionRequestModel struct {
	auction.AuctionRequest
}

func MakeEndpointUpdateAuction(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(UpdateAuctionRequestModel)
		if !ok {
			return nil, fmt.Errorf("MakeEndpointUpdateAuction.failed parsing request")
		}

		err = s.Update(ctx, req.AuctionRequest)

		if err != nil {
			return nil, err
		}

		return nil, nil
	}
}

type DeleteAuctionRequestModel struct {
	id string
}

func MakeEndpointDeleteAuction(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(DeleteAuctionRequestModel)
		if !ok {
			return nil, fmt.Errorf("MakeEndpointDeleteAuction.failed parsing request")
		}

		err = s.Delete(ctx, req.id)

		if err != nil {
			return nil, err
		}

		return nil, nil
	}
}
