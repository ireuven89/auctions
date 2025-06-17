package domain

import "time"

type Bid struct {
	ID        string
	AuctionID string
	BidderID  string
	Price     float64
	CreateAt  time.Time
	Winner    bool
}

type PlaceBidRequest struct {
	ID        string
	AuctionID string
	BidderID  string
	Amount    float64
	CreateAt  time.Time
	Winner    bool
}
