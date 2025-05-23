package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/ireuven89/auctions/auth-service/user"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/julienschmidt/httprouter"
)

func NewTransport(router *httprouter.Router, s Service) Transport {

	RegisterRoutes(router, s)
	return Transport{
		router: router,
		s:      s,
	}
}

type Transport struct {
	router *httprouter.Router
	s      Service
}

func (t *Transport) ListenAndServe(port string) {
	log.Printf("Starting auth server on port %s...", port)
	err := http.ListenAndServe(":"+port, t.router)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func RegisterRoutes(router *httprouter.Router, s Service) {

	registerUserHandler := kithttp.NewServer(
		MakeEndpointRegisterUser(s),
		decodeRegisterUserRequest,
		encodeRegisterUserResponse,
	)

	loginHandler := kithttp.NewServer(
		MakeEndpointLogin(s),
		decodeLoginRequest,
		encodeLoginUserResponse,
	)

	refreshHandler := kithttp.NewServer(
		MakeEndpointRefreshToken(s),
		decodeRefreshRequest,
		encodeRefreshResponse,
	)

	logoutHandler := kithttp.NewServer(
		MakeEndpointLogout(s),
		decodeRegisterUserRequest,
		encodeRegisterUserResponse,
	)

	publicKeyHandler := kithttp.NewServer(
		MakeEndpointGetPublicKey(s),
		decodeGetPublicRequest,
		encodeGetPublicResponse,
	)

	router.Handler(http.MethodPost, "/auth/register", registerUserHandler)
	router.Handler(http.MethodPost, "/auth/login", loginHandler)
	router.Handler(http.MethodPost, "/auth/refresh", refreshHandler)
	router.Handler(http.MethodPost, "/auth/logout", logoutHandler)
	router.Handler(http.MethodGet, "/auth/publicKey", publicKeyHandler)
}

func decodeRegisterUserRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var userInfo user.User

	if err := json.NewDecoder(r.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("decodeRegisterUserRequest failed parsing request %w", err)
	}

	return RegisterUserRequest{
		user: userInfo,
	}, nil
}

func encodeRegisterUserResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	res, ok := response.(RegisterUserResponse)

	if !ok {
		return fmt.Errorf("encodeRegisterUserResponse.failed casting response")
	}

	formatted := map[string]interface{}{
		"accessToken":  res.AccessToken,
		"refreshToken": res.RefreshToken,
	}

	w.Header().Set("Content-Type", "application/json")

	return json.NewEncoder(w).Encode(&formatted)
}

func decodeLoginRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req LoginRequestModel

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("decodeLoginRequest failed parsing request %w", err)
	}

	return req, nil
}

func encodeLoginUserResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	res, ok := response.(LoginResponseModel)

	if !ok {
		return fmt.Errorf("encodeLoginUserResponse.failed casting response")
	}

	formatted := map[string]interface{}{
		"token":        res.AccessToken,
		"refreshToken": res.RefreshToken,
	}

	w.Header().Set("Content-Type", "application/json")

	return json.NewEncoder(w).Encode(&formatted)
}

func decodeGetPublicRequest(ctx context.Context, r *http.Request) (interface{}, error) {

	return nil, nil
}

func encodeGetPublicResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	res, ok := response.(GetPublicKeyResponse)

	if !ok {
		return fmt.Errorf("encodeGetPublicResponse failed casting response")
	}

	formatted := map[string]interface{}{
		"jwks": res.publicKey,
	}

	w.Header().Set("Content-Type", "application/json")

	return json.NewEncoder(w).Encode(&formatted)
}

func decodeRefreshRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var refreshRequest RefreshRequestModel

	if err := json.NewDecoder(r.Body).Decode(&refreshRequest); err != nil {
		return nil, fmt.Errorf("decodeGetRefreshRequest failed encdoing request %w", err)
	}

	return refreshRequest, nil
}

func encodeRefreshResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	refreshResponse, ok := response.(RefreshResponseModel)

	if !ok {
		return fmt.Errorf("encodeRefreshResponse failed casting response")
	}

	formatted := map[string]interface{}{
		"token": refreshResponse.AccessToken,
	}

	w.Header().Set("Content-Type", "application/json")

	return json.NewEncoder(w).Encode(formatted)
}
