package internal

import (
	"context"
	"errors"
	"fmt"

	"github.com/ireuven89/auctions/auth-service/user"

	"github.com/go-kit/kit/endpoint"
	"github.com/ireuven89/auctions/auth-service/key"
)

type GetPublicKeyResponse struct {
	publicKey key.JWK
}

func MakeEndpointGetPublicKey(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {

		publicKey := s.GetPublicKey(ctx)

		return GetPublicKeyResponse{
			publicKey: publicKey,
		}, nil
	}
}

type RegisterUserRequest struct {
	user user.User
}

type RegisterUserResponse struct {
	AccessToken  string
	RefreshToken string
}

func MakeEndpointRegisterUser(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(RegisterUserRequest)

		if !ok {
			return nil, fmt.Errorf("MakeEndpointRegisterUser failed parsing request")
		}

		accessToken, refreshToken, err := s.Register(ctx, req.user)

		if err != nil {
			return nil, fmt.Errorf("MakeEndpointRegisterUser %w", err)
		}

		return RegisterUserResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}, nil
	}
}

type LoginRequestModel struct {
	Identifier string
	Password   string
}

type LoginResponseModel struct {
	AccessToken  string
	RefreshToken string
}

func MakeEndpointLogin(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(LoginRequestModel)

		if !ok {
			return nil, fmt.Errorf("MakeEndpointLogin failed casting request")
		}

		token, err := s.Login(ctx, req.Identifier, req.Password)

		if err != nil {
			// Pass through ErrUnauthorized for the transport to handle as 401
			if errors.Is(err, key.ErrInvalidCredentials) || errors.Is(err, key.ErrUserNotFound) {
				return nil, err
			}
			// Wrap all other errors
			return nil, fmt.Errorf("MakeEndpointLogin: %w", err)
		}

		return LoginResponseModel{
			AccessToken:  token.Access,
			RefreshToken: token.Refresh,
		}, nil
	}
}

func MakeEndpointLogout(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(LoginRequestModel)

		if !ok {
			return nil, fmt.Errorf("MakeEndpointLogout failed casting request")
		}
		token, err := s.Login(ctx, req.Identifier, req.Password)

		if err != nil {
			return nil, fmt.Errorf("MakeEndpointLogout.failed  %w", err)
		}

		return token, nil
	}
}

type RefreshRequestModel struct {
	Refresh string
}

type RefreshResponseModel struct {
	AccessToken string
}

func MakeEndpointRefreshToken(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(RefreshRequestModel)

		if !ok {
			return nil, fmt.Errorf("MakeEndpointRefreshToken failed casting request")
		}
		accessToken, err := s.RefreshToken(ctx, req.Refresh)

		if err != nil {
			return nil, fmt.Errorf("MakeEndpointRefreshToken %w", err)
		}

		return RefreshResponseModel{
			AccessToken: accessToken,
		}, nil
	}
}
