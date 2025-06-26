package internal

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/ireuven89/auctions/auction-service/auction"
	"github.com/ireuven89/auctions/auction-service/db"
	"go.uber.org/zap"
)

type Service interface {
	Fetch(ctx context.Context, id string) (*auction.Auction, error)
	Search(ctx context.Context, request auction.AuctionRequest) ([]auction.Auction, error)
	Update(ctx context.Context, auction auction.AuctionRequest) error
	Create(ctx context.Context, auction auction.AuctionRequest) (string, error)
	Delete(ctx context.Context, id string) error
	DeleteMany(ctx context.Context, ids []string) error
}

type AuctionService struct {
	repo   db.Repository
	logger *zap.Logger
}

func NewService(repo db.Repository, logger *zap.Logger) Service {

	return &AuctionService{
		logger: logger,
		repo:   repo,
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

	if err := s.repo.Create(ctx, auction); err != nil {
		s.logger.Error("AuctionService.Failed to create auction ", zap.Error(err))
		return "", fmt.Errorf("AuctionService.Create failed creating %w", err)
	}

	return auction.ID, nil
}

func validateAuction(auction auction.AuctionRequest) error {
	if auction.Item == "" {
		return fmt.Errorf("AuctionService.ValidateAuction failed invalid auction item")
	}

	return nil
}

func generateID() string {

	return uuid.New().String()
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
