package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/ireuven89/auctions/bidder-service/bidder"
	"github.com/julienschmidt/httprouter"
)

type Transport struct {
	router *httprouter.Router
	s      Service
}

type Router interface {
	Handle(method, path string, handler http.Handler)
}

func NewTransport(router *httprouter.Router, s Service) Transport {
	transport := Transport{
		router: router,
		s:      s,
	}

	RegisterRoutes(router, s)

	return transport
}

func (t *Transport) ListenAndServe(port string) {
	log.Printf("Starting bidder server on port %s...", port)
	err := http.ListenAndServe(":"+port, t.router)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func RegisterRoutes(router *httprouter.Router, s Service) {

	getBidderHandler := kithttp.NewServer(
		MakeEndpointGetBidder(s),
		decodeGetBidderRequest,
		encodeGetBidderResponse,
	)

	getBiddersHandler := kithttp.NewServer(
		MakeEndpointGetBidders(s),
		decodeGetBiddersRequest,
		encodeGetBiddersResponse,
	)

	createBidderHandler := kithttp.NewServer(
		MakeEndpointCreateBidder(s),
		decodeCreateBidderRequest,
		encodeCreateBidderResponse,
	)

	updateBidderHandler := kithttp.NewServer(
		MakeEndpointUpdateBidder(s),
		decodeUpdateBidderRequest,
		kithttp.EncodeJSONResponse,
	)

	deleteBidderHandler := kithttp.NewServer(
		MakeEndpointDeleteBidder(s),
		decodeDeleteBidderRequest,
		kithttp.EncodeJSONResponse,
	)

	deleteBiddersHandler := kithttp.NewServer(
		MakeEndpointDeleteBidders(s),
		decodeDeleteBiddersRequest,
		kithttp.EncodeJSONResponse,
	)

	router.Handler(http.MethodGet, "/bidders/:id", getBidderHandler)
	router.Handler(http.MethodGet, "/bidders", getBiddersHandler)
	router.Handler(http.MethodPost, "/bidders", createBidderHandler)
	router.Handler(http.MethodPut, "/bidders/:id", updateBidderHandler)
	router.Handler(http.MethodDelete, "/bidders/:id", deleteBidderHandler)
	router.Handler(http.MethodDelete, "/bidders", deleteBiddersHandler)

}

func decodeGetBidderRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	id := httprouter.ParamsFromContext(ctx).ByName("id")

	return GetBidderRequestModel{
		id: id,
	}, nil
}

func encodeGetBidderResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	res, ok := response.(GetBidderResponseModel)

	if !ok {
		return fmt.Errorf("encodeGetBidderResponse failed parsing reponse")
	}

	formmated := formatBidder(res.bidder)

	return json.NewEncoder(w).Encode(formmated)
}

func decodeGetBiddersRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req bidder.BiddersRequest
	name := r.URL.Query().Get("name")

	req.Name = name

	return GetBiddersRequestModel{
		req,
	}, nil
}

func encodeGetBiddersResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	res, ok := response.(GetBiddersResponseModel)

	if !ok {
		return fmt.Errorf("Transport.encodeGetBidderResponse failed parsing reponse")
	}

	var formatted []map[string]interface{}

	for _, bidder := range res.bidders {
		formattedBidder := formatBidder(bidder)
		formatted = append(formatted, formattedBidder)
	}

	return json.NewEncoder(w).Encode(formatted)
}

func decodeCreateBidderRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req bidder.Bidder

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("decodeCreateBidderRequest failed decoding request %v", err)
	}

	return CreateBidderRequestModel{
		bidder: req,
	}, nil
}

func encodeCreateBidderResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	res, ok := response.(CreateBidderResponseModel)

	if !ok {
		return fmt.Errorf("encodeGetBidderResponse failed casting response")
	}

	return json.NewEncoder(w).Encode(map[string]interface{}{"id": res.id})
}

func decodeUpdateBidderRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req bidder.Bidder
	id := httprouter.ParamsFromContext(ctx).ByName("id")

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("decodeUpdateBidderRequest failed decoding request %v", err)
	}

	req.ID = id

	return UpdateBidderRequestModel{
		bidder: req,
	}, nil
}

func decodeDeleteBidderRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	id := httprouter.ParamsFromContext(ctx).ByName("id")

	return DeleteBidderRequestModel{
		id: id,
	}, nil
}

func decodeDeleteBiddersRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var ids []string

	if err = json.NewDecoder(r.Body).Decode(&ids); err != nil {
		return nil, fmt.Errorf("decodeDeleteBiddersRequest failed casting request")
	}

	return DeleteBiddersRequestModel{
		ids: ids,
	}, nil
}
