package internal

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/endpoint"
	"github.com/ireuven89/auctions/bidder-service/bidder"
)

type GetBidderRequestModel struct {
	id string
}

type GetBidderResponseModel struct {
	bidder bidder.Bidder
}

func MakeEndpointGetBidder(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(GetBidderRequestModel)

		result, err := s.GetBidder(ctx, req.id)

		if err != nil {
			return nil, fmt.Errorf("MakeEndpointGetBidder failed getting bidder %w", err)
		}

		return GetBidderResponseModel{
			result,
		}, nil
	}
}

type GetBiddersRequestModel struct {
	bidder.BiddersRequest
}

type GetBiddersResponseModel struct {
	bidders []bidder.Bidder
}

func MakeEndpointGetBidders(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(GetBiddersRequestModel)

		result, err := s.SearchBidders(ctx, req.BiddersRequest)

		if err != nil {
			return nil, fmt.Errorf("MakeEndpointGetBidder failed getting bidder %w", err)
		}

		return GetBiddersResponseModel{
			result,
		}, nil
	}
}

type UpdateBidderRequestModel struct {
	bidder bidder.Bidder
}

func MakeEndpointUpdateBidder(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(UpdateBidderRequestModel)

		if err = s.UpdateBidder(ctx, req.bidder); err != nil {
			return nil, fmt.Errorf("MakeEndpointUpdateBidder %w", err)
		}

		return nil, nil
	}
}

type CreateBidderRequestModel struct {
	bidder bidder.Bidder
}

type CreateBidderResponseModel struct {
	id string
}

func MakeEndpointCreateBidder(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(CreateBidderRequestModel)

		if !ok {
			return nil, fmt.Errorf("MakeEndpointCreateBidder failed casting request ")
		}

		id, err := s.CreateBidder(ctx, req.bidder)

		if err != nil {
			return nil, fmt.Errorf("MakeEndpointCreateBidder %w", err)
		}

		return CreateBidderResponseModel{id: id}, nil
	}
}

type DeleteBidderRequestModel struct {
	id string
}

func MakeEndpointDeleteBidder(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(DeleteBidderRequestModel)

		if !ok {
			return nil, fmt.Errorf("MakeEndpointDeleteBidder failed parsing request")
		}

		if err = s.DeleteBidder(ctx, req.id); err != nil {
			return nil, fmt.Errorf("MakeEndpointDeleteBidder %w", err)
		}

		return nil, nil
	}
}

type DeleteBiddersRequestModel struct {
	ids []string
}

func MakeEndpointDeleteBidders(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(DeleteBiddersRequestModel)

		if !ok {
			return nil, fmt.Errorf("MakeEndpointDeleteBidder failed parsing request")
		}

		if err = s.DeleteBidders(ctx, req.ids); err != nil {
			return nil, fmt.Errorf("MakeEndpointDeleteBidder %w", err)
		}

		return nil, nil
	}
}
