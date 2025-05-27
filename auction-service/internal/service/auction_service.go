package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"os"
	"sync"
	"time"

	aws_config "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/ireuven89/auctions/auction-service/domain"
	"github.com/ireuven89/auctions/shared/config"

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
	FindWithPictures(ctx context.Context, auctionId string) ([]domain.ItemPicture, error)
	Update(ctx context.Context, item domain.ItemRequest) error
	Create(ctx context.Context, auction domain.ItemRequest) error
	CreateBulk(ctx context.Context, request []domain.Item) error
	CreateItemPicture(ctx context.Context, picture domain.ItemPicture) error
	Delete(ctx context.Context, id string) error
	FindByAuctionId(ctx context.Context, auctionId string) ([]domain.Item, error)
}

type ItemPictureRepository interface {
	CreateItemPicture(ctx context.Context, picture domain.ItemPicture) error
	DeleteItemPicture(ctx context.Context, id string) error
	CreateItemPictureBulk(ctx context.Context, picture []domain.ItemPicture) error
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
	CreateAuctionItems(ctx context.Context, auctionId string, items []domain.Item) error
	CreateAuctionPictures(ctx context.Context, id string, request []*multipart.FileHeader) error
	Delete(ctx context.Context, id string) error
	DeleteMany(ctx context.Context, ids []string) error
	PlaceBid(ctx context.Context, bid domain.PlaceBidRequest) error
}

type AuctionService struct {
	repo            Repository
	itemRepo        ItemRepository
	itemPictureRepo ItemPictureRepository
	bidRepo         BidRepository
	logger          *zap.Logger
	awsConfig       config.AWSConfig
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

func (s *AuctionService) CreateAuctionItems(ctx context.Context, auctionId string, items []domain.Item) error {

	for _, item := range items {
		item.ID = generateID()
		item.AuctionID = auctionId
		item.CreatedAt = time.Now()
	}

	if err := s.itemRepo.CreateBulk(ctx, items); err != nil {
		return fmt.Errorf("AuctionService.CreateAuctionItems %w", err)
	}

	return nil
}

func (s *AuctionService) Create(ctx context.Context, auction domain.AuctionRequest) (string, error) {
	//validate auction
	if ok := s.validateAuction(auction); !ok {
		return "", domain.ErrBadRequest
	}

	auction.ID = generateID()
	auction.CreatedAt = time.Now()
	auction.UpdatedAt = time.Now()

	if auction.Regions == nil {
		auction.Regions, _ = json.Marshal("world")
	}

	if err := s.repo.Create(ctx, auction); err != nil {
		s.logger.Error("AuctionService.Failed to create auction ", zap.Error(err))
		return "", fmt.Errorf("AuctionService.Create failed creating %w", err)
	}

	return auction.ID, nil
}

func generateID() string {

	return uuid.New().String()
}

func (s *AuctionService) validateAuction(auction domain.AuctionRequest) bool {

	return auction.Description != "" && auction.InitialOffer != 0 && auction.MinIncrement != 0
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

func (s *AuctionService) CreateAuctionPictures(ctx context.Context, itemId string, files []*multipart.FileHeader) error {
	downloadUrlChannel := make(chan string, len(files))
	var wg *sync.WaitGroup
	for _, file := range files {
		wg.Add(1)
		uploadImageToS3(ctx, file, itemId, s.awsConfig.S3Buckets.Primary, downloadUrlChannel, wg)
	}

	wg.Wait()

	var itemPictures []domain.ItemPicture

	for res := range downloadUrlChannel {
		itemPictures = append(itemPictures, domain.ItemPicture{ID: generateID(), ItemID: itemId, DownloadUrl: res})
	}

	if err := s.itemPictureRepo.CreateItemPictureBulk(ctx, itemPictures); err != nil {
		return fmt.Errorf("AuctionService.CreateAuctionPictures %w", err)
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

func (s *AuctionService) ExecuteInTransaction(ctx context.Context, txFunc func(txCtx context.Context) error) error {

	return nil
}

func (s *AuctionService) validateAuctionForBidding(auction domain.Auction) error {

	if auction.Status != domain.Active {
		return fmt.Errorf("AuctionService.validateAuctionForBidding auction %s is past due bidding", auction.ID)
	}

	return nil
}

// ✅ PURE BUSINESS VALIDATION
func (s *AuctionService) validateBidAmount(amount float64, auction *domain.Auction) error {
	minRequired := auction.CurrentBid + auction.MinIncrement
	if amount < minRequired {
		return fmt.Errorf("minimum bid is %.2f", minRequired)
	}
	return nil
}

// Upload a single image and send its S3 URL through the channel
func uploadImageToS3(ctx context.Context, image *multipart.FileHeader, itemID, bucketName string, urlChan chan string, wg *sync.WaitGroup) {
	defer wg.Done() // Mark goroutine as done

	cfg, err := aws_config.LoadDefaultConfig(context.TODO(), aws_config.WithRegion("region"))
	if err != nil {
		fmt.Println("AWS config error:", err)
		return
	}

	file, err := image.Open()

	if err != nil {
		fmt.Printf("uploadImageToS3 %v", err)
		return
	}

	client := s3.NewFromConfig(cfg)

	// Generate unique file name
	fileKey := fmt.Sprintf("%s/%s", itemID, generateID())

	// Upload to S3
	_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: &bucketName,
		Key:    &fileKey,
		Body:   file,
	})
	if err != nil {
		fmt.Println("Upload failed:", err)
		return
	}

	// Generate public URL
	imageURL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", bucketName, "region", fileKey)

	// Send result through channel
	urlChan <- imageURL
}

func (s *AuctionService) uploadItemImagesToS3(ctx context.Context, images []*os.File, itemID string) ([]string, error) {
	cfg, err := aws_config.LoadDefaultConfig(ctx, aws_config.WithRegion("region"))
	if err != nil {
		return nil, fmt.Errorf("AWS config error: %v", err)
	}

	client := s3.NewFromConfig(cfg)
	var uploadedURLs []string

	for _, image := range images {
		// Generate unique file name
		fileKey := fmt.Sprintf("%s/%s", itemID, uuid.New().String())

		// Upload image
		_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
			Bucket: &s.awsConfig.S3Buckets.Primary,
			Key:    &fileKey,
			Body:   image,
		})
		if err != nil {
			return nil, fmt.Errorf("upload failed: %v", err)
		}

		// Generate public URL'
		imageURL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", "bucketName", "region", fileKey)
		uploadedURLs = append(uploadedURLs, imageURL)

	}

	return uploadedURLs, nil
}
