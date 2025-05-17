package db

import "github.com/ireuven89/auctions/bidder-service/bidder"

type BidderDb struct {
	ID   string `db:"id"`
	Name string `db:"name"`
	Item string `db:"item"`
}

func toBidder(db BidderDb) bidder.Bidder {

	return bidder.Bidder{
		ID:   db.ID,
		Name: db.Name,
		Item: db.Item,
	}
}

func fromBidder(bidder bidder.Bidder) BidderDb {

	return BidderDb{
		ID:   bidder.ID,
		Name: bidder.Name,
		Item: bidder.Item,
	}
}
