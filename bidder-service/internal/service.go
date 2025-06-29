package internal

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/ireuven89/auctions/bidder-service/bidder"
	"github.com/ireuven89/auctions/bidder-service/db"
	"go.uber.org/zap"
)

type Service interface {
	GetBidder(ctx context.Context, id string) (bidder.Bidder, error)
	SearchBidders(ctx context.Context, request bidder.BiddersRequest) ([]bidder.Bidder, error)
	DeleteBidder(ctx context.Context, id string) error
	DeleteBidders(ctx context.Context, ids []string) error
	CreateBidder(ctx context.Context, bidder bidder.Bidder) (string, error)
	UpdateBidder(ctx context.Context, bidder bidder.Bidder) error
}

type BidderService struct {
	repo   db.Repository
	logger *zap.Logger
}

func NewService(repo db.Repository, logger *zap.Logger) Service {

	return &BidderService{
		repo:   repo,
		logger: logger,
	}
}

func (s *BidderService) GetBidder(ctx context.Context, id string) (bidder.Bidder, error) {
	result, err := s.repo.Find(ctx, id)

	if err != nil {
		s.logger.Error("BidderService.GetBidder failed updating", zap.Error(err))
		return bidder.Bidder{}, fmt.Errorf("BidderService.GetBidder %w", err)
	}
	return result, nil
}

func (s *BidderService) CreateBidder(ctx context.Context, bidder bidder.Bidder) (string, error) {
	bidder.ID = generateID()

	if err := s.repo.Create(ctx, bidder); err != nil {
		s.logger.Error("BidderService.CreateBidder failed updating", zap.Error(err))
		return "", fmt.Errorf("BidderService.CreateBidder %w", err)
	}

	return bidder.ID, nil
}

func generateID() string {

	return uuid.New().String()
}

func (s *BidderService) UpdateBidder(ctx context.Context, bidder bidder.Bidder) error {

	if err := s.repo.Update(ctx, bidder); err != nil {
		s.logger.Error("BidderService.UpdateBidder failed updating", zap.Error(err))
		return fmt.Errorf("BidderService.UpdateBidder %w", err)
	}
	return nil
}

func (s *BidderService) DeleteBidder(ctx context.Context, id string) error {

	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Error("BidderService.DeleteBidder failed deleting ", zap.Error(err), zap.String("id", id))
		return fmt.Errorf("BidderService.DeleteBidder %w", err)
	}

	return nil
}

func (s *BidderService) SearchBidders(ctx context.Context, req bidder.BiddersRequest) ([]bidder.Bidder, error) {

	bidders, err := s.repo.FindAll(ctx, req)
	if err != nil {
		s.logger.Error("BidderService.DeleteBidder failed deleting ", zap.Error(err))
		return nil, fmt.Errorf("BidderService.DeleteBidder %w", err)
	}

	return bidders, nil
}

func (s *BidderService) DeleteBidders(ctx context.Context, ids []string) error {
	var vals []interface{}

	for _, id := range ids {
		vals = append(vals, id)
	}

	if err := s.repo.DeleteMany(ctx, vals); err != nil {
		return fmt.Errorf("BidderService.DeleteBidders %w", err)
	}

	return nil
}
