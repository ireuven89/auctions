package internal

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/ireuven89/auctions/auction-service/domain"
	"github.com/ireuven89/auctions/auction-service/internal/service"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/julienschmidt/httprouter"
)

func NewTransport(s service.Service, router *httprouter.Router) Transport {

	transport := Transport{
		router: router,
		s:      s,
	}
	RegisterRoutes(router, s) // Register routes during initialization
	return transport
}

type Transport struct {
	router *httprouter.Router
	s      service.Service
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

func RegisterRoutes(router *httprouter.Router, s service.Service) {

	getAuctionHandler := kithttp.NewServer(
		MakeEndpointGetAuction(s),
		decodeGetAuctionRequest,
		encodeGetAuctionResponse,
		kithttp.ServerErrorEncoder(errorEncoder),
	)

	getAuctionsHandler := kithttp.NewServer(
		MakeEndpointGetAuctions(s),
		decodeGetAuctionsRequest,
		encodeGetAuctionsResponse,
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

	auctionItemsHandler := kithttp.NewServer(
		MakeEndpointCreateAuctionItems(s),
		decodeCreateItemRequest,
		kithttp.EncodeJSONResponse,
		kithttp.ServerErrorEncoder(errorEncoder),
	)

	AuctionItemsPicturesHandler := kithttp.NewServer(
		MakeEndpointCreateAuctionItemPictures(s),
		decodeCreateItemPicturesRequest,
		kithttp.EncodeJSONResponse,
		kithttp.ServerErrorEncoder(errorEncoder),
	)

	router.Handler(http.MethodGet, "/auctions/:id", getAuctionHandler)
	router.Handler(http.MethodGet, "/auctions", getAuctionsHandler)
	router.Handler(http.MethodPost, "/auctions", createAuctionHandler)
	router.Handler(http.MethodPut, "/auctions/:id", updateAuctionHandler)
	router.Handler(http.MethodDelete, "/auctions/:id", deleteAuctionHandler)
	router.Handler(http.MethodPost, "/auctions/:id/items", auctionItemsHandler)
	router.Handler(http.MethodPost, "/auctions/:id/items/:itemId/pictures", AuctionItemsPicturesHandler)

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

func decodeGetAuctionsRequest(c context.Context, r *http.Request) (interface{}, error) {

	name := r.URL.Query().Get("name")

	return GetAuctionsRequestModel{
		AuctionRequest: domain.AuctionRequest{
			Description: name,
		},
	}, nil
}

func encodeGetAuctionsResponse(c context.Context, w http.ResponseWriter, response interface{}) error {
	res, ok := response.(GetAuctionsResponseModel)

	if !ok {
		return fmt.Errorf("encodeGetAuctionResponse failed parsing reponse")
	}

	var formatted []map[string]interface{}

	for _, auct := range res.auctions {
		formatted = append(formatted, formatAuction(&auct))
	}

	return json.NewEncoder(w).Encode(&formatted)
}

func decodeCreateAuctionRequest(c context.Context, r *http.Request) (interface{}, error) {
	var req CreateAuctionRequestModel

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Printf("decodeCreateAuctionRequest failed decoding request %v", err)
		return nil, fmt.Errorf("decodeCreateAuctionRequest failed casting request %w", err)
	}

	return req, nil
}

func decodeCreateItemRequest(c context.Context, r *http.Request) (interface{}, error) {
	var items []domain.Item
	json.NewDecoder(r.Body)
	return AuctionItemsRequestModel{
		AuctionID: r.PathValue("id"),
		items:     items,
	}, nil
}

func decodeCreateItemPicturesRequest(c context.Context, r *http.Request) (interface{}, error) {

	if err := r.ParseMultipartForm(20 << 20); err != nil {
		return nil, fmt.Errorf("decodeCreateItemPicturesRequest %w", err)
	}

	files := r.MultipartForm.File["images"]

	if len(files) > 5 {
		return nil, fmt.Errorf("decodeCreateAuctionRequest images upload is restircted to 5 images")
	}

	return AuctionPicturesRequestModel{
		Files: files,
	}, nil
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
		return nil, fmt.Errorf("decodeUpdateAuctionRequest %w", err)
	}

	req.ID = httprouter.ParamsFromContext(c).ByName("id")

	return req, nil
}

func decodeDeleteAuctionRequest(c context.Context, r *http.Request) (interface{}, error) {
	var req DeleteAuctionRequestModel

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Printf("decodeDeleteAuctionRequest failed decoding request %v", err)
		return nil, fmt.Errorf("decodeDeleteAuctionRequest %w", err)
	}

	return req, nil
}

func errorEncoder(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")

	switch {
	case errors.Is(err, domain.ErrNotFound):
		w.WriteHeader(http.StatusNotFound)
	case errors.Is(err, domain.ErrTooManyRequests):
		w.WriteHeader(http.StatusTooManyRequests)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(map[string]string{
		"error": err.Error(),
	})
}
