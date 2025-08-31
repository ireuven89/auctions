package internal

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"sync"
	"time"

	"github.com/ireuven89/auctions/shared/jwksprovider"

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
	GetPublicKey(ctx context.Context) jwksprovider.JWKS
	Register(ctx context.Context, user user.User) (string, string, error)
}

type service struct {
	logger       *zap.Logger
	privateKey   *rsa.PrivateKey
	mu           sync.RWMutex
	publicKey    jwksprovider.JWKS
	RotateTicker *time.Ticker
	KeyMutex     sync.RWMutex
	repository   db.Repository
}

const refreshTokenTTL = 24 * 30 * time.Hour
const refreshMaxRate = 3
const accessTokenTTL = 15 * time.Minute

func NewAuthService(logger *zap.Logger, repo db.Repository, secretName string) (Service, error) {

	privateKey, err := loadPrivateKeyFromLocal()
	if err != nil {
		return nil, fmt.Errorf("failed starting service %w", err)
	}
	publicKey, err := loadPublicKeyFromLocal()

	if err != nil {
		return nil, fmt.Errorf("failed starting service %w", err)
	}

	s := service{privateKey: privateKey, publicKey: generateJWKSFromPublicKey(publicKey), logger: logger, repository: repo, RotateTicker: time.NewTicker(10 * time.Minute)}

	//todo remove this when key rotation is implmented in shared
	//#	go s.startKeyRotation()

	return &s, nil
}

// TODO this should called in NewAuthService  when secret key is done
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

// TODO this should be removed when secret key is done
func (s *service) startPrivateKeyRefresher(interval time.Duration) {
	go func() {
		for {
			key, err := loadPrivateKeyFromSecretsManager("")
			if err != nil {
				log.Printf("Failed to reload private key: %v", err)
			} else {
				s.mu.Lock()
				s.privateKey = key
				s.mu.Unlock()
				log.Println("Private key reloaded successfully")
			}
			time.Sleep(interval)
		}
	}()
}

func loadPrivateKeyFromLocal() (*rsa.PrivateKey, error) {
	privateKeyPath := os.Getenv("JWT_PRIVATE_KEY_PATH")

	privateKeyFile, err := os.ReadFile(privateKeyPath)

	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(privateKeyFile)

	if block == nil || block.Type != "PRIVATE KEY" {
		return nil, errors.New("invalid PEM block")
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyFile)
	if err != nil {
		return nil, fmt.Errorf("loadPrivateKeyFromLocal %w", err)
	}

	return privateKey, nil
}

// TODO this should be removed when secret key is done
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
		return nil, fmt.Errorf("loadPublicKeyFromLocal %w", err)
	}

	return publicKey, nil
}

// TODO this should called in NewAuthService when secret key is done
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

func (s *service) SignToken(ctx context.Context, userInfo user.User) (string, error) {
	claims := jwt.MapClaims{
		"sub":   userInfo.ID,
		"exp":   time.Now().Add(accessTokenTTL).Unix(),
		"iat":   time.Now().Unix(),
		"email": userInfo.Email,
	}

	jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	return token.SignedString(s.privateKey)
}

func (s *service) GetPublicKey(ctx context.Context) jwksprovider.JWKS {

	return s.publicKey
}

func (s *service) startKeyRotation() {
	for range s.RotateTicker.C {
		newKey, _ := rsa.GenerateKey(rand.Reader, 2048)
		s.KeyMutex.Lock()
		s.privateKey = newKey
		s.publicKey = generateJWKSFromPublicKey(&newKey.PublicKey)
		s.KeyMutex.Unlock()
		log.Println("🔄 AuthService rotated RSA key")
	}
}

// generateJWKSFromPublicKey - this method exposes on ly the public key
func generateJWKSFromPublicKey(pub *rsa.PublicKey) jwksprovider.JWKS {
	jwk := map[string]interface{}{
		"kty": "RSA",
		"kid": "current",
		"n":   pub.N.String(),
		"e":   "AQAB",
	}
	raw, _ := json.Marshal(jwk)
	return jwksprovider.JWKS{
		Keys: []json.RawMessage{raw},
	}
}

func (s *service) Register(ctx context.Context, userCredentials user.User) (string, string, error) {
	if ok := validateEmail(userCredentials.Email); !ok {
		return "", "", fmt.Errorf("invalid email pattern")
	}

	hashedPassword, err := hashPassword(userCredentials.Password)

	if err != nil {
		return "", "", fmt.Errorf("service.Register %w", err)
	}

	userID := generateID()

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

func validateEmail(email string) bool {
	matched, err := regexp.MatchString("^[\\w0-9]+@[\\w0-9]+\\.[a-z]{2,}$", email)

	if err != nil {
		return false
	}

	return matched
}

func (s *service) RefreshToken(ctx context.Context, refreshToken string) (string, error) {
	if ok := s.IsRefreshAllowed(ctx, refreshToken); !ok {
		return "", key.ErrTooManyRequests
	}

	userId, err := s.repository.GetToken(ctx, "refresh:"+refreshToken)
	if err != nil {
		if errors.Is(err, key.ErrExpiredToken) {
			s.logger.Error("service.RefreshToken", zap.Error(err), zap.String("token", refreshToken))
			return "", key.ErrExpiredToken
		}

		return "", fmt.Errorf("service.RefreshToken %w", err)
	}

	user, err := s.repository.FindUser(ctx, userId)

	if err != nil {
		s.logger.Error("service.RefreshToken", zap.Error(err))
		return "", fmt.Errorf("failed fetching user info from refresh token %w", err)
	}

	accessToken, err := s.SignToken(ctx, *user)

	if err != nil {
		s.logger.Error("service.RefreshToken", zap.Error(err))
		return "", fmt.Errorf("RefreshToken failed refreshing token %w", err)
	}

	return accessToken, nil
}

func (s *service) IsRefreshAllowed(ctx context.Context, refreshToken string) bool {
	rate, err := s.repository.GetRefreshRate(ctx, refreshToken)

	if err != nil {
		s.logger.Error("service.IsRefreshAllowed", zap.Error(err))
		return false
	}

	if rate > refreshMaxRate {
		return false
	}

	return true
}

// generate access token in and sign it
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

// Login - gets user and password and returns key token or error if process failed
func (s *service) Login(ctx context.Context, identifier, password string) (*key.Token, error) {

	user, err := s.repository.FindUserByCredentials(ctx, identifier)

	if err != nil {
		s.logger.Error("service.Login unauthorized user", zap.Error(err), zap.String("identifier", identifier))
		return nil, key.ErrInvalidCredentials
	}

	if ok := verifyUser(user.Password, password); !ok {
		s.logger.Error("service.Login unauthorized user", zap.Error(err), zap.String("identifier", identifier))
		return nil, key.ErrInvalidCredentials
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
