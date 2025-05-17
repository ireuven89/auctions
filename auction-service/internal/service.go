package internal

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/ireuven89/auctions/auction-service/auction"
	"github.com/ireuven89/auctions/auction-service/db"
	"go.uber.org/zap"
)

type Service interface {
	Fetch(ctx context.Context, id string) (*auction.Auction, error)
	Update(ctx context.Context, auction auction.AuctionRequest) error
	Create(ctx context.Context, auction auction.AuctionRequest) (string, error)
	Delete(ctx context.Context, id string) error
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
		s.logger.Error("AuctionService failed to fetch auction")
		if err == sql.ErrNoRows {
			return nil, auction.ErrNotFound
		}
		return nil, err
	}

	return &res, nil
}

func (s *AuctionService) Update(ctx context.Context, auction auction.AuctionRequest) error {

	if err := s.repo.Update(ctx, auction); err != nil {
		s.logger.Error("AuctionService failed to update auction")
		return err
	}

	return nil
}

func (s *AuctionService) Create(ctx context.Context, auction auction.AuctionRequest) (string, error) {
	id := uuid.New().String()
	auction.ID = id

	if err := s.repo.Create(ctx, auction); err != nil {
		s.logger.Error("failed to create auction ", zap.Error(err))
		return "", err
	}

	return id, nil
}

func (s *AuctionService) Delete(ctx context.Context, id string) error {

	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}
