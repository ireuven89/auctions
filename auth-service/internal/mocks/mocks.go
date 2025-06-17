package mocks

import (
	"context"
	"time"

	"github.com/ireuven89/auctions/auth-service/key"

	"github.com/ireuven89/auctions/auth-service/user"
)

// MockRepository mocks the repository interface
type MockRepo struct {
	CreateUserFunc            func(ctx context.Context, u user.User) error
	FindUserFunc              func(ctx context.Context, id string) error
	FindUserByCredentialsFunc func(ctx context.Context, identifier string) (*user.User, error)
	GetTokenFunc              func(ctx context.Context, token string) (string, error)
	SaveRefreshTokenFunc      func(ctx context.Context, token string, userId string, ttl time.Duration) error
	GetRefreshRateFunc        func(ctx context.Context, token string) (int, error)
}

func (m *MockRepo) GetRefreshRate(ctx context.Context, token string) (int, error) {
	return m.GetRefreshRateFunc(ctx, token)
}

func (m *MockRepo) FindUser(ctx context.Context, id string) (*user.User, error) {
	return m.FindUser(ctx, id)
}

func (m *MockRepo) FindUserByCredentials(ctx context.Context, identifier string) (*user.User, error) {
	return m.FindUserByCredentials(ctx, identifier)
}

func (m *MockRepo) SaveRefreshToken(ctx context.Context, token string, userId string, ttl time.Duration) error {
	return m.SaveRefreshTokenFunc(ctx, token, userId, ttl)
}

func (m *MockRepo) GetToken(ctx context.Context, token string) (string, error) {
	return m.GetToken(ctx, token)
}

func (m *MockRepo) CreateUser(ctx context.Context, u user.User) error {
	return m.CreateUserFunc(ctx, u)
}

// MockService embeds service.Service and mocks token functions
type MockService struct {
	PubKey key.JWK
	MockRepo
	signTokenFunc        func(ctx context.Context, u user.User) (string, error)
	generateRefreshToken func(ctx context.Context, id string) (string, error)
	LoginFunc            func(ctx context.Context, userIdentifier, password string) (*key.Token, error)
	RefreshTokenFunc     func(ctx context.Context, refreshToken string) (string, error)
	GetPublicKeyFunc     func(ctx context.Context) key.JWK
	RegisterFunc         func(ctx context.Context, user user.User) (string, string, error)
}

func (m *MockService) SignToken(ctx context.Context, u user.User) (string, error) {
	return m.signTokenFunc(ctx, u)
}
func (m *MockService) GenerateRefreshToken(ctx context.Context, id string) (string, error) {
	return m.generateRefreshToken(ctx, id)
}

func (m *MockService) Login(ctx context.Context, userIdentifier, password string) (*key.Token, error) {
	return m.LoginFunc(ctx, userIdentifier, password)
}
func (m *MockService) RefreshToken(ctx context.Context, refreshToken string) (string, error) {
	return m.RefreshTokenFunc(ctx, refreshToken)
}
func (m *MockService) GetPublicKey(ctx context.Context) key.JWK {
	return m.GetPublicKeyFunc(ctx)
}

func (m *MockService) Register(ctx context.Context, user user.User) (string, string, error) {
	return m.RegisterFunc(ctx, user)
}
