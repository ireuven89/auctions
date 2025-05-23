package internal

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/ireuven89/auctions/auth-service/user"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/ireuven89/auctions/auth-service/db"
	"github.com/ireuven89/auctions/auth-service/key"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	SignToken(ctx context.Context, user user.User) (string, error)
	Login(ctx context.Context, userIdentifier, password string) (*key.Token, error)
	RefreshToken(ctx context.Context, refreshToken string) (string, error)
	GenerateRefreshToken(ctx context.Context, userInfo string) (string, error)
	GetPublicKey(ctx context.Context) key.JWK
	Register(ctx context.Context, user user.User) (string, string, error)
}

type service struct {
	logger      *zap.Logger
	privateKey  *rsa.PrivateKey
	publicKey   *rsa.PublicKey
	publicKeyId string
	repository  db.Repository
}

const refreshTokenTTL = 24 * 30 * time.Hour
const accessTokenTTL = 15 * time.Minute

func NewAuthService(logger *zap.Logger, repo db.Repository, secretName string) (Service, error) {
	privateKey, err := loadPrivateKeyFromLocal()
	if err != nil {
		return nil, fmt.Errorf("failed starting service %w", err)
	}

	publicKey, err := loadPublicKeyFromLocal()
	return &service{privateKey: privateKey, publicKey: publicKey, logger: logger, repository: repo}, nil
}

func loadPrivateKeyFromSecretsManager(secretName string) (*rsa.PrivateKey, error) {
	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}

	client := secretsmanager.NewFromConfig(cfg)
	out, err := client.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
		SecretId: &secretName,
	})
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode([]byte(*out.SecretString))
	if block == nil {
		return nil, errors.New("invalid PEM block")
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(*out.SecretString))
	if err != nil {
		return nil, fmt.Errorf("loadPrivateKeyFromSecretsManager %w", err)
	}

	return privateKey, nil
}

func loadPrivateKeyFromLocal() (*rsa.PrivateKey, error) {
	privateKeyPath := os.Getenv("JWT_PRIVATE_KEY_PATH")

	privateKeyFile, err := os.ReadFile(privateKeyPath)

	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(privateKeyFile)
	if block == nil {
		return nil, errors.New("invalid PEM block")
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyFile)
	if err != nil {
		return nil, fmt.Errorf("loadPrivateKeyFromSecretsManager %w", err)
	}

	return privateKey, nil
}

func loadPublicKeyFromLocal() (*rsa.PublicKey, error) {
	publicKeyPath := os.Getenv("JWT_PUBLIC_KEY_PATH")

	publicKeyFile, err := os.ReadFile(publicKeyPath)

	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(publicKeyFile)
	if block == nil {
		return nil, errors.New("invalid PEM block")
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyFile)
	if err != nil {
		return nil, fmt.Errorf("loadPrivateKeyFromSecretsManager %w", err)
	}

	return publicKey, nil
}

func loadPublicKeyFromSecretsManager(secretName string) (*rsa.PublicKey, error) {
	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}

	client := secretsmanager.NewFromConfig(cfg)
	out, err := client.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
		SecretId: &secretName,
	})
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode([]byte(*out.SecretString))
	if block == nil {
		return nil, errors.New("invalid PEM block")
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(*out.SecretString))
	if err != nil {
		return nil, fmt.Errorf("loadPrivateKeyFromSecretsManager %w", err)
	}

	return publicKey, nil
}

func (s *service) SignToken(ctx context.Context, user user.User) (string, error) {
	claims := jwt.MapClaims{
		"sub":   user.ID,
		"exp":   time.Now().Add(accessTokenTTL).Unix(),
		"iat":   time.Now().Unix(),
		"email": user.Email,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(s.privateKey)
}

func (s *service) GetPublicKey(ctx context.Context) key.JWK {

	return key.JWK{
		E:   base64.RawURLEncoding.EncodeToString(big.NewInt(int64(s.publicKey.E)).Bytes()),
		N:   base64.RawURLEncoding.EncodeToString((s.publicKey.N).Bytes()),
		Kty: "RSA",
		Alg: jwt.SigningMethodRS256.Name,
		Use: "sig",
		Kid: "default",
	}
}

func (s *service) Register(ctx context.Context, userCredentials user.User) (string, string, error) {
	userID := generateID()

	hashedPassword, err := hashPassword(userCredentials.Password)

	if err != nil {
		return "", "", fmt.Errorf("service.Register %w", err)
	}

	userCredentials.ID = userID
	userCredentials.Password = hashedPassword

	err = s.repository.CreateUser(ctx, userCredentials)

	if err != nil {
		return "", "", fmt.Errorf("service.Register failed %w", err)
	}

	token, err := s.SignToken(ctx, userCredentials)

	if err != nil {
		return "", "", fmt.Errorf("service.Register failed %w", err)
	}

	refreshToken, err := s.GenerateRefreshToken(ctx, userID)

	return token, refreshToken, nil
}

func (s *service) RefreshToken(ctx context.Context, refreshToken string) (string, error) {
	userId, err := s.repository.GetToken(ctx, "refresh:"+refreshToken)
	if err != nil {
		s.logger.Error("service.RefreshToken", zap.Error(err))
		return "", fmt.Errorf("RefreshToken invalid or expired token %w", err)
	}

	user, err := s.repository.FindUser(ctx, userId)

	if err != nil {
		s.logger.Error("service.RefreshToken", zap.Error(err))
		return "", fmt.Errorf("failed fetching user info from refresh token %w", err)
	}

	accessToken, err := s.SignToken(ctx, *user)

	if err != nil {
		s.logger.Error("service.RefreshToken", zap.Error(err))
		return "", fmt.Errorf("RefreshToken failed refreshong token %w", err)
	}

	return accessToken, nil
}

func (s *service) GenerateAccessToken(ctx context.Context, userId string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userId,
		"exp": time.Now().Add(time.Minute * 15).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(s.privateKey)
}

// TODO	change to user info - save in redis user as JSON
func (s *service) GenerateRefreshToken(ctx context.Context, userID string) (string, error) {

	token := uuid.New().String()

	if err := s.repository.SaveRefreshToken(ctx, fmt.Sprintf("refresh:%s", token), userID, refreshTokenTTL); err != nil {
		return "", fmt.Errorf("GenerateRefreshToken %w", err)
	}

	return token, nil
}

func generateID() string {

	return uuid.New().String()
}

func (s *service) Login(ctx context.Context, identifier, password string) (*key.Token, error) {

	user, err := s.repository.FindUserByCredentials(ctx, identifier)

	if err != nil {
		return nil, fmt.Errorf("service.Login user not found result %w", err)
	}

	if ok := verifyUser(user.Password, password); !ok {
		return nil, fmt.Errorf("service.Login invalid password")
	}

	accessToken, err := s.SignToken(ctx, *user)

	if err != nil {
		return nil, fmt.Errorf("service.Login failed creating access token %w", err)
	}

	refreshToken, err := s.GenerateRefreshToken(ctx, user.ID)

	if err != nil {
		return nil, fmt.Errorf("service.Login failed creating refresh token %w", err)

	}

	return &key.Token{
		Access:  accessToken,
		Refresh: refreshToken,
	}, nil
}

func verifyUser(hashedPassword string, password string) bool {

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return false
	}

	return true
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}
