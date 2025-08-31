package domain

import (
	"mime/multipart"
	"os"
	"strings"
	"time"
)

type Category string

const (
	Electronics Category = "Electronics"
	Clothing    Category = "Clothing"
	Furniture   Category = "Furniture"
	Vintage     Category = "Vintage"
)

func (c Category) IsValid() bool {
	switch c {
	case Clothing, Electronics, Furniture, Vintage:
		return true
	default:
		return false
	}
}

func AllowedCategories() string {
	return strings.Join([]string{
		string(Electronics),
		string(Clothing),
		string(Furniture),
		string(Vintage),
	}, ",")
}

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
	ID           string
	ItemID       string
	AuctionID    string
	Name         string
	DownloadLink string
}
