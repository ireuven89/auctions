package auction

import "errors"

type Auction struct {
	ID       string
	Name     string
	BidderId string
}

type AuctionRequest struct {
	ID   string `json:"-"`
	Name string `json:"name"`
	Item string `json:"item"`
}

var ErrNotFound = errors.New("resource not found")
var DuplicateKey = errors.New("duplicate error key")
