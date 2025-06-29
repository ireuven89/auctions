package auction

import (
	"errors"
	"time"
)

type Auction struct {
	ID          string
	Name        string
	Description string
	UserId      string
	Active      bool
	EndTime     int64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type AuctionRequest struct {
	ID          string `json:"-"`
	Name        string `json:"name"`
	Description string `json:"description"`
	UserId      string `json:"user_id"`
	Active      *bool  `json:"active"`
	EndTime     int64  `json:"end_time"`
	CreatedAt   int64  `json:"created_at"`
	UpdatedAt   int64  `json:"updated_at"`
}

var ErrNotFound = errors.New("resource not found")
var DuplicateKey = errors.New("duplicate error key")
