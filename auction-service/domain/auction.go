package domain

import (
	"encoding/json"
	"fmt"
	"time"
)

type Auction struct {
	ID           string
	Description  string
	Category     Category
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
	MinIncrement int64           `json:"minIncrement"`
	Status       string          `json:"status"`
	Category     Category        `json:"category"`
	SellerId     string          `json:"sellerId"`
	WinnerId     string          `json:"winnerId"`
	CurrentBid   float64         `json:"currentBid"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
	Items        []ItemRequest   `json:"items"`
}

type AppError struct {
	Kind    string
	Message string
}

func (e *AppError) Error() string {
	return e.Message
}

func ErrNotFound(resource string) error {
	return &AppError{Kind: "not_found", Message: fmt.Sprintf("%s not found", resource)}
}
func ErrTooManyRequests(action string) error {
	return &AppError{Kind: "too_many_requests", Message: fmt.Sprintf("too many requests while %s", action)}
}

func ErrUnAuthorized(user string) error {
	return &AppError{Kind: "unauthorized", Message: fmt.Sprintf("unauthorized access by %s", user)}
}

func ErrBadRequest(reason string) error {
	return &AppError{Kind: "bad_request", Message: fmt.Sprintf("bad request: %s", reason)}
}

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
