package internal

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/ireuven89/auctions/auction-service/auction"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

var errNotFound = errors.New("not found")

func NewTransport(s Service, router *httprouter.Router) Transport {

	transport := Transport{
		router: router,
		s:      s,
	}
	RegisterRoutes(router, s) // Register routes during initialization
	return transport
}

type Transport struct {
	router *httprouter.Router
	s      Service
}

type Router interface {
	Handle(method, path string, handler http.Handler)
}

func (t *Transport) ListenAndServe(port string) {
	log.Printf("Starting auctions server on port %s...", port)
	err := http.ListenAndServe(":"+port, t.router)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func RegisterRoutes(router *httprouter.Router, s Service) {

	getAuctionHandler := kithttp.NewServer(
		MakeEndpointGetAuction(s),
		decodeGetAuctionRequest,
		encodeGetAuctionResponse,
		kithttp.ServerErrorEncoder(errorEncoder),
	)

	createAuctionHandler := kithttp.NewServer(
		MakeEndpointCreateAuction(s),
		decodeCreateAuctionRequest,
		encodeCreateAuctionResponse,
		kithttp.ServerErrorEncoder(errorEncoder),
	)

	updateAuctionHandler := kithttp.NewServer(
		MakeEndpointUpdateAuction(s),
		decodeUpdateAuctionRequest,
		kithttp.EncodeJSONResponse,
		kithttp.ServerErrorEncoder(errorEncoder),
	)

	deleteAuctionHandler := kithttp.NewServer(
		MakeEndpointDeleteAuction(s),
		decodeDeleteAuctionRequest,
		kithttp.EncodeJSONResponse,
		kithttp.ServerErrorEncoder(errorEncoder),
	)

	router.Handler(http.MethodGet, "/auctions/:id", getAuctionHandler)
	router.Handler(http.MethodPost, "/auctions", createAuctionHandler)
	router.Handler(http.MethodPut, "/auctions/:id", updateAuctionHandler)
	router.Handler(http.MethodDelete, "/auctions/:id", deleteAuctionHandler)
}

func decodeGetAuctionRequest(c context.Context, r *http.Request) (interface{}, error) {
	id := httprouter.ParamsFromContext(c).ByName("id")
	return GetAuctionRequestModel{
		id: id,
	}, nil
}

func encodeGetAuctionResponse(c context.Context, w http.ResponseWriter, response interface{}) error {
	res, ok := response.(GetAuctionResponseModel)

	if !ok {
		return fmt.Errorf("encodeGetAuctionResponse failed parsing reponse")
	}

	formatted := formatAuction(res.auction)

	return json.NewEncoder(w).Encode(&formatted)
}

func decodeCreateAuctionRequest(c context.Context, r *http.Request) (interface{}, error) {
	var req CreateAuctionRequestModel

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Printf("decodeCreateAuctionRequest failed decoding request %v", err)
		return nil, err
	}

	return req, nil
}

func encodeCreateAuctionResponse(c context.Context, w http.ResponseWriter, response interface{}) error {
	res, ok := response.(CreateAuctionResponseModel)

	if !ok {
		return fmt.Errorf("encodeCreateAuctionResponse failed parsing response")
	}

	formatted := map[string]interface{}{"id": res.id}

	return json.NewEncoder(w).Encode(formatted)
}

func decodeUpdateAuctionRequest(c context.Context, r *http.Request) (interface{}, error) {
	var req UpdateAuctionRequestModel

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Printf("decodeCreateAuctionRequest failed decoding request %v", err)
		return nil, err
	}

	req.ID = httprouter.ParamsFromContext(c).ByName("id")

	return req, nil
}

func decodeDeleteAuctionRequest(c context.Context, r *http.Request) (interface{}, error) {
	var req DeleteAuctionRequestModel

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Printf("decodeDeleteAuctionRequest failed decoding request %v", err)
		return nil, err
	}

	return req, nil
}

func errorEncoder(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")

	switch {
	case errors.Is(err, auction.ErrNotFound):
		w.WriteHeader(http.StatusNotFound)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(map[string]string{
		"error": err.Error(),
	})
}
