package internal

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/ireuven89/auctions/shared/jwksprovider"
	"testing"

	"github.com/ireuven89/auctions/auth-service/internal/mocks"
	"github.com/ireuven89/auctions/auth-service/key"
	user2 "github.com/ireuven89/auctions/auth-service/user"
	"github.com/stretchr/testify/assert"
)

// GET PUBLIC KEY
func TestMakeEndpointGetPublicKey(t *testing.T) {
	mock := &mocks.MockService{GetPublicKeyFunc: func(ctx context.Context) jwksprovider.JWKS {
		return jwksprovider.JWKS{
			Keys: []json.RawMessage{[]byte{'E'}},
		}
	}}
	endpoint := MakeEndpointGetPublicKey(mock)

	resp, err := endpoint(context.Background(), nil)
	assert.NoError(t, err)
	assert.Equal(t, []json.RawMessage{[]byte{'E'}}, resp.(GetPublicKeyResponse).PublicKey)
}

// REGISTER USER
func TestMakeEndpointRegisterUser_Success(t *testing.T) {
	mock := &mocks.MockService{
		RegisterFunc: func(ctx context.Context, user user2.User) (string, string, error) {
			return "acc", "ref", nil
		},
	}
	endpoint := MakeEndpointRegisterUser(mock)

	req := RegisterUserRequest{user: user2.User{Name: "foo", Password: "bar"}}
	resp, err := endpoint(context.Background(), req)
	assert.NoError(t, err)
	r := resp.(RegisterUserResponse)
	assert.Equal(t, "acc", r.AccessToken)
	assert.Equal(t, "ref", r.RefreshToken)
}

func TestMakeEndpointRegisterUser_Error(t *testing.T) {
	mock := &mocks.MockService{
		RegisterFunc: func(ctx context.Context, user user2.User) (string, string, error) {
			return "", "", errors.New("register error")
		},
	}
	endpoint := MakeEndpointRegisterUser(mock)

	req := RegisterUserRequest{user: user2.User{Name: "foo"}}
	resp, err := endpoint(context.Background(), req)
	assert.Nil(t, resp)
	assert.Error(t, err)
}

func TestMakeEndpointRegisterUser_BadRequest(t *testing.T) {
	mock := &mocks.MockService{}
	endpoint := MakeEndpointRegisterUser(mock)

	resp, err := endpoint(context.Background(), "not a request")
	assert.Nil(t, resp)
	assert.Error(t, err)
}

// LOGIN
func TestMakeEndpointLogin_Success(t *testing.T) {
	mock := &mocks.MockService{
		LoginFunc: func(ctx context.Context, identifier, password string) (*key.Token, error) {
			return &key.Token{Access: "acc", Refresh: "ref"}, nil
		},
	}
	endpoint := MakeEndpointLogin(mock)

	req := LoginRequestModel{Identifier: "foo", Password: "bar"}
	resp, err := endpoint(context.Background(), req)
	assert.NoError(t, err)
	r := resp.(LoginResponseModel)
	assert.Equal(t, "acc", r.AccessToken)
	assert.Equal(t, "ref", r.RefreshToken)
}

func TestMakeEndpointLogin_Error(t *testing.T) {
	mock := &mocks.MockService{
		LoginFunc: func(ctx context.Context, identifier, password string) (*key.Token, error) {
			return &key.Token{}, errors.New("login error")
		},
	}
	endpoint := MakeEndpointLogin(mock)

	req := LoginRequestModel{Identifier: "foo", Password: "bar"}
	resp, err := endpoint(context.Background(), req)
	assert.Nil(t, resp)
	assert.Error(t, err)
}

func TestMakeEndpointLogin_BadRequest(t *testing.T) {
	mock := &mocks.MockService{}
	endpoint := MakeEndpointLogin(mock)

	resp, err := endpoint(context.Background(), 42)
	assert.Nil(t, resp)
	assert.Error(t, err)
}

// LOGOUT
/*func TestMakeEndpointLogout_Success(t *testing.T) {
	mock := &mocks.MockService{
		LoginFunc: func(ctx context.Context, identifier, password string) (*key.Token, error) {
			return &key.Token{Access: "acc", Refresh: "ref"}, nil
		},
	}
	endpoint := MakeEndpointLogout(mock)

	req := LoginRequestModel{Identifier: "foo", Password: "bar"}
	resp, err := endpoint(context.Background(), req)
	assert.NoError(t, err)
	r := resp.(key.Token)
	assert.Equal(t, "acc", r.Access)
	assert.Equal(t, "ref", r.Refresh)
}

func TestMakeEndpointLogout_Error(t *testing.T) {
	mock := &mocks.MockService{
		LO: func(ctx context.Context, identifier, password string) (Token, error) {
			return Token{}, errors.New("logout error")
		},
	}
	endpoint := MakeEndpointLogout(mock)

	req := LoginRequestModel{Identifier: "foo", Password: "bar"}
	resp, err := endpoint(context.Background(), req)
	assert.Nil(t, resp)
	assert.Error(t, err)
}*/

func TestMakeEndpointLogout_BadRequest(t *testing.T) {
	mock := &mocks.MockService{}
	endpoint := MakeEndpointLogout(mock)

	resp, err := endpoint(context.Background(), 123)
	assert.Nil(t, resp)
	assert.Error(t, err)
}

// REFRESH TOKEN
func TestMakeEndpointRefreshToken_Success(t *testing.T) {
	mock := &mocks.MockService{
		RefreshTokenFunc: func(ctx context.Context, refresh string) (string, error) {
			return "new-access", nil
		},
	}
	endpoint := MakeEndpointRefreshToken(mock)

	req := RefreshRequestModel{Refresh: "ref"}
	resp, err := endpoint(context.Background(), req)
	assert.NoError(t, err)
	r := resp.(RefreshResponseModel)
	assert.Equal(t, "new-access", r.AccessToken)
}

func TestMakeEndpointRefreshToken_Error(t *testing.T) {
	mock := &mocks.MockService{
		RefreshTokenFunc: func(ctx context.Context, refresh string) (string, error) {
			return "", errors.New("refresh error")
		},
	}
	endpoint := MakeEndpointRefreshToken(mock)

	req := RefreshRequestModel{Refresh: "ref"}
	resp, err := endpoint(context.Background(), req)
	assert.Nil(t, resp)
	assert.Error(t, err)
}

func TestMakeEndpointRefreshToken_BadRequest(t *testing.T) {
	mock := &mocks.MockService{}
	endpoint := MakeEndpointRefreshToken(mock)

	resp, err := endpoint(context.Background(), 1)
	assert.Nil(t, resp)
	assert.Error(t, err)
}
