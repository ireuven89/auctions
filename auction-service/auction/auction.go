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
}

var ErrNotFound = errors.New("resource not found")
var DBErr = errors.New("db error")
