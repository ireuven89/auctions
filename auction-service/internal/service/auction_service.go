package service

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/ireuven89/auctions/auction-service/domain"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Repository interface {
	Find(ctx context.Context, id string) (domain.Auction, error)
	FindAll(ctx context.Context, request domain.AuctionRequest) ([]domain.Auction, error)
	Update(ctx context.Context, auction domain.AuctionRequest) error
	Create(ctx context.Context, auction domain.AuctionRequest) error
	Delete(ctx context.Context, id string) error
	DeleteMany(ctx context.Context, ids []interface{}) error
}

type ItemRepository interface {
	Find(ctx context.Context, id string) (domain.Item, error)
	FindWithPictures(ctx context.Context, auctionId string) ([]domain.ItemPictureResponse, error)
	Update(ctx context.Context, item domain.ItemRequest) error
	Create(ctx context.Context, auction domain.ItemRequest) error
	CreateBulk(ctx context.Context, request []domain.ItemRequest) error
	Delete(ctx context.Context, id string) error
	FindByAuctionId(ctx context.Context, auctionId string) ([]domain.Item, error)
}

type BidRepository interface {
	Find(ctx context.Context, id string) (domain.Bid, error)
	Create(ctx context.Context, bid domain.Bid) error
}

type Service interface {
	Fetch(ctx context.Context, id string) (*domain.Auction, error)
	Search(ctx context.Context, request domain.AuctionRequest) ([]domain.Auction, error)
	Update(ctx context.Context, auction domain.AuctionRequest) error
	Create(ctx context.Context, auction domain.AuctionRequest) (string, error)
	Delete(ctx context.Context, id string) error
	DeleteMany(ctx context.Context, ids []string) error
	PlaceBid(ctx context.Context, bid domain.PlaceBidRequest) error
}

type AuctionService struct {
	repo     Repository
	itemRepo ItemRepository
	bidRepo  BidRepository
	logger   *zap.Logger
}

func NewService(repo Repository, itemRepo ItemRepository, logger *zap.Logger) Service {

	return &AuctionService{
		logger:   logger,
		repo:     repo,
		itemRepo: itemRepo,
	}
}

func (s *AuctionService) Fetch(ctx context.Context, id string) (*domain.Auction, error) {
	res, err := s.repo.Find(ctx, id)

	if err != nil {
		s.logger.Error("AuctionService failed to fetch auction", zap.Error(err))
		if err == sql.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("AuctionService.Fetch failed fetching bidder %w", err)
	}

	return &res, nil
}

func (s *AuctionService) Search(ctx context.Context, request domain.AuctionRequest) ([]domain.Auction, error) {

	res, err := s.repo.FindAll(ctx, request)

	if err != nil {
		s.logger.Error("AuctionService.Search failed to search", zap.Error(err))
		return nil, fmt.Errorf("AuctionService.Search %w", err)
	}

	return res, nil
}

func (s *AuctionService) Update(ctx context.Context, auction domain.AuctionRequest) error {
	auction.UpdatedAt = time.Now()
	if err := s.repo.Update(ctx, auction); err != nil {
		s.logger.Error("AuctionService failed to update auction", zap.Error(err))
		return fmt.Errorf("AuctionService.Update failed updating %w", err)
	}

	return nil
}

func (s *AuctionService) Create(ctx context.Context, auction domain.AuctionRequest) (string, error) {
	auction.ID = generateID()
	auction.CreatedAt = time.Now()
	auction.UpdatedAt = time.Now()

	if err := s.repo.Create(ctx, auction); err != nil {
		s.logger.Error("AuctionService.Failed to create auction ", zap.Error(err))
		return "", fmt.Errorf("AuctionService.Create failed creating %w", err)
	}

	for _, i := range auction.Items {
		i.ID = generateID()
		for _, p := range i.Pictures {
			p.ID = generateID()
		}
	}

	if err := s.itemRepo.CreateBulk(ctx, auction.Items); err != nil {
		return "", fmt.Errorf("")
	}

	return auction.ID, nil
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

func (s *AuctionService) FindAuctionItems(ctx context.Context, auctionId string) ([]domain.Item, error) {

	res, err := s.itemRepo.FindByAuctionId(ctx, auctionId)

	if err != nil {
		s.logger.Error("AuctionService", zap.Error(err), zap.String("id", auctionId))
		return nil, fmt.Errorf("AuctionService.FindAuctionItems %w", err)
	}

	return res, nil
}

func (s *AuctionService) PlaceBid(ctx context.Context, req domain.PlaceBidRequest) error {
	return s.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		// 1. Get auction (delegates to repository)
		auction, err := s.repo.Find(ctx, req.AuctionID)

		if err = s.validateAuctionForBidding(auction); err != nil {
			return err
		}

		if err = s.validateBidAmount(req.Amount, &auction); err != nil {
			return err
		}

		// All database operations delegated to repositories
		bid := domain.Bid{
			ID:        generateID(),
			AuctionID: auction.ID,
			CreateAt:  time.Now(),
			Winner:    false,
			BidderID:  req.BidderID,
			Price:     req.Amount,
		}
		if err = s.bidRepo.Create(ctx, bid); err != nil {
			return err
		}
		auction.CurrentBid = req.Amount
		auctionUpdate := domain.AuctionRequest{ID: auction.ID, CurrentBid: req.Amount}
		if err = s.repo.Update(ctx, auctionUpdate); err != nil {
			return err
		}

		return nil
	})
}

func (s *AuctionService) ExecuteInTransaction(context.Context, func(ctx context.Context) error) error {

	return nil
}

func (s *AuctionService) validateAuctionForBidding(auction domain.Auction) error {

	if auction.Status != domain.Active {
		return fmt.Errorf("AuctionService.validateAuctionForBidding failed")
	}

	return nil
}

// ✅ PURE BUSINESS VALIDATION - No Database Calls
func (s *AuctionService) validateBidAmount(amount float64, auction *domain.Auction) error {
	minRequired := auction.CurrentBid + auction.MinIncrement
	if amount < minRequired {
		return fmt.Errorf("minimum bid is %.2f", minRequired)
	}
	return nil
}
