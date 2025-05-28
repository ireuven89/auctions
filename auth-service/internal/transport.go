package internal

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/ireuven89/auctions/auth-service/key"

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
		kithttp.ServerErrorEncoder(errorEncoder),
	)

	loginHandler := kithttp.NewServer(
		MakeEndpointLogin(s),
		decodeLoginRequest,
		encodeLoginUserResponse,
		kithttp.ServerErrorEncoder(errorEncoder),
	)

	refreshHandler := kithttp.NewServer(
		MakeEndpointRefreshToken(s),
		decodeRefreshRequest,
		encodeRefreshResponse,
		kithttp.ServerErrorEncoder(errorEncoder),
	)

	logoutHandler := kithttp.NewServer(
		MakeEndpointLogout(s),
		decodeRegisterUserRequest,
		encodeRegisterUserResponse,
		kithttp.ServerErrorEncoder(errorEncoder),
	)

	publicKeyHandler := kithttp.NewServer(
		MakeEndpointGetPublicKey(s),
		decodeGetPublicRequest,
		encodeGetPublicResponse,
		kithttp.ServerErrorEncoder(errorEncoder),
	)

	router.Handler(http.MethodPost, "/auth/register", registerUserHandler)
	router.Handler(http.MethodPost, "/auth/login", loginHandler)
	router.Handler(http.MethodPost, "/auth/refresh", refreshHandler)
	router.Handler(http.MethodPost, "/auth/logout", logoutHandler)
	router.Handler(http.MethodGet, "/auth/jwks", publicKeyHandler)
	router.Handler(http.MethodDelete, "/auth/user/:id", publicKeyHandler)

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

func decodeLogoutRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req LoginRequestModel

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("decodeLoginRequest failed parsing request %w", err)
	}

	return req, nil
}

func encodeLogoutUserResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
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
		"jwks": res.PublicKey.Keys,
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

func errorEncoder(ctx context.Context, err error, w http.ResponseWriter) {
	switch {
	case errors.Is(err, key.ErrUserNotFound),
		errors.Is(err, key.ErrInvalidCredentials):
		w.WriteHeader(http.StatusUnauthorized) // 401

	case errors.Is(err, key.ErrInvalidToken),
		errors.Is(err, key.ErrExpiredToken):
		w.WriteHeader(http.StatusUnauthorized) // 401

	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})

}
