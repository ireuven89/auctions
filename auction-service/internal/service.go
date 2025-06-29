package internal

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ireuven89/auctions/auction-service/auction"
	"github.com/ireuven89/auctions/auction-service/db"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type Service interface {
	Fetch(ctx context.Context, id string) (*auction.Auction, error)
	Search(ctx context.Context, request auction.AuctionRequest) ([]auction.Auction, error)
	Update(ctx context.Context, auction auction.AuctionRequest) error
	Activate(ctx context.Context, auction auction.AuctionRequest) error
	Create(ctx context.Context, auction auction.AuctionRequest) (string, error)
	Delete(ctx context.Context, id string) error
	DeleteMany(ctx context.Context, ids []string) error
}

var (
	activeAuctionsQueue    = "active_auctions"
	expired_auctions_queue = "expired_auctions"
)

type AuctionService struct {
	repo   db.Repository
	redis  *redis.Client
	logger *zap.Logger
}

func NewService(repo db.Repository, redisConn *redis.Client, logger *zap.Logger) Service {

	return &AuctionService{
		logger: logger,
		repo:   repo,
		redis:  redisConn,
	}
}

func (s *AuctionService) Fetch(ctx context.Context, id string) (*auction.Auction, error) {
	res, err := s.repo.Find(ctx, id)

	if err != nil {
		s.logger.Error("AuctionService failed to fetch auction", zap.Error(err))
		if errors.Is(err, sql.ErrNoRows) {
			return nil, auction.ErrNotFound
		}
		return nil, fmt.Errorf("AuctionService.Fetch failed fetching bidder %w", err)
	}

	return &res, nil
}

func (s *AuctionService) Search(ctx context.Context, request auction.AuctionRequest) ([]auction.Auction, error) {
	res, err := s.repo.FindAll(ctx, request)

	if err != nil {
		s.logger.Error("AuctionService.Search failed to search", zap.Error(err))
		return nil, fmt.Errorf("AuctionService.Search %w", err)
	}

	return res, nil
}

func (s *AuctionService) Update(ctx context.Context, auction auction.AuctionRequest) error {

	if err := s.repo.Update(ctx, auction); err != nil {
		s.logger.Error("AuctionService failed to update auction", zap.Error(err))
		return fmt.Errorf("AuctionService.Update failed updating %w", err)
	}

	return nil
}

func (s *AuctionService) Create(ctx context.Context, auction auction.AuctionRequest) (string, error) {
	if err := validateAuction(auction); err != nil {
		return "", err
	}

	auction.ID = generateID()
	auction.CreatedAt = time.Now().Unix()
	auction.UpdatedAt = time.Now().Unix()

	if err := s.repo.Create(ctx, auction); err != nil {
		s.logger.Error("AuctionService.Failed to create auction ", zap.Error(err))
		return "", fmt.Errorf("AuctionService.Create failed creating %w", err)
	}

	return auction.ID, nil
}

func validateAuction(auction auction.AuctionRequest) error {
	if auction.Description == "" {
		return fmt.Errorf("AuctionService.ValidateAuction failed invalid auction description")
	}
	if auction.Name == "" {
		return fmt.Errorf("AuctionService.ValidateAuction failed invalid auction name")
	}

	if auction.UserId == "" {
		return fmt.Errorf("AuctionService.ValidateAuction failed invalid auction user_id")
	}

	return nil
}

func generateID() string {

	return uuid.New().String()
}

func (s *AuctionService) Activate(ctx context.Context, auction auction.AuctionRequest) error {
	if err := validateActivateAuction(auction); err != nil {
		return fmt.Errorf("AuctionService.Activate failed invalid auction")
	}

	err := s.repo.WithTransactionContext(ctx, func(txRepo db.Repository) error {
		if err := s.repo.Update(ctx, auction); err != nil {
			return fmt.Errorf("AuctionService.Activate failed activating auction %w", err)
		}

		//publish to redis the auction
		key := fmt.Sprintf("%s:%s", "queue", activeAuctionsQueue)
		if err := s.redis.LPush(ctx, key, auction.ID, time.Until(time.UnixMilli(auction.EndTime))).Err(); err != nil {
			return fmt.Errorf("AuctionService.Activate failed activating auction %w", err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("AuctionService.Activate failed activating auction %w", err)
	}

	return nil
}

func validateActivateAuction(auction auction.AuctionRequest) error {
	if auction.ID == "" {
		return fmt.Errorf("AuctionService.Activate failed invalid auction ID")
	}

	if auction.EndTime <= time.Now().Unix() {
		return fmt.Errorf("AuctionService.Activate failed invalid auction end_time")
	}

	return nil
}

func (s *AuctionService) Delete(ctx context.Context, id string) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Error("AuctionService.Delete failed deleting bidder", zap.Error(err))
		return fmt.Errorf("AuctionService.Delete failed deleting %w", err)
	}

	return nil
}

func (s *AuctionService) DeleteMany(ctx context.Context, ids []string) error {
	var vals []interface{}

	for _, id := range ids {
		vals = append(vals, id)
	}

	if err := s.repo.DeleteMany(ctx, vals); err != nil {
		s.logger.Error("AuctionService.DeleteMany failed deleting ids", zap.Error(err))
		return fmt.Errorf("AuctionService.DeleteMany failed deleting %w", err)
	}

	return nil
}
