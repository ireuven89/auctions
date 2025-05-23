package mocks

import (
	"context"
	"time"

	"github.com/ireuven89/auctions/auth-service/internal"
	"github.com/ireuven89/auctions/auth-service/user"
)

// MockRepository mocks the repository interface
type MockRepo struct {
	CreateUserFunc            func(ctx context.Context, u user.User) error
	FindUserFunc              func(ctx context.Context, id string) error
	FindUserByCredentialsFunc func(ctx context.Context, identifier string) (*user.User, error)
	GetTokenFunc              func(ctx context.Context, token string) (string, error)
}

func (m *MockRepo) FindUser(ctx context.Context, id string) (*user.User, error) {
	return m.FindUser(ctx, id)
}

func (m *MockRepo) FindUserByCredentials(ctx context.Context, identifier string) (*user.User, error) {
	return m.FindUserByCredentials(ctx, identifier)
}

func (m *MockRepo) SaveRefreshToken(ctx context.Context, token string, userInfo []byte, ttl time.Duration) error {
	return m.SaveRefreshToken(ctx, token, userInfo, ttl)
}

func (m *MockRepo) GetToken(ctx context.Context, token string) (string, error) {
	return m.GetToken(ctx, token)
}

func (m *MockRepo) CreateUser(ctx context.Context, u user.User) error {
	return m.CreateUserFunc(ctx, u)
}

// MockService embeds service.Service and mocks token functions
type MockService struct {
	internal.Service // Embed actual service if needed
	MockRepo
	signTokenFunc         func(ctx context.Context, u user.User) (string, error)
	generateRefreshTokenF func(ctx context.Context, id string) (string, error)
}

func (m *MockService) SignToken(ctx context.Context, u user.User) (string, error) {
	return m.signTokenFunc(ctx, u)
}
func (m *MockService) GenerateRefreshToken(ctx context.Context, id string) (string, error) {
	return m.generateRefreshTokenF(ctx, id)
}
