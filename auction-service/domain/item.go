package domain

import "time"

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
	Pictures    []ItemPicture
}

type ItemPicture struct {
	ID           string
	DownloadLink string
	ItemID       string
}
