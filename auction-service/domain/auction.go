package domain

import (
	"encoding/json"
	"errors"
	"time"
)

type Auction struct {
	ID           string
	Description  string
	Regions      []byte
	InitialOffer float64
	CurrentBid   float64
	MinIncrement float64
	Status       AuctionStatus
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type AuctionRequest struct {
	ID           string          `json:"-"`
	Description  string          `json:"description"`
	Regions      json.RawMessage `json:"regions"`
	InitialOffer int64           `json:"initialOffer"`
	Status       string          `json:"status"`
	SellerId     string          `json:"sellerId"`
	WinnerId     string          `json:"winnerId"`
	CurrentBid   float64         `json:"currentBid"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
	//	Items        []ItemRequest   `json:"items"`
}

var ErrNotFound = errors.New("resource not found")
var ErrTooManyRequests = errors.New("too many requests")

type AuctionStatus int

const (
	Pending AuctionStatus = iota
	Active
	Completed
	Cancelled
)

// Optionally, implement Stringer interface for pretty printing
func (s AuctionStatus) String() string {
	switch s {
	case Pending:
		return "Pending"
	case Active:
		return "Active"
	case Completed:
		return "Completed"
	case Cancelled:
		return "Cancelled"
	default:
		return "Unknown"
	}
}

func FromString(status string) AuctionStatus {
	switch status {
	case "Pending":
		return Pending
	case "Active":
		return Active
	case "Completed":
		return Completed
	case "Cancelled":
		return Cancelled
	default:
		return Pending
	}
}
