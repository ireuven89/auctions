package domain

import (
	"mime/multipart"
	"os"
	"time"
)

type Item struct {
	ID          string
	Description string
	AuctionID   string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type ItemPictureResponse struct {
	ID           string
	Description  string
	AuctionID    string
	DownloadLink string
}

type ItemRequest struct {
	ID          string
	Description string
	AuctionID   string
	Pictures    []*os.File
}

type ItemPictureRequest struct {
	ID     string
	File   *multipart.FileHeader
	ItemID string
}

type ItemPicture struct {
	ID          string
	DownloadUrl string
	ItemID      string
}
