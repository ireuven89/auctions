package internal

import (
	"context"
	"errors"
	"fmt"
	"github.com/ireuven89/auctions/auth-service/user"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/ireuven89/auctions/auth-service/internal/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func setupTestKeys(t *testing.T) (privPath, pubPath string) {
	t.Helper()
	privPath = filepath.Join("testdata", "test_private.pem")
	pubPath = filepath.Join("testdata", "test_public.pem")
	if _, err := os.Stat(privPath); err != nil {
		t.Fatalf("private key missing: %v", err)
	}
	if _, err := os.Stat(pubPath); err != nil {
		t.Fatalf("public key missing: %v", err)
	}
	os.Setenv("JWT_PRIVATE_KEY_PATH", privPath)
	os.Setenv("JWT_PUBLIC_KEY_PATH", pubPath)
	return
}

/*func TestNewAuthService(t *testing.T) {
	setupTestKeys(t)
	logger := zap.NewNop()
	repo := &mocks.MockRepo{}

	svc, err := NewAuthService(logger, repo, "ignored")
	assert.NoError(t, err)
	assert.NotNil(t, svc)
}

func TestNewAuthService_PrivateKeyFail(t *testing.T) {
	os.Setenv("JWT_PRIVATE_KEY_PATH", "nonexistent.pem")
	os.Setenv("JWT_PUBLIC_KEY_PATH", "nonexistent_pub.pem")
	logger := zap.NewNop()
	repo := &mocks.MockRepo{}
	svc, err := NewAuthService(logger, repo, "ignored")
	assert.Error(t, err)
	assert.Nil(t, svc)
}*/

func TestNewAuthService_PublicKeyFail(t *testing.T) {
	// Write only private key (with invalid data)
	tmpPriv := "test_private.pem"
	defer os.Remove(tmpPriv)
	os.Setenv("JWT_PRIVATE_KEY_PATH", tmpPriv)
	os.Setenv("JWT_PUBLIC_KEY_PATH", "nonexistent_pub.pem")
	assert.NoError(t, os.WriteFile(tmpPriv, []byte("BAD DATA"), 0600))
	logger := zap.NewNop()
	repo := &mocks.MockRepo{}
	svc, err := NewAuthService(logger, repo, "ignored")
	assert.Error(t, err)
	assert.Nil(t, svc)
}

/*// Register success
func TestService_Register_Success(t *testing.T) {
	setupTestKeys(t)
	//	defer os.Remove(privateKeyPath)

	logger := zap.NewNop()
	repo := &mocks.MockRepo{
		CreateUserFunc:       func(ctx context.Context, user user.User) error { return nil },
		SaveRefreshTokenFunc: func(ctx context.Context, key, userId string, ttl time.Duration) error { return nil },
	}

	// Example: If NewAuthService takes key path as param
	svc, err := NewAuthService(logger, repo, "")
	assert.NoError(t, err)
	_, _, err = svc.Register(context.Background(), user.User{Email: "foo@bar.com", Password: "pass"})
	assert.NoError(t, err)
}*/

// Register fails on CreateUser
func TestService_Register_CreateUserFail(t *testing.T) {
	setupTestKeys(t)
	logger := zap.NewNop()
	repo := &mocks.MockRepo{
		CreateUserFunc: func(ctx context.Context, user user.User) error { return errors.New("fail create") },
	}
	svc, err := NewAuthService(logger, repo, "ignored")
	assert.NoError(t, err)
	_, _, err = svc.Register(context.Background(), user.User{Email: "foo@bar.com", Password: "pass"})
	assert.Error(t, err)
}

// RefreshToken fails if token not found
func TestService_RefreshToken_TokenNotFound(t *testing.T) {
	setupTestKeys(t)
	logger := zap.NewNop()
	repo := &mocks.MockRepo{
		GetTokenFunc:       func(ctx context.Context, key string) (string, error) { return "", errors.New("not found") },
		GetRefreshRateFunc: func(ctx context.Context, token string) (int, error) { return 1, nil },
	}
	svc, err := NewAuthService(logger, repo, "ignored")
	assert.NoError(t, err)
	_, err = svc.RefreshToken(context.Background(), "badtoken")
	assert.Error(t, err)
}

type EmailTest struct {
	pattern string
	valid   bool
}

func TestRegex(t *testing.T) {
	testEmails := []EmailTest{
		{
			valid:   true,
			pattern: "test@example.com",
		}, {
			valid:   false,
			pattern: "invalid-email.com",
		},
		{
			valid:   true,
			pattern: "user@domain.org",
		},
		{
			valid:   false,
			pattern: "another@domain",
		},
		{
			valid:   true,
			pattern: "foo@bar.baz",
		},
		{
			valid:   true,
			pattern: fmt.Sprintf("test_%d@example.com", time.Now().UnixNano()),
		},
	}

	for _, test := range testEmails {
		assert.Equal(t, test.valid, validateEmail(test.pattern))
	}
}
